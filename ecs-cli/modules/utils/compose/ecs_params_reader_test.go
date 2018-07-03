// Copyright 2015-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package utils

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/docker/libcompose/yaml"
	"github.com/stretchr/testify/assert"
)

func TestReadECSParams(t *testing.T) {
	ecsParamsString := `version: 1
task_definition:
  ecs_network_mode: host
  task_role_arn: arn:aws:iam::123456789012:role/my_role`

	content := []byte(ecsParamsString)

	tmpfile, err := ioutil.TempFile("", "ecs-params")
	assert.NoError(t, err, "Could not create ecs fields tempfile")

	ecsParamsFileName := tmpfile.Name()
	defer os.Remove(ecsParamsFileName)

	_, err = tmpfile.Write(content)
	assert.NoError(t, err, "Could not write data to ecs fields tempfile")

	err = tmpfile.Close()
	assert.NoError(t, err, "Could not close tempfile")

	ecsParams, err := ReadECSParams(ecsParamsFileName)

	if assert.NoError(t, err) {
		assert.Equal(t, "1", ecsParams.Version, "Expected version to match")
		taskDef := ecsParams.TaskDefinition
		assert.Equal(t, "host", taskDef.NetworkMode, "Expected network mode to match")
		assert.Equal(t, "arn:aws:iam::123456789012:role/my_role", taskDef.TaskRoleArn, "Expected task role ARN to match")
		// Should still populate other fields with empty values
		assert.Empty(t, taskDef.ExecutionRole)
		awsvpcConfig := ecsParams.RunParams.NetworkConfiguration.AwsVpcConfiguration
		assert.Empty(t, awsvpcConfig.Subnets)
		assert.Empty(t, awsvpcConfig.SecurityGroups)
	}
}

func TestReadECSParams_FileDoesNotExist(t *testing.T) {
	_, err := ReadECSParams("nonexistant.yml")
	assert.Error(t, err)
}

func TestReadECSParams_NoFile(t *testing.T) {
	ecsParams, err := ReadECSParams("")
	if assert.NoError(t, err) {
		assert.Nil(t, ecsParams)
	}
}

func TestReadECSParams_WithServices(t *testing.T) {
	ecsParamsString := `version: 1
task_definition:
  ecs_network_mode: host
  task_role_arn: arn:aws:iam::123456789012:role/my_role
  services:
    mysql:
      essential: false
      cpu_shares: 100
      mem_limit: 524288000
      mem_reservation: 500mb
    wordpress:
      essential: true`

	content := []byte(ecsParamsString)

	tmpfile, err := ioutil.TempFile("", "ecs-params")
	assert.NoError(t, err, "Could not create ecs fields tempfile")

	ecsParamsFileName := tmpfile.Name()
	defer os.Remove(ecsParamsFileName)

	_, err = tmpfile.Write(content)
	assert.NoError(t, err, "Could not write data to ecs fields tempfile")

	err = tmpfile.Close()
	assert.NoError(t, err, "Could not close tempfile")

	ecsParams, err := ReadECSParams(ecsParamsFileName)

	if assert.NoError(t, err) {
		taskDef := ecsParams.TaskDefinition
		assert.Equal(t, "host", ecsParams.TaskDefinition.NetworkMode, "Expected NetworkMode to match")
		assert.Equal(t, "arn:aws:iam::123456789012:role/my_role", taskDef.TaskRoleArn, "Expected TaskRoleArn to match")

		containerDefs := taskDef.ContainerDefinitions
		assert.Equal(t, 2, len(containerDefs), "Expected 2 containers")

		mysql := containerDefs["mysql"]
		wordpress := containerDefs["wordpress"]

		assert.False(t, mysql.Essential, "Expected container to not be essential")
		assert.Equal(t, int64(100), mysql.Cpu)
		assert.Equal(t, yaml.MemStringorInt(524288000), mysql.Memory)
		assert.Equal(t, yaml.MemStringorInt(524288000), mysql.MemoryReservation)
		assert.True(t, wordpress.Essential, "Expected container to be essential")
	}
}

