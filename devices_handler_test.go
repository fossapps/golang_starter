//go:generate mockgen -destination=./mocks/mock_pushy_client.go -package=mocks github.com/cyberhck/pushy IPushyClient

package crazy_nl_backend_test

import (
	"testing"
	"github.com/golang/mock/gomock"
	"crazy_nl_backend/mocks"
)

func TestServer_RegisterHandlerReturnsInvalidRequestIfJsonInvalid(t *testing.T) {
}

func TestServer_RegisterHandlerReturnsInvalidTokenIfTokenIsInvalid(t *testing.T) {

}

func TestServer_RegisterHandlerIfExistsItReturnsBadRequest(t *testing.T) {
	mockDbCtrl := gomock.NewController(t)
	mockDevicesCtrl := gomock.NewController(t)
	mockDb := mocks.NewMockDb(mockDbCtrl)
	mockDeviceManager := mocks.NewMockIDeviceManager(mockDevicesCtrl)
}
