//go:generate mockgen -destination=../mocks/mock_redis.go -package=mocks crazy_nl_backend/helpers IRedisClient

package helpers_test

import (
	"testing"
	"crazy_nl_backend/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/golang/mock/gomock"
	"crazy_nl_backend/mocks"
)

func TestQueueDeviceRegistrationReturnsErrorIfTokenAlreadyExists(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRedis := mocks.NewMockIRedisClient(controller)
	mockRedis.EXPECT().SIsMember("registration", "token").Return(true, nil)
	err := helpers.QueueDeviceRegistration("token", mockRedis)
	assert.Error(t, err, "token already exists")
}

func TestQueueDeviceRegistrationReturnsNoErrorAndAddsToClientIfNotAlreadyPresent(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRedis := mocks.NewMockIRedisClient(controller)
	mockRedis.EXPECT().SIsMember("registration", "token").Return(false, nil)
	mockRedis.EXPECT().SAdd("registration", "token").Return(int64(1), nil)
	helpers.QueueDeviceRegistration("token", mockRedis)
}

func TestRegisterDevice(t *testing.T) {
	
}