func TestReadECSParams_WithRunParams(t *testing.T) {
	ecsParamsString := `version: 1
task_definition:
  ecs_network_mode: awsvpc
run_params:
  network_configuration:
    awsvpc_configuration:
      subnets: [subnet-feedface, subnet-deadbeef]
      security_groups:
        - sg-bafff1ed
        - sg-c0ffeefe`

	content := []byte(ecsParamsString)

	tmpfile, err := ioutil.TempFile("", "ecs-params")
	assert.NoError(t, err, "Could not create ecs fields tempfile")

	ecsParamsFileName := tmpfile.Name()
	defer os.Remove(ecsParamsFileName)

	_, err = tmpfile.Write(content)
	assert.NoError(t, err, "Could not write data to ecs fields tempfile")

	err = tmpfile.Close()
	assert.NoError(t, err, "Could not close tempfile")

	ecsParams, err := ReadECSParams(ecsParamsFileName)

	if assert.NoError(t, err) {
		taskDef := ecsParams.TaskDefinition
		assert.Equal(t, "awsvpc", taskDef.NetworkMode, "Expected network mode to match")

		awsvpcConfig := ecsParams.RunParams.NetworkConfiguration.AwsVpcConfiguration
		assert.Equal(t, 2, len(awsvpcConfig.Subnets), "Expected 2 subnets")
		assert.Equal(t, []string{"subnet-feedface", "subnet-deadbeef"}, awsvpcConfig.Subnets, "Expected subnets to match")
		assert.Equal(t, 2, len(awsvpcConfig.SecurityGroups), "Expected 2 securityGroups")
		assert.Equal(t, []string{"sg-bafff1ed", "sg-c0ffeefe"}, awsvpcConfig.SecurityGroups, "Expected security groups to match")
		assert.Equal(t, AssignPublicIp(""), awsvpcConfig.AssignPublicIp, "Expected AssignPublicIP to be empty")
	}
}

// Task Size, Task Execution Role, and Assign Public Ip are required for Fargate tasks
func TestReadECSParams_WithFargateRunParams(t *testing.T) {
	ecsParamsString := `version: 1
task_definition:
  ecs_network_mode: awsvpc
  task_execution_role: arn:aws:iam::123456789012:role/fargate_role
  task_size:
    mem_limit: 0.5GB
    cpu_limit: 256
run_params:
  network_configuration:
    awsvpc_configuration:
      subnets: [subnet-feedface, subnet-deadbeef]
      security_groups:
        - sg-bafff1ed
        - sg-c0ffeefe
      assign_public_ip: ENABLED`

	content := []byte(ecsParamsString)

	tmpfile, err := ioutil.TempFile("", "ecs-params")
	assert.NoError(t, err, "Could not create ecs fields tempfile")

	ecsParamsFileName := tmpfile.Name()
	defer os.Remove(ecsParamsFileName)

	_, err = tmpfile.Write(content)
	assert.NoError(t, err, "Could not write data to ecs fields tempfile")

	err = tmpfile.Close()
	assert.NoError(t, err, "Could not close tempfile")

	ecsParams, err := ReadECSParams(ecsParamsFileName)

	if assert.NoError(t, err) {
		taskDef := ecsParams.TaskDefinition
		assert.Equal(t, "awsvpc", taskDef.NetworkMode, "Expected network mode to match")
		assert.Equal(t, "arn:aws:iam::123456789012:role/fargate_role", taskDef.ExecutionRole)
		assert.Equal(t, "0.5GB", taskDef.TaskSize.Memory)
		assert.Equal(t, "256", taskDef.TaskSize.Cpu)

		awsvpcConfig := ecsParams.RunParams.NetworkConfiguration.AwsVpcConfiguration
		assert.Equal(t, 2, len(awsvpcConfig.Subnets), "Expected 2 subnets")
		assert.Equal(t, []string{"subnet-feedface", "subnet-deadbeef"}, awsvpcConfig.Subnets, "Expected subnets to match")
		assert.Equal(t, 2, len(awsvpcConfig.SecurityGroups), "Expected 2 securityGroups")
		assert.Equal(t, []string{"sg-bafff1ed", "sg-c0ffeefe"}, awsvpcConfig.SecurityGroups, "Expected security groups to match")
		assert.Equal(t, Enabled, awsvpcConfig.AssignPublicIp, "Expected AssignPublicIp to match")
	}
}

