// Code generated by MockGen. DO NOT EDIT.
// Source: crazy_nl_backend/helpers (interfaces: IRedisClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	redis "github.com/go-redis/redis"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockIRedisClient is a mock of IRedisClient interface
type MockIRedisClient struct {
	ctrl     *gomock.Controller
	recorder *MockIRedisClientMockRecorder
}

// MockIRedisClientMockRecorder is the mock recorder for MockIRedisClient
type MockIRedisClientMockRecorder struct {
	mock *MockIRedisClient
}

// NewMockIRedisClient creates a new mock instance
func NewMockIRedisClient(ctrl *gomock.Controller) *MockIRedisClient {
	mock := &MockIRedisClient{ctrl: ctrl}
	mock.recorder = &MockIRedisClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIRedisClient) EXPECT() *MockIRedisClientMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockIRedisClient) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockIRedisClientMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockIRedisClient)(nil).Close))
}

// Expire mocks base method
func (m *MockIRedisClient) Expire(arg0 string, arg1 time.Duration) (bool, error) {
	ret := m.ctrl.Call(m, "Expire", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Expire indicates an expected call of Expire
func (mr *MockIRedisClientMockRecorder) Expire(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Expire", reflect.TypeOf((*MockIRedisClient)(nil).Expire), arg0, arg1)
}

// SAdd mocks base method
func (m *MockIRedisClient) SAdd(arg0 string, arg1 ...interface{}) (int64, error) {
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SAdd", varargs...)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SAdd indicates an expected call of SAdd
func (mr *MockIRedisClientMockRecorder) SAdd(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SAdd", reflect.TypeOf((*MockIRedisClient)(nil).SAdd), varargs...)
}

// SIsMember mocks base method
func (m *MockIRedisClient) SIsMember(arg0 string, arg1 interface{}) (bool, error) {
	ret := m.ctrl.Call(m, "SIsMember", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SIsMember indicates an expected call of SIsMember
func (mr *MockIRedisClientMockRecorder) SIsMember(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SIsMember", reflect.TypeOf((*MockIRedisClient)(nil).SIsMember), arg0, arg1)
}

// SMembers mocks base method
func (m *MockIRedisClient) SMembers(arg0 string) ([]string, error) {
	ret := m.ctrl.Call(m, "SMembers", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SMembers indicates an expected call of SMembers
func (mr *MockIRedisClientMockRecorder) SMembers(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SMembers", reflect.TypeOf((*MockIRedisClient)(nil).SMembers), arg0)
}

// SPop mocks base method
func (m *MockIRedisClient) SPop(arg0 string) (string, error) {
	ret := m.ctrl.Call(m, "SPop", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SPop indicates an expected call of SPop
func (mr *MockIRedisClientMockRecorder) SPop(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SPop", reflect.TypeOf((*MockIRedisClient)(nil).SPop), arg0)
}

// SRem mocks base method
func (m *MockIRedisClient) SRem(arg0 string, arg1 ...interface{}) (int64, error) {
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SRem", varargs...)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SRem indicates an expected call of SRem
func (mr *MockIRedisClientMockRecorder) SRem(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SRem", reflect.TypeOf((*MockIRedisClient)(nil).SRem), varargs...)
}

// ZAdd mocks base method
func (m *MockIRedisClient) ZAdd(arg0 string, arg1 ...redis.Z) (int64, error) {
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ZAdd", varargs...)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ZAdd indicates an expected call of ZAdd
func (mr *MockIRedisClientMockRecorder) ZAdd(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ZAdd", reflect.TypeOf((*MockIRedisClient)(nil).ZAdd), varargs...)
}

// ZCard mocks base method
func (m *MockIRedisClient) ZCard(arg0 string) (int64, error) {
	ret := m.ctrl.Call(m, "ZCard", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ZCard indicates an expected call of ZCard
func (mr *MockIRedisClientMockRecorder) ZCard(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ZCard", reflect.TypeOf((*MockIRedisClient)(nil).ZCard), arg0)
}

// ZRemRangeByScore mocks base method
func (m *MockIRedisClient) ZRemRangeByScore(arg0, arg1, arg2 string) (int64, error) {
	ret := m.ctrl.Call(m, "ZRemRangeByScore", arg0, arg1, arg2)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ZRemRangeByScore indicates an expected call of ZRemRangeByScore
func (mr *MockIRedisClientMockRecorder) ZRemRangeByScore(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ZRemRangeByScore", reflect.TypeOf((*MockIRedisClient)(nil).ZRemRangeByScore), arg0, arg1, arg2)
}
