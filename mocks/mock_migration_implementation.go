// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/fossapps/starter/migrations (interfaces: Migration)

// Package mocks is a generated GoMock package.
package mocks

import (
	db "github.com/fossapps/starter/db"
	mgo "github.com/globalsign/mgo"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockMigration is a mock of Migration interface
type MockMigration struct {
	ctrl     *gomock.Controller
	recorder *MockMigrationMockRecorder
}

// MockMigrationMockRecorder is the mock recorder for MockMigration
type MockMigrationMockRecorder struct {
	mock *MockMigration
}

// NewMockMigration creates a new mock instance
func NewMockMigration(ctrl *gomock.Controller) *MockMigration {
	mock := &MockMigration{ctrl: ctrl}
	mock.recorder = &MockMigrationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMigration) EXPECT() *MockMigrationMockRecorder {
	return m.recorder
}

// Apply mocks base method
func (m *MockMigration) Apply(arg0 db.DB) {
	m.ctrl.Call(m, "Apply", arg0)
}

// Apply indicates an expected call of Apply
func (mr *MockMigrationMockRecorder) Apply(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Apply", reflect.TypeOf((*MockMigration)(nil).Apply), arg0)
}

// GetDescription mocks base method
func (m *MockMigration) GetDescription() string {
	ret := m.ctrl.Call(m, "GetDescription")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetDescription indicates an expected call of GetDescription
func (mr *MockMigrationMockRecorder) GetDescription() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDescription", reflect.TypeOf((*MockMigration)(nil).GetDescription))
}

// GetKey mocks base method
func (m *MockMigration) GetKey() string {
	ret := m.ctrl.Call(m, "GetKey")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetKey indicates an expected call of GetKey
func (mr *MockMigrationMockRecorder) GetKey() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKey", reflect.TypeOf((*MockMigration)(nil).GetKey))
}

// Remove mocks base method
func (m *MockMigration) Remove(arg0 *mgo.Database) {
	m.ctrl.Call(m, "Remove", arg0)
}

// Remove indicates an expected call of Remove
func (mr *MockMigrationMockRecorder) Remove(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockMigration)(nil).Remove), arg0)
}
