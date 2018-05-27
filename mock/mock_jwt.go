// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/fossapps/starter/jwt (interfaces: Manager)

// Package mock is a generated GoMock package.
package mock

import (
	http "net/http"
	reflect "reflect"

	db "github.com/fossapps/starter/db"
	jwt "github.com/fossapps/starter/jwt"
	gomock "github.com/golang/mock/gomock"
)

// MockJwtManager is a mock of Manager interface
type MockJwtManager struct {
	ctrl     *gomock.Controller
	recorder *MockJwtManagerMockRecorder
}

// MockJwtManagerMockRecorder is the mock recorder for MockJwtManager
type MockJwtManagerMockRecorder struct {
	mock *MockJwtManager
}

// NewMockJwtManager creates a new mock instance
func NewMockJwtManager(ctrl *gomock.Controller) *MockJwtManager {
	mock := &MockJwtManager{ctrl: ctrl}
	mock.recorder = &MockJwtManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockJwtManager) EXPECT() *MockJwtManagerMockRecorder {
	return m.recorder
}

// CreateForUser mocks base method
func (m *MockJwtManager) CreateForUser(arg0 *db.User) (string, error) {
	ret := m.ctrl.Call(m, "CreateForUser", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateForUser indicates an expected call of CreateForUser
func (mr *MockJwtManagerMockRecorder) CreateForUser(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateForUser", reflect.TypeOf((*MockJwtManager)(nil).CreateForUser), arg0)
}

// GetJwtDataFromRequest mocks base method
func (m *MockJwtManager) GetJwtDataFromRequest(arg0 *http.Request) (*jwt.Claims, error) {
	ret := m.ctrl.Call(m, "GetJwtDataFromRequest", arg0)
	ret0, _ := ret[0].(*jwt.Claims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJwtDataFromRequest indicates an expected call of GetJwtDataFromRequest
func (mr *MockJwtManagerMockRecorder) GetJwtDataFromRequest(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJwtDataFromRequest", reflect.TypeOf((*MockJwtManager)(nil).GetJwtDataFromRequest), arg0)
}
