// Copyright 2015-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/aws/amazon-ecs-cli/ecs-cli/modules/clients/aws/cloudformation (interfaces: CloudformationClient)

package mock_cloudformation

import (
	cloudformation "github.com/aws/amazon-ecs-cli/ecs-cli/modules/clients/aws/cloudformation"
	cloudformation0 "github.com/aws/aws-sdk-go/service/cloudformation"
	gomock "github.com/golang/mock/gomock"
)

// Mock of CloudformationClient interface
type MockCloudformationClient struct {
	ctrl     *gomock.Controller
	recorder *_MockCloudformationClientRecorder
}

// Recorder for MockCloudformationClient (not exported)
type _MockCloudformationClientRecorder struct {
	mock *MockCloudformationClient
}

func NewMockCloudformationClient(ctrl *gomock.Controller) *MockCloudformationClient {
	mock := &MockCloudformationClient{ctrl: ctrl}
	mock.recorder = &_MockCloudformationClientRecorder{mock}
	return mock
}

func (_m *MockCloudformationClient) EXPECT() *_MockCloudformationClientRecorder {
	return _m.recorder
}

func (_m *MockCloudformationClient) CreateStack(_param0 string, _param1 string, _param2 bool, _param3 *cloudformation.CfnStackParams) (string, error) {
	ret := _m.ctrl.Call(_m, "CreateStack", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCloudformationClientRecorder) CreateStack(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateStack", arg0, arg1, arg2, arg3)
}

func (_m *MockCloudformationClient) DeleteStack(_param0 string) error {
	ret := _m.ctrl.Call(_m, "DeleteStack", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCloudformationClientRecorder) DeleteStack(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteStack", arg0)
}

func (_m *MockCloudformationClient) DescribeNetworkResources(_param0 string) error {
	ret := _m.ctrl.Call(_m, "DescribeNetworkResources", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCloudformationClientRecorder) DescribeNetworkResources(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DescribeNetworkResources", arg0)
}

func (_m *MockCloudformationClient) DescribeStacks(_param0 string) (*cloudformation0.DescribeStacksOutput, error) {
	ret := _m.ctrl.Call(_m, "DescribeStacks", _param0)
	ret0, _ := ret[0].(*cloudformation0.DescribeStacksOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCloudformationClientRecorder) DescribeStacks(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DescribeStacks", arg0)
}

func (_m *MockCloudformationClient) GetStackParameters(_param0 string) ([]*cloudformation0.Parameter, error) {
	ret := _m.ctrl.Call(_m, "GetStackParameters", _param0)
	ret0, _ := ret[0].([]*cloudformation0.Parameter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCloudformationClientRecorder) GetStackParameters(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetStackParameters", arg0)
}

func (_m *MockCloudformationClient) UpdateStack(_param0 string, _param1 *cloudformation.CfnStackParams) (string, error) {
	ret := _m.ctrl.Call(_m, "UpdateStack", _param0, _param1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCloudformationClientRecorder) UpdateStack(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateStack", arg0, arg1)
}

func (_m *MockCloudformationClient) ValidateStackExists(_param0 string) error {
	ret := _m.ctrl.Call(_m, "ValidateStackExists", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCloudformationClientRecorder) ValidateStackExists(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ValidateStackExists", arg0)
}

func (_m *MockCloudformationClient) WaitUntilCreateComplete(_param0 string) error {
	ret := _m.ctrl.Call(_m, "WaitUntilCreateComplete", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCloudformationClientRecorder) WaitUntilCreateComplete(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "WaitUntilCreateComplete", arg0)
}

func (_m *MockCloudformationClient) WaitUntilDeleteComplete(_param0 string) error {
	ret := _m.ctrl.Call(_m, "WaitUntilDeleteComplete", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCloudformationClientRecorder) WaitUntilDeleteComplete(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "WaitUntilDeleteComplete", arg0)
}

func (_m *MockCloudformationClient) WaitUntilUpdateComplete(_param0 string) error {
	ret := _m.ctrl.Call(_m, "WaitUntilUpdateComplete", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCloudformationClientRecorder) WaitUntilUpdateComplete(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "WaitUntilUpdateComplete", arg0)
}
