//go:generate mockgen -destination=./mocks/mock_pushy_client.go -package=mocks github.com/cyberhck/pushy IPushyClient

package starter_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fossapps/starter/mocks"

	"github.com/cyberhck/pushy"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/fossapps/starter"
)

func TestServer_RegisterHandlerReturnsBadRequestRequestIfJsonInvalid(t *testing.T) {
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
	mockDb := mocks.NewMockDB(mockDbCtrl)
	mockDeviceManager := mocks.NewMockDeviceManager(mockDevicesCtrl)
	token := "some_random_large_token_which_is_checked"
	pushyCtrl := gomock.NewController(t)
	defer pushyCtrl.Finish()
	mockPushy := mocks.NewMockIPushyClient(pushyCtrl)
	mockPushy.EXPECT().DeviceInfo(token).Return(nil, nil, nil)
	mockDeviceManager.EXPECT().Exists(token).Times(1).Return(true)
	mockDb.EXPECT().Devices().Times(1).Return(mockDeviceManager)
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
	expect.Equal("already registered", res.Message)
}

func TestServer_RegisterHandlerReturnsInternalServerIfDbError(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	defer mockDbCtrl.Finish()
	mockDevicesCtrl := gomock.NewController(t)
	defer mockDevicesCtrl.Finish()
	mockDb := mocks.NewMockDB(mockDbCtrl)
	mockDeviceManager := mocks.NewMockDeviceManager(mockDevicesCtrl)
	pushyCtrl := gomock.NewController(t)
	defer pushyCtrl.Finish()
	mockPushy := mocks.NewMockIPushyClient(pushyCtrl)
	token := "some_random_large_token_which_is_checked"
	mockPushy.EXPECT().DeviceInfo(token).Return(nil, nil, nil)
	mockDeviceManager.EXPECT().Exists(token).Times(1).Return(false)
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
	mockPushy := mocks.NewMockIPushyClient(pushyCtrl)
	token := "some_random_large_token_which_is_checked"
	mockPushy.EXPECT().DeviceInfo(token).Return(nil, nil, nil)
	mockDb := mocks.NewMockDB(mockDbCtrl)
	mockDeviceManager := mocks.NewMockDeviceManager(mockDevicesCtrl)
	mockDeviceManager.EXPECT().Exists(token).Times(1).Return(false)
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
	mockDb := mocks.NewMockDB(mockDbCtrl)
	pushyCtrl := gomock.NewController(t)
	defer pushyCtrl.Finish()
	mockPushy := mocks.NewMockIPushyClient(pushyCtrl)
	token := "some_random_large_token_which_is_checked"
	mockPushyError := pushy.Error{
		Error: "We could not find a device with that token linked to your account.",
	}
	mockPushy.EXPECT().DeviceInfo(token).Return(nil, &mockPushyError, nil)
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
