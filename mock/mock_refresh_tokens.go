// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/fossapps/starter/db (interfaces: RefreshTokenManager)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	db "github.com/fossapps/starter/db"
	gomock "github.com/golang/mock/gomock"
)

// MockRefreshTokenManager is a mock of RefreshTokenManager interface
type MockRefreshTokenManager struct {
	ctrl     *gomock.Controller
	recorder *MockRefreshTokenManagerMockRecorder
}

// MockRefreshTokenManagerMockRecorder is the mock recorder for MockRefreshTokenManager
type MockRefreshTokenManagerMockRecorder struct {
	mock *MockRefreshTokenManager
}

// NewMockRefreshTokenManager creates a new mock instance
func NewMockRefreshTokenManager(ctrl *gomock.Controller) *MockRefreshTokenManager {
	mock := &MockRefreshTokenManager{ctrl: ctrl}
	mock.recorder = &MockRefreshTokenManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRefreshTokenManager) EXPECT() *MockRefreshTokenManagerMockRecorder {
	return m.recorder
}

// Add mocks base method
func (m *MockRefreshTokenManager) Add(arg0, arg1 string) error {
	ret := m.ctrl.Call(m, "Add", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add
func (mr *MockRefreshTokenManagerMockRecorder) Add(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockRefreshTokenManager)(nil).Add), arg0, arg1)
}

// Delete mocks base method
func (m *MockRefreshTokenManager) Delete(arg0 string) error {
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockRefreshTokenManagerMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRefreshTokenManager)(nil).Delete), arg0)
}

// FindOne mocks base method
func (m *MockRefreshTokenManager) FindOne(arg0 string) (*db.RefreshToken, error) {
	ret := m.ctrl.Call(m, "FindOne", arg0)
	ret0, _ := ret[0].(*db.RefreshToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOne indicates an expected call of FindOne
func (mr *MockRefreshTokenManagerMockRecorder) FindOne(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOne", reflect.TypeOf((*MockRefreshTokenManager)(nil).FindOne), arg0)
}

// List mocks base method
func (m *MockRefreshTokenManager) List(arg0 string) ([]db.RefreshToken, error) {
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].([]db.RefreshToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockRefreshTokenManagerMockRecorder) List(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRefreshTokenManager)(nil).List), arg0)
}
