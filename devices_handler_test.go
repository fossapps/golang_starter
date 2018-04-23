//go:generate mockgen -destination=./mocks/mock_pushy_client.go -package=mocks github.com/cyberhck/pushy IPushyClient

package crazy_nl_backend_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"crazy_nl_backend/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"crazy_nl_backend"
)

func TestServer_RegisterHandlerReturnsInvalidRequestIfJsonInvalid(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	server := crazy_nl_backend.Server{}
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

func TestServer_RegisterHandlerReturnsInvalidTokenIfTokenIsInvalid(t *testing.T) {
	expect := assert.New(t)
	server := crazy_nl_backend.Server{}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(crazy_nl_backend.NewRegistration{
		Token: "token",
	})
	request := httptest.NewRequest("POST", "/", buffer)
	responseRecorder := httptest.NewRecorder()
	server.RegisterHandler()(responseRecorder, request)
	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
}

func TestServer_RegisterHandlerIfExistsItReturnsBadRequest(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	mockDevicesCtrl := gomock.NewController(t)
	mockDb := mocks.NewMockDb(mockDbCtrl)
	mockDeviceManager := mocks.NewMockIDeviceManager(mockDevicesCtrl)
	token := "some_random_large_token_which_is_checked"
	mockDeviceManager.EXPECT().Exists(token).Times(1).Return(true)
	mockDb.EXPECT().Devices().Times(1).Return(mockDeviceManager)
	server := crazy_nl_backend.Server{
		Db: mockDb,
	}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(crazy_nl_backend.NewRegistration{
		Token: token,
	})
	request := httptest.NewRequest("POST", "/", buffer)
	responseRecorder := httptest.NewRecorder()
	server.RegisterHandler()(responseRecorder, request)
	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
	res := new(crazy_nl_backend.ErrorResponse)
	json.NewDecoder(responseRecorder.Body).Decode(&res)
	expect.Equal("already registered", res.Message)
}

func TestServer_RegisterHandlerReturnsInternalServerIfDbError(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	mockDevicesCtrl := gomock.NewController(t)
	mockDb := mocks.NewMockDb(mockDbCtrl)
	mockDeviceManager := mocks.NewMockIDeviceManager(mockDevicesCtrl)
	token := "some_random_large_token_which_is_checked"
	mockDeviceManager.EXPECT().Exists(token).Times(1).Return(false)
	mockDeviceManager.EXPECT().Register(token).Times(1).Return(errors.New("db error"))
	mockDb.EXPECT().Devices().MinTimes(1).Return(mockDeviceManager)
	server := crazy_nl_backend.Server{
		Db: mockDb,
	}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(crazy_nl_backend.NewRegistration{
		Token: token,
	})
	request := httptest.NewRequest("POST", "/", buffer)
	responseRecorder := httptest.NewRecorder()
	server.RegisterHandler()(responseRecorder, request)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
	res := new(crazy_nl_backend.ErrorResponse)
	json.NewDecoder(responseRecorder.Body).Decode(&res)
	expect.Equal("db error", res.Message)
}

func TestServer_RegisterHandlerRegisters(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	mockDevicesCtrl := gomock.NewController(t)
	mockDb := mocks.NewMockDb(mockDbCtrl)
	mockDeviceManager := mocks.NewMockIDeviceManager(mockDevicesCtrl)
	token := "some_random_large_token_which_is_checked"
	mockDeviceManager.EXPECT().Exists(token).Times(1).Return(false)
	mockDeviceManager.EXPECT().Register(token).Times(1).Return(nil)
	mockDb.EXPECT().Devices().MinTimes(1).Return(mockDeviceManager)
	server := crazy_nl_backend.Server{
		Db: mockDb,
	}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(crazy_nl_backend.NewRegistration{
		Token: token,
	})
	request := httptest.NewRequest("POST", "/", buffer)
	responseRecorder := httptest.NewRecorder()
	server.RegisterHandler()(responseRecorder, request)
	expect.Equal(http.StatusOK, responseRecorder.Code)
	res := new(crazy_nl_backend.RegistrationResponse)
	json.NewDecoder(responseRecorder.Body).Decode(&res)
	expect.Equal("success", res.Status)
}
