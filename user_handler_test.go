package crazy_nl_backend_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"bytes"
	"encoding/json"
	"crazy_nl_backend"
	"net/http"
	"github.com/golang/mock/gomock"
	"crazy_nl_backend/mocks"
	"crazy_nl_backend/db"
	"errors"
)

// region User.Create
func TestServer_CreateUserReturnsBadRequestIfNoBody(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	crazy_nl_backend.Server{}.CreateUser()(responseRecorder, request)
	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
}

func TestServer_CreateUserReturnsBadRequestIfUserIsInvalid(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(crazy_nl_backend.NewUser{
		Email: "invalid",
		Password: "pass",
	})
	request := httptest.NewRequest("POST", "/", buffer)
	crazy_nl_backend.Server{}.CreateUser()(responseRecorder, request)
	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
}

func TestServer_CreateUserReturnsConflictStatusIfUserAlreadyPresent(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	buffer := new(bytes.Buffer)
	mockUser := crazy_nl_backend.NewUser{
		Email: "user@example.com",
		Password: "password",
	}
	json.NewEncoder(buffer).Encode(mockUser)
	request := httptest.NewRequest("POST", "/", buffer)
	// mock user manager
	// mock db manager
	userCtrl := gomock.NewController(t)
	defer userCtrl.Finish()
	dbCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	userManager := mocks.NewMockIUserManager(userCtrl)
	userManager.EXPECT().FindByEmail(mockUser.Email).AnyTimes().Return(&db.User{
		Email:mockUser.Email,
		Password:mockUser.Password,
	})
	dbManager := mocks.NewMockDb(dbCtrl)
	dbManager.EXPECT().Clone().Times(1).Return(dbManager)
	dbManager.EXPECT().Close().Times(1)
	dbManager.EXPECT().Users().AnyTimes().Return(userManager)
	crazy_nl_backend.Server{Db:dbManager}.CreateUser()(responseRecorder, request)
	expect.Equal(http.StatusConflict, responseRecorder.Code)
}
func TestServer_CreateUserRespondsWithInternalServerErrorIfDbError(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	buffer := new(bytes.Buffer)
	mockUser := crazy_nl_backend.NewUser{
		Email: "user@example.com",
		Password: "password",
	}
	json.NewEncoder(buffer).Encode(mockUser)
	request := httptest.NewRequest("POST", "/", buffer)
	// mock user manager
	// mock db manager
	userCtrl := gomock.NewController(t)
	defer userCtrl.Finish()
	dbCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	userManager := mocks.NewMockIUserManager(userCtrl)
	userManager.EXPECT().FindByEmail(mockUser.Email).AnyTimes().Return(nil)
	userManager.EXPECT().Create(gomock.Any()).Return(errors.New("db error"))
	dbManager := mocks.NewMockDb(dbCtrl)
	dbManager.EXPECT().Clone().Times(1).Return(dbManager)
	dbManager.EXPECT().Close().Times(1)
	dbManager.EXPECT().Users().AnyTimes().Return(userManager)
	crazy_nl_backend.Server{Db:dbManager}.CreateUser()(responseRecorder, request)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
}

func TestServer_CreateUserRespondsWithStatusCreated(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	buffer := new(bytes.Buffer)
	mockUser := crazy_nl_backend.NewUser{
		Email: "user@example.com",
		Password: "password",
	}
	json.NewEncoder(buffer).Encode(mockUser)
	request := httptest.NewRequest("POST", "/", buffer)
	// mock user manager
	// mock db manager
	userCtrl := gomock.NewController(t)
	defer userCtrl.Finish()
	dbCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	userManager := mocks.NewMockIUserManager(userCtrl)
	userManager.EXPECT().FindByEmail(mockUser.Email).AnyTimes().Return(nil)
	userManager.EXPECT().Create(gomock.Any()).Return(nil)
	dbManager := mocks.NewMockDb(dbCtrl)
	dbManager.EXPECT().Clone().Times(1).Return(dbManager)
	dbManager.EXPECT().Close().Times(1)
	dbManager.EXPECT().Users().AnyTimes().Return(userManager)
	crazy_nl_backend.Server{Db:dbManager}.CreateUser()(responseRecorder, request)
	expect.Equal(http.StatusCreated, responseRecorder.Code)
}
// endregion

// region User.List
func TestServer_ListUsersReturnsInternalServerErrorIfDbError(t *testing.T) {
	expect := assert.New(t)
	dbCtrl := gomock.NewController(t)
	userCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	defer userCtrl.Finish()
	mockDb := mocks.NewMockDb(dbCtrl)
	mockUserManager := mocks.NewMockIUserManager(userCtrl)
	mockUserManager.EXPECT().List().AnyTimes().Return(nil, errors.New("db error"))
	mockDb.EXPECT().Users().AnyTimes().Return(mockUserManager)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	crazy_nl_backend.Server{
		Db: mockDb,
	}.ListUsers()(responseRecorder, request)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
}

func TestServer_ListUsersReturnsListOfUsers(t *testing.T) {
	mockUsers := []db.User{
		{Email:"mail@example.com", Permissions: []string{"sudo"}, Password: "test_password"},
		{Email:"mail2@example.com", Permissions: []string{"sudo"}, Password: "test_password2"},
	}
	expect := assert.New(t)
	dbCtrl := gomock.NewController(t)
	userCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	defer userCtrl.Finish()
	mockDb := mocks.NewMockDb(dbCtrl)
	mockUserManager := mocks.NewMockIUserManager(userCtrl)
	mockUserManager.EXPECT().List().AnyTimes().Return(mockUsers, nil)
	mockDb.EXPECT().Users().AnyTimes().Return(mockUserManager)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	crazy_nl_backend.Server{
		Db: mockDb,
	}.ListUsers()(responseRecorder, request)
	expect.Equal(http.StatusOK, responseRecorder.Code)
	var resUsers []db.User = nil
	json.NewDecoder(responseRecorder.Body).Decode(&resUsers)
	expect.Equal(mockUsers, resUsers)
}
// endregion
