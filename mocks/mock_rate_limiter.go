// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/fossapps/starter/adapters (interfaces: IRateLimiter)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockIRateLimiter is a mock of IRateLimiter interface
type MockIRateLimiter struct {
	ctrl     *gomock.Controller
	recorder *MockIRateLimiterMockRecorder
}

// MockIRateLimiterMockRecorder is the mock recorder for MockIRateLimiter
type MockIRateLimiterMockRecorder struct {
	mock *MockIRateLimiter
}

// NewMockIRateLimiter creates a new mock instance
func NewMockIRateLimiter(ctrl *gomock.Controller) *MockIRateLimiter {
	mock := &MockIRateLimiter{ctrl: ctrl}
	mock.recorder = &MockIRateLimiterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIRateLimiter) EXPECT() *MockIRateLimiterMockRecorder {
	return m.recorder
}

// Count mocks base method
func (m *MockIRateLimiter) Count(arg0 string) (int64, error) {
	ret := m.ctrl.Call(m, "Count", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count
func (mr *MockIRateLimiterMockRecorder) Count(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockIRateLimiter)(nil).Count), arg0)
}

// Hit mocks base method
func (m *MockIRateLimiter) Hit(arg0 string) (int64, error) {
	ret := m.ctrl.Call(m, "Hit", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Hit indicates an expected call of Hit
func (mr *MockIRateLimiterMockRecorder) Hit(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hit", reflect.TypeOf((*MockIRateLimiter)(nil).Hit), arg0)
}
