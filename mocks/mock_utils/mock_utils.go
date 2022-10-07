// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/netapp/trident/utils (interfaces: LUKSDeviceInterface)

// Package mock_utils is a generated GoMock package.
package mock_utils

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockLUKSDeviceInterface is a mock of LUKSDeviceInterface interface.
type MockLUKSDeviceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockLUKSDeviceInterfaceMockRecorder
}

// MockLUKSDeviceInterfaceMockRecorder is the mock recorder for MockLUKSDeviceInterface.
type MockLUKSDeviceInterfaceMockRecorder struct {
	mock *MockLUKSDeviceInterface
}

// NewMockLUKSDeviceInterface creates a new mock instance.
func NewMockLUKSDeviceInterface(ctrl *gomock.Controller) *MockLUKSDeviceInterface {
	mock := &MockLUKSDeviceInterface{ctrl: ctrl}
	mock.recorder = &MockLUKSDeviceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLUKSDeviceInterface) EXPECT() *MockLUKSDeviceInterfaceMockRecorder {
	return m.recorder
}

// DevicePath mocks base method.
func (m *MockLUKSDeviceInterface) DevicePath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DevicePath")
	ret0, _ := ret[0].(string)
	return ret0
}

// DevicePath indicates an expected call of DevicePath.
func (mr *MockLUKSDeviceInterfaceMockRecorder) DevicePath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DevicePath", reflect.TypeOf((*MockLUKSDeviceInterface)(nil).DevicePath))
}

// IsLUKSFormatted mocks base method.
func (m *MockLUKSDeviceInterface) IsLUKSFormatted(arg0 context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsLUKSFormatted", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsLUKSFormatted indicates an expected call of IsLUKSFormatted.
func (mr *MockLUKSDeviceInterfaceMockRecorder) IsLUKSFormatted(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsLUKSFormatted", reflect.TypeOf((*MockLUKSDeviceInterface)(nil).IsLUKSFormatted), arg0)
}

// IsOpen mocks base method.
func (m *MockLUKSDeviceInterface) IsOpen(arg0 context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsOpen", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsOpen indicates an expected call of IsOpen.
func (mr *MockLUKSDeviceInterfaceMockRecorder) IsOpen(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsOpen", reflect.TypeOf((*MockLUKSDeviceInterface)(nil).IsOpen), arg0)
}

// LUKSDeviceName mocks base method.
func (m *MockLUKSDeviceInterface) LUKSDeviceName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LUKSDeviceName")
	ret0, _ := ret[0].(string)
	return ret0
}

// LUKSDeviceName indicates an expected call of LUKSDeviceName.
func (mr *MockLUKSDeviceInterfaceMockRecorder) LUKSDeviceName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LUKSDeviceName", reflect.TypeOf((*MockLUKSDeviceInterface)(nil).LUKSDeviceName))
}

// LUKSDevicePath mocks base method.
func (m *MockLUKSDeviceInterface) LUKSDevicePath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LUKSDevicePath")
	ret0, _ := ret[0].(string)
	return ret0
}

// LUKSDevicePath indicates an expected call of LUKSDevicePath.
func (mr *MockLUKSDeviceInterfaceMockRecorder) LUKSDevicePath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LUKSDevicePath", reflect.TypeOf((*MockLUKSDeviceInterface)(nil).LUKSDevicePath))
}

// LUKSFormat mocks base method.
func (m *MockLUKSDeviceInterface) LUKSFormat(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LUKSFormat", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// LUKSFormat indicates an expected call of LUKSFormat.
func (mr *MockLUKSDeviceInterfaceMockRecorder) LUKSFormat(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LUKSFormat", reflect.TypeOf((*MockLUKSDeviceInterface)(nil).LUKSFormat), arg0, arg1)
}

// Open mocks base method.
func (m *MockLUKSDeviceInterface) Open(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Open indicates an expected call of Open.
func (mr *MockLUKSDeviceInterfaceMockRecorder) Open(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockLUKSDeviceInterface)(nil).Open), arg0, arg1)
}
