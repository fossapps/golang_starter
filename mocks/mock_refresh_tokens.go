// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/fossapps/starter/db (interfaces: IRefreshTokenManager)

// Package mocks is a generated GoMock package.
package mocks

import (
	db "github.com/fossapps/starter/db"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockIRefreshTokenManager is a mock of IRefreshTokenManager interface
type MockIRefreshTokenManager struct {
	ctrl     *gomock.Controller
	recorder *MockIRefreshTokenManagerMockRecorder
}

// MockIRefreshTokenManagerMockRecorder is the mock recorder for MockIRefreshTokenManager
type MockIRefreshTokenManagerMockRecorder struct {
	mock *MockIRefreshTokenManager
}

// NewMockIRefreshTokenManager creates a new mock instance
func NewMockIRefreshTokenManager(ctrl *gomock.Controller) *MockIRefreshTokenManager {
	mock := &MockIRefreshTokenManager{ctrl: ctrl}
	mock.recorder = &MockIRefreshTokenManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIRefreshTokenManager) EXPECT() *MockIRefreshTokenManagerMockRecorder {
	return m.recorder
}

// Add mocks base method
func (m *MockIRefreshTokenManager) Add(arg0, arg1 string) {
	m.ctrl.Call(m, "Add", arg0, arg1)
}

// Add indicates an expected call of Add
func (mr *MockIRefreshTokenManagerMockRecorder) Add(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockIRefreshTokenManager)(nil).Add), arg0, arg1)
}

// Delete mocks base method
func (m *MockIRefreshTokenManager) Delete(arg0 string) error {
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockIRefreshTokenManagerMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIRefreshTokenManager)(nil).Delete), arg0)
}

// FindOne mocks base method
func (m *MockIRefreshTokenManager) FindOne(arg0 string) *db.RefreshToken {
	ret := m.ctrl.Call(m, "FindOne", arg0)
	ret0, _ := ret[0].(*db.RefreshToken)
	return ret0
}

// FindOne indicates an expected call of FindOne
func (mr *MockIRefreshTokenManagerMockRecorder) FindOne(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOne", reflect.TypeOf((*MockIRefreshTokenManager)(nil).FindOne), arg0)
}

// List mocks base method
func (m *MockIRefreshTokenManager) List(arg0 string) ([]db.RefreshToken, error) {
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].([]db.RefreshToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockIRefreshTokenManagerMockRecorder) List(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIRefreshTokenManager)(nil).List), arg0)
}
