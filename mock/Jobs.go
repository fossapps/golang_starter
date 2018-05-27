// Code generated by MockGen. DO NOT EDIT.
// Source: starter/worker (interfaces: ICronJob)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	captain "github.com/cyberhck/captain"
	gomock "github.com/golang/mock/gomock"
)

// MockICronJob is a mock of ICronJob interface
type MockICronJob struct {
	ctrl     *gomock.Controller
	recorder *MockICronJobMockRecorder
}

// MockICronJobMockRecorder is the mock recorder for MockICronJob
type MockICronJobMockRecorder struct {
	mock *MockICronJob
}

// NewMockICronJob creates a new mock instance
func NewMockICronJob(ctrl *gomock.Controller) *MockICronJob {
	mock := &MockICronJob{ctrl: ctrl}
	mock.recorder = &MockICronJobMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockICronJob) EXPECT() *MockICronJobMockRecorder {
	return m.recorder
}

// Job mock base method
func (m *MockICronJob) Job() captain.Worker {
	ret := m.ctrl.Call(m, "Job")
	ret0, _ := ret[0].(captain.Worker)
	return ret0
}

// Job indicates an expected call of Job
func (mr *MockICronJobMockRecorder) Job() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Job", reflect.TypeOf((*MockICronJob)(nil).Job))
}

// LockProvider mock base method
func (m *MockICronJob) LockProvider() captain.LockProvider {
	ret := m.ctrl.Call(m, "LockProvider")
	ret0, _ := ret[0].(captain.LockProvider)
	return ret0
}

// LockProvider indicates an expected call of LockProvider
func (mr *MockICronJobMockRecorder) LockProvider() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LockProvider", reflect.TypeOf((*MockICronJob)(nil).LockProvider))
}

// ResultProcessor mock base method
func (m *MockICronJob) ResultProcessor() captain.ResultProcessor {
	ret := m.ctrl.Call(m, "ResultProcessor")
	ret0, _ := ret[0].(captain.ResultProcessor)
	return ret0
}

// ResultProcessor indicates an expected call of ResultProcessor
func (mr *MockICronJobMockRecorder) ResultProcessor() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResultProcessor", reflect.TypeOf((*MockICronJob)(nil).ResultProcessor))
}

// RuntimeProcessor mock base method
func (m *MockICronJob) RuntimeProcessor() captain.RuntimeProcessor {
	ret := m.ctrl.Call(m, "RuntimeProcessor")
	ret0, _ := ret[0].(captain.RuntimeProcessor)
	return ret0
}

// RuntimeProcessor indicates an expected call of RuntimeProcessor
func (mr *MockICronJobMockRecorder) RuntimeProcessor() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RuntimeProcessor", reflect.TypeOf((*MockICronJob)(nil).RuntimeProcessor))
}

// IsApplied mock base method
func (m *MockICronJob) ShouldRun(arg0 string) bool {
	ret := m.ctrl.Call(m, "IsApplied", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsApplied indicates an expected call of IsApplied
func (mr *MockICronJobMockRecorder) ShouldRun(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsApplied", reflect.TypeOf((*MockICronJob)(nil).ShouldRun), arg0)
}
