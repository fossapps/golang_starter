package helpers_test

import (
	"testing"
	"github.com/stretchr/testify/mock"
	"crazy_nl_backend/helpers"
	"github.com/stretchr/testify/assert"
)

type redis struct {
	mock.Mock
}

func (m redis) SIsMember(key string, member interface{}) (bool, error) {
	args := m.Called(key, member)
	return args.Bool(0), args.Error(1)
}

func (m redis) SAdd(key string, value ...interface{}) (int64, error) {
	c := append([] interface{}{key}, value...)
	args := m.Called(c...)
	return int64(args.Int(0)), args.Error(1)
}

func (m redis) Close() error {
	return nil
}

func (m redis) SPop(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func TestQueueDeviceRegistrationReturnsErrorIfTokenAlreadyExists(t *testing.T) {
	mockRedis := new(redis)
	mockRedis.On("SIsMember", "registration", "token").
		Return(true, nil)
	err := helpers.QueueDeviceRegistration("token", mockRedis)
	assert.Error(t, err, "token already exists")
	mockRedis.AssertExpectations(t)
}

func TestQueueDeviceRegistrationReturnsNoErrorAndAddsToClientIfNotAlreadyPresent(t *testing.T) {
	mockRedis := new(redis)
	mockRedis.On("SIsMember", "registration", "token").
		Return(false, nil)
	mockRedis.On("SAdd", "registration", "token").Return(1, nil)
	helpers.QueueDeviceRegistration("token", mockRedis)
	mockRedis.AssertExpectations(t)
}
