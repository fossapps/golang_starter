// Code generated by MockGen. DO NOT EDIT.
// Source: crazy_nl_backend/db (interfaces: IDeviceManager)

// Package mocks is a generated GoMock package.
package mocks

import (
	db "crazy_nl_backend/db"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockIDeviceManager is a mock of IDeviceManager interface
type MockIDeviceManager struct {
	ctrl     *gomock.Controller
	recorder *MockIDeviceManagerMockRecorder
}

// MockIDeviceManagerMockRecorder is the mock recorder for MockIDeviceManager
type MockIDeviceManagerMockRecorder struct {
	mock *MockIDeviceManager
}

// NewMockIDeviceManager creates a new mock instance
func NewMockIDeviceManager(ctrl *gomock.Controller) *MockIDeviceManager {
	mock := &MockIDeviceManager{ctrl: ctrl}
	mock.recorder = &MockIDeviceManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIDeviceManager) EXPECT() *MockIDeviceManagerMockRecorder {
	return m.recorder
}

// Exists mocks base method
func (m *MockIDeviceManager) Exists(arg0 string) bool {
	ret := m.ctrl.Call(m, "Exists", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Exists indicates an expected call of Exists
func (mr *MockIDeviceManagerMockRecorder) Exists(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockIDeviceManager)(nil).Exists), arg0)
}

// FindByToken mocks base method
func (m *MockIDeviceManager) FindByToken(arg0 string) *db.Device {
	ret := m.ctrl.Call(m, "FindByToken", arg0)
	ret0, _ := ret[0].(*db.Device)
	return ret0
}

// FindByToken indicates an expected call of FindByToken
func (mr *MockIDeviceManagerMockRecorder) FindByToken(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByToken", reflect.TypeOf((*MockIDeviceManager)(nil).FindByToken), arg0)
}

// Register mocks base method
func (m *MockIDeviceManager) Register(arg0 string) error {
	ret := m.ctrl.Call(m, "Register", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register
func (mr *MockIDeviceManagerMockRecorder) Register(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockIDeviceManager)(nil).Register), arg0)
}
