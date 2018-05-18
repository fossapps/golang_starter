// Code generated by MockGen. DO NOT EDIT.
// Source: golang_starter/db (interfaces: IUserManager)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	db "golang_starter/db"
	reflect "reflect"
)

// MockIUserManager is a mock of IUserManager interface
type MockIUserManager struct {
	ctrl     *gomock.Controller
	recorder *MockIUserManagerMockRecorder
}

// MockIUserManagerMockRecorder is the mock recorder for MockIUserManager
type MockIUserManagerMockRecorder struct {
	mock *MockIUserManager
}

// NewMockIUserManager creates a new mock instance
func NewMockIUserManager(ctrl *gomock.Controller) *MockIUserManager {
	mock := &MockIUserManager{ctrl: ctrl}
	mock.recorder = &MockIUserManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIUserManager) EXPECT() *MockIUserManagerMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockIUserManager) Create(arg0 db.User) error {
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockIUserManagerMockRecorder) Create(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIUserManager)(nil).Create), arg0)
}

// Edit mocks base method
func (m *MockIUserManager) Edit(arg0 string, arg1 db.User) error {
	ret := m.ctrl.Call(m, "Edit", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Edit indicates an expected call of Edit
func (mr *MockIUserManagerMockRecorder) Edit(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Edit", reflect.TypeOf((*MockIUserManager)(nil).Edit), arg0, arg1)
}

// FindByEmail mocks base method
func (m *MockIUserManager) FindByEmail(arg0 string) *db.User {
	ret := m.ctrl.Call(m, "FindByEmail", arg0)
	ret0, _ := ret[0].(*db.User)
	return ret0
}

// FindByEmail indicates an expected call of FindByEmail
func (mr *MockIUserManagerMockRecorder) FindByEmail(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockIUserManager)(nil).FindByEmail), arg0)
}

// FindById mocks base method
func (m *MockIUserManager) FindById(arg0 string) *db.User {
	ret := m.ctrl.Call(m, "FindById", arg0)
	ret0, _ := ret[0].(*db.User)
	return ret0
}

// FindById indicates an expected call of FindById
func (mr *MockIUserManagerMockRecorder) FindById(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockIUserManager)(nil).FindById), arg0)
}

// List mocks base method
func (m *MockIUserManager) List() ([]db.User, error) {
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockIUserManagerMockRecorder) List() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIUserManager)(nil).List))
}