func TestReadECSParams_MemoryWithUnits(t *testing.T) {
	ecsParamsString := `version: 1
task_definition:
  ecs_network_mode: awsvpc
  task_size:
    mem_limit: 0.5GB
    cpu_limit: 256`

	content := []byte(ecsParamsString)

	tmpfile, err := ioutil.TempFile("", "ecs-params")
	assert.NoError(t, err, "Could not create ecs fields tempfile")

	ecsParamsFileName := tmpfile.Name()
	defer os.Remove(ecsParamsFileName)

	_, err = tmpfile.Write(content)
	assert.NoError(t, err, "Could not write data to ecs fields tempfile")

	err = tmpfile.Close()
	assert.NoError(t, err, "Could not close tempfile")

	ecsParams, err := ReadECSParams(ecsParamsFileName)

	if assert.NoError(t, err) {
		taskSize := ecsParams.TaskDefinition.TaskSize
		assert.Equal(t, "256", taskSize.Cpu, "Expected CPU limit to match")
		assert.Equal(t, "0.5GB", taskSize.Memory, "Expected Memory limit to match")
	}
}

// Task Size must match specific CPU/Memory buckets, but we leave validation to ECS.
func TestReadECSParams_WithTaskSize(t *testing.T) {
	ecsParamsString := `version: 1
task_definition:
  task_size:
    mem_limit: 1024
    cpu_limit: 256`

	content := []byte(ecsParamsString)

	tmpfile, err := ioutil.TempFile("", "ecs-params")
	assert.NoError(t, err, "Could not create ecs fields tempfile")

	ecsParamsFileName := tmpfile.Name()
	defer os.Remove(ecsParamsFileName)

	_, err = tmpfile.Write(content)
	assert.NoError(t, err, "Could not write data to ecs fields tempfile")

	err = tmpfile.Close()
	assert.NoError(t, err, "Could not close tempfile")

	ecsParams, err := ReadECSParams(ecsParamsFileName)

	if assert.NoError(t, err) {
		taskSize := ecsParams.TaskDefinition.TaskSize
		assert.Equal(t, "256", taskSize.Cpu, "Expected CPU limit to match")
		assert.Equal(t, "1024", taskSize.Memory, "Expected Memory limit to match")
	}
}

/** ConvertToECSNetworkConfiguration tests **/

func TestConvertToECSNetworkConfiguration(t *testing.T) {
	taskDef := EcsTaskDef{NetworkMode: "awsvpc"}
	subnets := []string{"subnet-feedface"}
	securityGroups := []string{"sg-c0ffeefe"}
	awsVpconfig := AwsVpcConfiguration{
		Subnets:        subnets,
		SecurityGroups: securityGroups,
	}

	networkConfig := NetworkConfiguration{
		AwsVpcConfiguration: awsVpconfig,
	}

	ecsParams := &ECSParams{
		TaskDefinition: taskDef,
		RunParams: RunParams{
			NetworkConfiguration: networkConfig,
		},
	}

	ecsNetworkConfig, err := ConvertToECSNetworkConfiguration(ecsParams)

	if assert.NoError(t, err) {
		ecsAwsConfig := ecsNetworkConfig.AwsvpcConfiguration
		assert.Equal(t, subnets[0], aws.StringValue(ecsAwsConfig.Subnets[0]), "Expected subnets to match")
		assert.Equal(t, securityGroups[0], aws.StringValue(ecsAwsConfig.SecurityGroups[0]), "Expected securityGroups to match")
		assert.Nil(t, ecsAwsConfig.AssignPublicIp, "Expected AssignPublicIp to be nil")
	}
}

