// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cyberhck/pushy (interfaces: IPushyClient)

// Package mock is a generated GoMock package.
package mock

import (
	pushy "github.com/cyberhck/pushy"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockIPushyClient is a mock of IPushyClient interface
type MockIPushyClient struct {
	ctrl     *gomock.Controller
	recorder *MockIPushyClientMockRecorder
}

// MockIPushyClientMockRecorder is the mock recorder for MockIPushyClient
type MockIPushyClientMockRecorder struct {
	mock *MockIPushyClient
}

// NewMockIPushyClient creates a new mock instance
func NewMockIPushyClient(ctrl *gomock.Controller) *MockIPushyClient {
	mock := &MockIPushyClient{ctrl: ctrl}
	mock.recorder = &MockIPushyClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIPushyClient) EXPECT() *MockIPushyClientMockRecorder {
	return m.recorder
}

// DeleteNotification mocks base method
func (m *MockIPushyClient) DeleteNotification(arg0 string) (*pushy.SimpleSuccess, *pushy.Error, error) {
	ret := m.ctrl.Call(m, "DeleteNotification", arg0)
	ret0, _ := ret[0].(*pushy.SimpleSuccess)
	ret1, _ := ret[1].(*pushy.Error)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DeleteNotification indicates an expected call of DeleteNotification
func (mr *MockIPushyClientMockRecorder) DeleteNotification(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNotification", reflect.TypeOf((*MockIPushyClient)(nil).DeleteNotification), arg0)
}

// DeviceInfo mocks base method
func (m *MockIPushyClient) DeviceInfo(arg0 string) (*pushy.DeviceInfo, *pushy.Error, error) {
	ret := m.ctrl.Call(m, "DeviceInfo", arg0)
	ret0, _ := ret[0].(*pushy.DeviceInfo)
	ret1, _ := ret[1].(*pushy.Error)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DeviceInfo indicates an expected call of DeviceInfo
func (mr *MockIPushyClientMockRecorder) DeviceInfo(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeviceInfo", reflect.TypeOf((*MockIPushyClient)(nil).DeviceInfo), arg0)
}

// DevicePresence mocks base method
func (m *MockIPushyClient) DevicePresence(arg0 ...string) (*pushy.DevicePresenceResponse, *pushy.Error, error) {
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DevicePresence", varargs...)
	ret0, _ := ret[0].(*pushy.DevicePresenceResponse)
	ret1, _ := ret[1].(*pushy.Error)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DevicePresence indicates an expected call of DevicePresence
func (mr *MockIPushyClientMockRecorder) DevicePresence(arg0 ...interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DevicePresence", reflect.TypeOf((*MockIPushyClient)(nil).DevicePresence), arg0...)
}

// GetHTTPClient mocks base method
func (m *MockIPushyClient) GetHTTPClient() pushy.IHTTPClient {
	ret := m.ctrl.Call(m, "GetHTTPClient")
	ret0, _ := ret[0].(pushy.IHTTPClient)
	return ret0
}

// GetHTTPClient indicates an expected call of GetHTTPClient
func (mr *MockIPushyClientMockRecorder) GetHTTPClient() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHTTPClient", reflect.TypeOf((*MockIPushyClient)(nil).GetHTTPClient))
}

// NotificationStatus mocks base method
func (m *MockIPushyClient) NotificationStatus(arg0 string) (*pushy.NotificationStatus, *pushy.Error, error) {
	ret := m.ctrl.Call(m, "NotificationStatus", arg0)
	ret0, _ := ret[0].(*pushy.NotificationStatus)
	ret1, _ := ret[1].(*pushy.Error)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// NotificationStatus indicates an expected call of NotificationStatus
func (mr *MockIPushyClientMockRecorder) NotificationStatus(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotificationStatus", reflect.TypeOf((*MockIPushyClient)(nil).NotificationStatus), arg0)
}

// NotifyDevice mocks base method
func (m *MockIPushyClient) NotifyDevice(arg0 pushy.SendNotificationRequest) (*pushy.NotificationResponse, *pushy.Error, error) {
	ret := m.ctrl.Call(m, "NotifyDevice", arg0)
	ret0, _ := ret[0].(*pushy.NotificationResponse)
	ret1, _ := ret[1].(*pushy.Error)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// NotifyDevice indicates an expected call of NotifyDevice
func (mr *MockIPushyClientMockRecorder) NotifyDevice(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyDevice", reflect.TypeOf((*MockIPushyClient)(nil).NotifyDevice), arg0)
}

// SetHTTPClient mocks base method
func (m *MockIPushyClient) SetHTTPClient(arg0 pushy.IHTTPClient) {
	m.ctrl.Call(m, "SetHTTPClient", arg0)
}

// SetHTTPClient indicates an expected call of SetHTTPClient
func (mr *MockIPushyClientMockRecorder) SetHTTPClient(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHTTPClient", reflect.TypeOf((*MockIPushyClient)(nil).SetHTTPClient), arg0)
}

// SubscribeToTopic mocks base method
func (m *MockIPushyClient) SubscribeToTopic(arg0 string, arg1 ...string) (*pushy.SimpleSuccess, *pushy.Error, error) {
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SubscribeToTopic", varargs...)
	ret0, _ := ret[0].(*pushy.SimpleSuccess)
	ret1, _ := ret[1].(*pushy.Error)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SubscribeToTopic indicates an expected call of SubscribeToTopic
func (mr *MockIPushyClientMockRecorder) SubscribeToTopic(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeToTopic", reflect.TypeOf((*MockIPushyClient)(nil).SubscribeToTopic), varargs...)
}

// UnsubscribeFromTopic mocks base method
func (m *MockIPushyClient) UnsubscribeFromTopic(arg0 string, arg1 ...string) (*pushy.SimpleSuccess, *pushy.Error, error) {
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UnsubscribeFromTopic", varargs...)
	ret0, _ := ret[0].(*pushy.SimpleSuccess)
	ret1, _ := ret[1].(*pushy.Error)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// UnsubscribeFromTopic indicates an expected call of UnsubscribeFromTopic
func (mr *MockIPushyClientMockRecorder) UnsubscribeFromTopic(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnsubscribeFromTopic", reflect.TypeOf((*MockIPushyClient)(nil).UnsubscribeFromTopic), varargs...)
}
