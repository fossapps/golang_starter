package crazy_nl_backend_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"crazy_nl_backend"
	"crazy_nl_backend/db"
	"crazy_nl_backend/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestServer_ListPermissionsListsPermissionFromDb(t *testing.T) {
	expect := assert.New(t)
	permissionCtrl := gomock.NewController(t)
	defer permissionCtrl.Finish()
	permissionManager := mocks.NewMockIPermissionManager(permissionCtrl)
	dbCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	dbManager := mocks.NewMockDb(dbCtrl)
	dbManager.EXPECT().Clone().AnyTimes().Return(dbManager)
	dbManager.EXPECT().Close().Times(1)
	server := crazy_nl_backend.Server{
		Db: dbManager,
	}
	mockPermissionList := []db.Permission{
		{
			Description: "test description",
			Key:         "key",
		},
		{
			Description: "second description",
			Key:         "second.key",
		},
	}
	permissionManager.EXPECT().List().Times(1).Return(mockPermissionList, nil)
	dbManager.EXPECT().Permissions().Times(1).Return(permissionManager)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	server.ListPermissions()(responseRecorder, request)
	response := new([]db.Permission)
	json.NewDecoder(responseRecorder.Body).Decode(&response)
	expect.Equal(http.StatusOK, responseRecorder.Code)
	expect.Equal(response, &mockPermissionList)
}

func TestServer_ListPermissionsReturnsInternalServerErrorIfDbError(t *testing.T) {
	expect := assert.New(t)
	permissionCtrl := gomock.NewController(t)
	defer permissionCtrl.Finish()
	permissionManager := mocks.NewMockIPermissionManager(permissionCtrl)
	dbCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	dbManager := mocks.NewMockDb(dbCtrl)
	dbManager.EXPECT().Clone().AnyTimes().Return(dbManager)
	dbManager.EXPECT().Close().Times(1)
	server := crazy_nl_backend.Server{
		Db: dbManager,
	}
	permissionManager.EXPECT().List().Times(1).Return(nil, errors.New("db error"))
	dbManager.EXPECT().Permissions().Times(1).Return(permissionManager)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	server.ListPermissions()(responseRecorder, request)
	response := new([]db.Permission)
	json.NewDecoder(responseRecorder.Body).Decode(&response)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
}