func TestConvertToECSNetworkConfiguration_NoSecurityGroups(t *testing.T) {
	taskDef := EcsTaskDef{NetworkMode: "awsvpc"}
	subnets := []string{"subnet-feedface"}
	awsVpconfig := AwsVpcConfiguration{
		Subnets: subnets,
	}

	networkConfig := NetworkConfiguration{
		AwsVpcConfiguration: awsVpconfig,
	}

	ecsParams := &ECSParams{
		TaskDefinition: taskDef,
		RunParams: RunParams{
			NetworkConfiguration: networkConfig,
		},
	}

	ecsNetworkConfig, err := ConvertToECSNetworkConfiguration(ecsParams)

	if assert.NoError(t, err) {
		ecsAwsConfig := ecsNetworkConfig.AwsvpcConfiguration
		assert.Equal(t, subnets[0], aws.StringValue(ecsAwsConfig.Subnets[0]), "Expected subnets to match")
		assert.Nil(t, ecsAwsConfig.AssignPublicIp, "Expected AssignPublicIp to be nil")
	}
}

func TestConvertToECSNetworkConfiguration_ErrorWhenNoSubnets(t *testing.T) {
	taskDef := EcsTaskDef{NetworkMode: "awsvpc"}
	subnets := []string{}

	awsVpconfig := AwsVpcConfiguration{
		Subnets: subnets,
	}

	networkConfig := NetworkConfiguration{
		AwsVpcConfiguration: awsVpconfig,
	}

	ecsParams := &ECSParams{
		TaskDefinition: taskDef,
		RunParams: RunParams{
			NetworkConfiguration: networkConfig,
		},
	}

	_, err := ConvertToECSNetworkConfiguration(ecsParams)

	assert.Error(t, err)
}

func TestConvertToECSNetworkConfiguration_WhenNoECSParams(t *testing.T) {
	ecsParams, err := ConvertToECSNetworkConfiguration(nil)

	if assert.NoError(t, err) {
		assert.Nil(t, ecsParams)
	}
}

func TestConvertToECSNetworkConfiguration_WithAssignPublicIp(t *testing.T) {
	taskDef := EcsTaskDef{NetworkMode: "awsvpc"}
	subnets := []string{"subnet-feedface"}
	awsVpconfig := AwsVpcConfiguration{
		Subnets:        subnets,
		AssignPublicIp: Enabled,
	}

	networkConfig := NetworkConfiguration{
		AwsVpcConfiguration: awsVpconfig,
	}

	ecsParams := &ECSParams{
		TaskDefinition: taskDef,
		RunParams: RunParams{
			NetworkConfiguration: networkConfig,
		},
	}

	ecsNetworkConfig, err := ConvertToECSNetworkConfiguration(ecsParams)

	if assert.NoError(t, err) {
		ecsAwsConfig := ecsNetworkConfig.AwsvpcConfiguration
		assert.Equal(t, subnets[0], aws.StringValue(ecsAwsConfig.Subnets[0]), "Expected subnets to match")
		assert.Equal(t, "ENABLED", aws.StringValue(ecsAwsConfig.AssignPublicIp), "Expected AssignPublicIp to match")
	}
}

func TestConvertToECSNetworkConfiguration_NoNetworkConfig(t *testing.T) {
	taskDef := EcsTaskDef{NetworkMode: "bridge"}

	ecsParams := &ECSParams{
		TaskDefinition: taskDef,
	}

	ecsNetworkConfig, err := ConvertToECSNetworkConfiguration(ecsParams)

	if assert.NoError(t, err) {
		assert.Nil(t, ecsNetworkConfig, "Expected AssignPublicIp to be nil")
	}
}

