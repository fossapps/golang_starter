package helpers_test

import (
	"golang_starter/helpers"
	"golang_starter/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
