// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/fossapps/starter/db (interfaces: MigrationManager)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMigrationManager is a mock of MigrationManager interface
type MockMigrationManager struct {
	ctrl     *gomock.Controller
	recorder *MockMigrationManagerMockRecorder
}

// MockMigrationManagerMockRecorder is the mock recorder for MockMigrationManager
type MockMigrationManagerMockRecorder struct {
	mock *MockMigrationManager
}

// NewMockMigrationManager creates a new mock instance
func NewMockMigrationManager(ctrl *gomock.Controller) *MockMigrationManager {
	mock := &MockMigrationManager{ctrl: ctrl}
	mock.recorder = &MockMigrationManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMigrationManager) EXPECT() *MockMigrationManagerMockRecorder {
	return m.recorder
}

// IsApplied mocks base method
func (m *MockMigrationManager) IsApplied(arg0 string) (bool, error) {
	ret := m.ctrl.Call(m, "IsApplied", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsApplied indicates an expected call of IsApplied
func (mr *MockMigrationManagerMockRecorder) IsApplied(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsApplied", reflect.TypeOf((*MockMigrationManager)(nil).IsApplied), arg0)
}

// MarkApplied mocks base method
func (m *MockMigrationManager) MarkApplied(arg0, arg1 string) error {
	ret := m.ctrl.Call(m, "MarkApplied", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkApplied indicates an expected call of MarkApplied
func (mr *MockMigrationManagerMockRecorder) MarkApplied(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkApplied", reflect.TypeOf((*MockMigrationManager)(nil).MarkApplied), arg0, arg1)
}