func TestReadECSParams_WithHealthCheck(t *testing.T) {
	ecsParamsString := `version: 1
task_definition:
  services:
    mysql:
      healthcheck:
        test: ["CMD", "curl", "-f", "http://localhost"]
        interval: 1m30s
        timeout: 10s
        retries: 3
        start_period: 40s
    wordpress:
      healthcheck:
        command: ["CMD-SHELL", "curl -f http://localhost"]
        interval: 70
        timeout: 15
        retries: 5
        start_period: 40
    logstash:
      healthcheck:
        test: curl -f http://localhost
        interval: 10m
        timeout: 15s
        retries: 5
        start_period: 50`

	mysqlExpectedHealthCheck := ecs.HealthCheck{
		Command:     aws.StringSlice([]string{"CMD", "curl", "-f", "http://localhost"}),
		Interval:    aws.Int64(int64(90)),
		Timeout:     aws.Int64(int64(10)),
		Retries:     aws.Int64(int64(3)),
		StartPeriod: aws.Int64(int64(40)),
	}

	wordpressExpectedHealthCheck := ecs.HealthCheck{
		Command:     aws.StringSlice([]string{"CMD-SHELL", "curl -f http://localhost"}),
		Interval:    aws.Int64(int64(70)),
		Timeout:     aws.Int64(int64(15)),
		Retries:     aws.Int64(int64(5)),
		StartPeriod: aws.Int64(int64(40)),
	}

	logstashExpectedHealthCheck := ecs.HealthCheck{
		Command:     aws.StringSlice([]string{"CMD-SHELL", "curl -f http://localhost"}),
		Interval:    aws.Int64(int64(600)),
		Timeout:     aws.Int64(int64(15)),
		Retries:     aws.Int64(int64(5)),
		StartPeriod: aws.Int64(int64(50)),
	}

	content := []byte(ecsParamsString)

	tmpfile, err := ioutil.TempFile("", "ecs-params")
	assert.NoError(t, err, "Could not create ecs fields tempfile")

	ecsParamsFileName := tmpfile.Name()
	defer os.Remove(ecsParamsFileName)

	_, err = tmpfile.Write(content)
	assert.NoError(t, err, "Could not write data to ecs fields tempfile")

	err = tmpfile.Close()
	assert.NoError(t, err, "Could not close tempfile")

	ecsParams, err := ReadECSParams(ecsParamsFileName)

	if assert.NoError(t, err) {
		taskDef := ecsParams.TaskDefinition

		containerDefs := taskDef.ContainerDefinitions
		assert.Equal(t, 3, len(containerDefs), "Expected 3 containers")

		mysql := containerDefs["mysql"]
		wordpress := containerDefs["wordpress"]
		logstash := containerDefs["logstash"]

		verifyHealthCheck(t, mysqlExpectedHealthCheck, mysql.HealthCheck.HealthCheck)
		verifyHealthCheck(t, wordpressExpectedHealthCheck, wordpress.HealthCheck.HealthCheck)
		verifyHealthCheck(t, logstashExpectedHealthCheck, logstash.HealthCheck.HealthCheck)

	}
}

func TestReadECSParams_WithHealthCheckErrorCaseInvalidDuration(t *testing.T) {
	ecsParamsString := `version: 1
task_definition:
  services:
    mysql:
      healthcheck:
        test: ["CMD", "curl", "-f", "http://localhost"]
        interval: cat
        timeout: 10s
        retries: 3
        start_period: 40s`

	content := []byte(ecsParamsString)

	tmpfile, err := ioutil.TempFile("", "ecs-params")
	assert.NoError(t, err, "Could not create ecs fields tempfile")

	ecsParamsFileName := tmpfile.Name()
	defer os.Remove(ecsParamsFileName)

	_, err = tmpfile.Write(content)
	assert.NoError(t, err, "Could not write data to ecs fields tempfile")

	err = tmpfile.Close()
	assert.NoError(t, err, "Could not close tempfile")

	_, err = ReadECSParams(ecsParamsFileName)

	assert.Error(t, err)
}

func TestReadECSParams_WithHealthCheckErrorCaseTestAndCommand(t *testing.T) {
	ecsParamsString := `version: 1
task_definition:
  services:
    mysql:
      healthcheck:
        test: ["CMD", "curl", "-f", "http://localhost"]
        command: ["CMD", "curl", "-f", "http://localhost"]
        interval: cat
        timeout: 10s
        retries: 3
        start_period: 40s`

	content := []byte(ecsParamsString)

	tmpfile, err := ioutil.TempFile("", "ecs-params")
	assert.NoError(t, err, "Could not create ecs fields tempfile")

	ecsParamsFileName := tmpfile.Name()
	defer os.Remove(ecsParamsFileName)

	_, err = tmpfile.Write(content)
	assert.NoError(t, err, "Could not write data to ecs fields tempfile")

	err = tmpfile.Close()
	assert.NoError(t, err, "Could not close tempfile")

	_, err = ReadECSParams(ecsParamsFileName)

	assert.Error(t, err)
}

