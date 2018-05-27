// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/fossapps/starter/db (interfaces: DeviceManager)

// Package mock is a generated GoMock package.
package mock

import (
	db "github.com/fossapps/starter/db"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDeviceManager is a mock of DeviceManager interface
type MockDeviceManager struct {
	ctrl     *gomock.Controller
	recorder *MockDeviceManagerMockRecorder
}

// MockDeviceManagerMockRecorder is the mock recorder for MockDeviceManager
type MockDeviceManagerMockRecorder struct {
	mock *MockDeviceManager
}

// NewMockDeviceManager creates a new mock instance
func NewMockDeviceManager(ctrl *gomock.Controller) *MockDeviceManager {
	mock := &MockDeviceManager{ctrl: ctrl}
	mock.recorder = &MockDeviceManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDeviceManager) EXPECT() *MockDeviceManagerMockRecorder {
	return m.recorder
}

// Exists mocks base method
func (m *MockDeviceManager) Exists(arg0 string) (bool, error) {
	ret := m.ctrl.Call(m, "Exists", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists
func (mr *MockDeviceManagerMockRecorder) Exists(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockDeviceManager)(nil).Exists), arg0)
}

// FindByToken mocks base method
func (m *MockDeviceManager) FindByToken(arg0 string) (*db.Device, error) {
	ret := m.ctrl.Call(m, "FindByToken", arg0)
	ret0, _ := ret[0].(*db.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByToken indicates an expected call of FindByToken
func (mr *MockDeviceManagerMockRecorder) FindByToken(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByToken", reflect.TypeOf((*MockDeviceManager)(nil).FindByToken), arg0)
}

// Register mocks base method
func (m *MockDeviceManager) Register(arg0 string) error {
	ret := m.ctrl.Call(m, "Register", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register
func (mr *MockDeviceManagerMockRecorder) Register(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockDeviceManager)(nil).Register), arg0)
}
