package starter_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fossapps/starter/mock"

	"github.com/cyberhck/pushy"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/fossapps/starter"
)

func TestServer_RegisterHandlerReturnsUnprocessableEntityRequestIfJsonInvalid(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	server := starter.Server{}
	server.RegisterHandler()(responseRecorder, request)
	expect.Equal(http.StatusUnprocessableEntity, responseRecorder.Code)
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(struct {
		Body string
	}{
		Body: "some random thing",
	})
	request = httptest.NewRequest("POST", "/", buffer)
	responseRecorder = httptest.NewRecorder()
	server.RegisterHandler()(responseRecorder, request)
	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
}

func TestServer_RegisterHandlerReturnsBadRequestIfTokenIsInvalid(t *testing.T) {
	expect := assert.New(t)
	server := starter.Server{}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(starter.NewRegistration{
		Token: "token",
	})
	request := httptest.NewRequest("POST", "/", buffer)
	responseRecorder := httptest.NewRecorder()
	server.RegisterHandler()(responseRecorder, request)
	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
}

func TestServer_RegisterHandlerReturnsBadRequestIfDuplicate(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	defer mockDbCtrl.Finish()
	mockDevicesCtrl := gomock.NewController(t)
	defer mockDevicesCtrl.Finish()
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockDeviceManager := mock.NewMockDeviceManager(mockDevicesCtrl)
	token := "some_random_large_token_which_is_checked"
	mockDeviceManager.EXPECT().Exists(token).Times(1).Return(true, nil)
	mockDb.EXPECT().Devices().Times(1).Return(mockDeviceManager)
	server := starter.Server{
		Db:    mockDb,
	}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(starter.NewRegistration{
		Token: token,
	})
	request := httptest.NewRequest("POST", "/", buffer)
	responseRecorder := httptest.NewRecorder()
	server.RegisterHandler()(responseRecorder, request)
	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
	res := new(starter.SimpleResponse)
	json.NewDecoder(responseRecorder.Body).Decode(&res)
	expect.Equal("already registered", res.Message)
}

func TestServer_RegisterHandlerHandlesDbError(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	defer mockDbCtrl.Finish()
	mockDevicesCtrl := gomock.NewController(t)
	defer mockDevicesCtrl.Finish()
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockDeviceManager := mock.NewMockDeviceManager(mockDevicesCtrl)
	token := "some_random_large_token_which_is_checked"
	mockDeviceManager.EXPECT().Exists(token).Times(1).Return(false, errors.New("db error"))
	mockDb.EXPECT().Devices().Times(1).Return(mockDeviceManager)
	server := starter.Server{
		Db:    mockDb,
	}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(starter.NewRegistration{
		Token: token,
	})
	request := httptest.NewRequest("POST", "/", buffer)
	responseRecorder := httptest.NewRecorder()
	server.RegisterHandler()(responseRecorder, request)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
}

func TestServer_RegisterHandlerReturnsInternalServerIfDbError(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	defer mockDbCtrl.Finish()
	mockDevicesCtrl := gomock.NewController(t)
	defer mockDevicesCtrl.Finish()
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockDeviceManager := mock.NewMockDeviceManager(mockDevicesCtrl)
	pushyCtrl := gomock.NewController(t)
	defer pushyCtrl.Finish()
	mockPushy := mock.NewMockIPushyClient(pushyCtrl)
	token := "some_random_large_token_which_is_checked"
	mockPushy.EXPECT().DeviceInfo(token).Return(nil, nil, nil)
	mockDeviceManager.EXPECT().Exists(token).Times(1).Return(false, nil)
	mockDeviceManager.EXPECT().Register(token).Times(1).Return(errors.New("db error"))
	mockDb.EXPECT().Devices().MinTimes(1).Return(mockDeviceManager)
	server := starter.Server{
		Db:    mockDb,
		Pushy: mockPushy,
	}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(starter.NewRegistration{
		Token: token,
	})
	request := httptest.NewRequest("POST", "/", buffer)
	responseRecorder := httptest.NewRecorder()
	server.RegisterHandler()(responseRecorder, request)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
	res := new(starter.SimpleResponse)
	json.NewDecoder(responseRecorder.Body).Decode(&res)
	expect.Equal("db error", res.Message)
}

func TestServer_RegisterHandlerRegisters(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	defer mockDbCtrl.Finish()
	mockDevicesCtrl := gomock.NewController(t)
	defer mockDevicesCtrl.Finish()
	pushyCtrl := gomock.NewController(t)
	defer pushyCtrl.Finish()
	mockPushy := mock.NewMockIPushyClient(pushyCtrl)
	token := "some_random_large_token_which_is_checked"
	mockPushy.EXPECT().DeviceInfo(token).Return(nil, nil, nil)
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockDeviceManager := mock.NewMockDeviceManager(mockDevicesCtrl)
	mockDeviceManager.EXPECT().Exists(token).Times(1).Return(false, nil)
	mockDeviceManager.EXPECT().Register(token).Times(1).Return(nil)
	mockDb.EXPECT().Devices().MinTimes(1).Return(mockDeviceManager)
	server := starter.Server{
		Db:    mockDb,
		Pushy: mockPushy,
	}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(starter.NewRegistration{
		Token: token,
	})
	request := httptest.NewRequest("POST", "/", buffer)
	responseRecorder := httptest.NewRecorder()
	server.RegisterHandler()(responseRecorder, request)
	expect.Equal(http.StatusOK, responseRecorder.Code)
	res := new(starter.SimpleResponse)
	json.NewDecoder(responseRecorder.Body).Decode(&res)
	expect.True(res.Success)
	expect.Equal("success", res.Message)
}

func TestServer_RegisterHandlerRespondsWithBadRequestIfDeviceTokenInvalid(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	defer mockDbCtrl.Finish()
	mockDb := mock.NewMockDB(mockDbCtrl)
	deviceManagerCtrl := gomock.NewController(t)
	deviceManager := mock.NewMockDeviceManager(deviceManagerCtrl)
	pushyCtrl := gomock.NewController(t)
	defer pushyCtrl.Finish()
	mockPushy := mock.NewMockIPushyClient(pushyCtrl)
	token := "some_random_large_token_which_is_checked"
	mockPushyError := pushy.Error{
		Error: "We could not find a device with that token linked to your account.",
	}
	mockPushy.EXPECT().DeviceInfo(token).Return(nil, &mockPushyError, nil)
	deviceManager.EXPECT().Exists(token).AnyTimes().Return(false, nil)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockDb.EXPECT().Devices().AnyTimes().Return(deviceManager)
	server := starter.Server{
		Db:    mockDb,
		Pushy: mockPushy,
	}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(starter.NewRegistration{
		Token: token,
	})
	request := httptest.NewRequest("POST", "/", buffer)
	responseRecorder := httptest.NewRecorder()
	server.RegisterHandler()(responseRecorder, request)
	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
	res := new(starter.SimpleResponse)
	json.NewDecoder(responseRecorder.Body).Decode(&res)
	expect.Equal("invalid token", res.Message)
}