func TestReadECSParams_WithHealthCheck_MissingIntFields(t *testing.T) {
	ecsParamsString := `version: 1
task_definition:
  services:
    mysql:
      healthcheck:
        test: ["CMD", "curl", "-f", "http://localhost"]`

	mysqlExpectedHealthCheck := ecs.HealthCheck{
		Command:     aws.StringSlice([]string{"CMD", "curl", "-f", "http://localhost"}),
		Interval:    nil,
		Timeout:     nil,
		Retries:     aws.Int64(int64(3)), // default value
		StartPeriod: nil,
	}

	content := []byte(ecsParamsString)

	tmpfile, err := ioutil.TempFile("", "ecs-params")
	assert.NoError(t, err, "Could not create ecs fields tempfile")

	ecsParamsFileName := tmpfile.Name()
	defer os.Remove(ecsParamsFileName)

	_, err = tmpfile.Write(content)
	assert.NoError(t, err, "Could not write data to ecs fields tempfile")

	err = tmpfile.Close()
	assert.NoError(t, err, "Could not close tempfile")

	ecsParams, err := ReadECSParams(ecsParamsFileName)

	if assert.NoError(t, err) {
		taskDef := ecsParams.TaskDefinition

		containerDefs := taskDef.ContainerDefinitions
		assert.Equal(t, 1, len(containerDefs), "Expected 3 containers")

		mysql := containerDefs["mysql"]

		verifyHealthCheck(t, mysqlExpectedHealthCheck, mysql.HealthCheck.HealthCheck)
	}
}

func TestReadECSParams_WithHealthCheck_MissingTest(t *testing.T) {
	ecsParamsString := `version: 1
task_definition:
  services:
    mysql:
      healthcheck:
        start_period: 10`

	content := []byte(ecsParamsString)

	tmpfile, err := ioutil.TempFile("", "ecs-params")
	assert.NoError(t, err, "Could not create ecs fields tempfile")

	ecsParamsFileName := tmpfile.Name()
	defer os.Remove(ecsParamsFileName)

	_, err = tmpfile.Write(content)
	assert.NoError(t, err, "Could not write data to ecs fields tempfile")

	err = tmpfile.Close()
	assert.NoError(t, err, "Could not close tempfile")

	ecsParams, err := ReadECSParams(ecsParamsFileName)

	if assert.NoError(t, err) {
		taskDef := ecsParams.TaskDefinition

		containerDefs := taskDef.ContainerDefinitions
		assert.Equal(t, 1, len(containerDefs), "Expected 3 containers")

		mysql := containerDefs["mysql"]

		assert.Len(t, mysql.HealthCheck.HealthCheck.Command, 0, "Expected health check command to be empty")
		assert.Nil(t, mysql.HealthCheck.HealthCheck.Interval, "Expected interval to be nil")
		assert.Nil(t, mysql.HealthCheck.HealthCheck.Timeout, "Expected timeout to be nil")
		assert.Equal(t, aws.Int64(3), mysql.HealthCheck.HealthCheck.Retries, "Expected retries to be default")
		assert.Equal(t, aws.Int64(10), mysql.HealthCheck.HealthCheck.StartPeriod)
	}
}

func verifyHealthCheck(t *testing.T, expected, actual ecs.HealthCheck) {
	assert.ElementsMatch(t, aws.StringValueSlice(expected.Command), aws.StringValueSlice(actual.Command), "Expected healthcheck command to match")
	assert.Equal(t, aws.Int64Value(expected.Interval), aws.Int64Value(actual.Interval), "Expected healthcheck interval to match")
	assert.Equal(t, aws.Int64Value(expected.Retries), aws.Int64Value(actual.Retries), "Expected healthcheck retries to match")
	assert.Equal(t, aws.Int64Value(expected.StartPeriod), aws.Int64Value(actual.StartPeriod), "Expected healthcheck start_period  to match")
	assert.Equal(t, aws.Int64Value(expected.Timeout), aws.Int64Value(actual.Timeout), "Expected healthcheck timeout to match")
}
