package crazy_nl_backend_test

import (
	"bytes"
	"crazy_nl_backend"
	"crazy_nl_backend/db"
	"crazy_nl_backend/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/globalsign/mgo/bson"
	"crazy_nl_backend/transformers"
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
		Email:    "invalid",
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
		Email:    "user@example.com",
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
		Email:    mockUser.Email,
		Password: mockUser.Password,
	})
	dbManager := mocks.NewMockDb(dbCtrl)
	dbManager.EXPECT().Clone().Times(1).Return(dbManager)
	dbManager.EXPECT().Close().Times(1)
	dbManager.EXPECT().Users().AnyTimes().Return(userManager)
	crazy_nl_backend.Server{Db: dbManager}.CreateUser()(responseRecorder, request)
	expect.Equal(http.StatusConflict, responseRecorder.Code)
}
func TestServer_CreateUserRespondsWithInternalServerErrorIfDbError(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	buffer := new(bytes.Buffer)
	mockUser := crazy_nl_backend.NewUser{
		Email:    "user@example.com",
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
	crazy_nl_backend.Server{Db: dbManager}.CreateUser()(responseRecorder, request)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
}

func TestServer_CreateUserRespondsWithStatusCreated(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	buffer := new(bytes.Buffer)
	mockUser := crazy_nl_backend.NewUser{
		Email:    "user@example.com",
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
	crazy_nl_backend.Server{Db: dbManager}.CreateUser()(responseRecorder, request)
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
		{Email: "mail@example.com", Permissions: []string{"sudo"}, Password: ""},
		{Email: "mail2@example.com", Permissions: []string{"sudo"}, Password: ""},
	}
	expectedUsers := []transformers.ResponseUser {
		{Email: "mail@example.com", Permissions: []string{"sudo"}},
		{Email: "mail2@example.com", Permissions: []string{"sudo"}},
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
	var resUsers []transformers.ResponseUser = nil
	json.NewDecoder(responseRecorder.Body).Decode(&resUsers)
	expect.Equal(expectedUsers, resUsers)
}

// endregion

// region User.Availability
func TestServer_UserAvailabilityRespondsWithBadRequestIfRequestInvalid(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	crazy_nl_backend.Server{}.UserAvailability()(responseRecorder, request)
	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
}

func TestServer_UserAvailabilityReturnsFalseIfUnavailable(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	email := "admin@example.com"
	mockUser := db.User{
		Email: email,
	}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(mockUser)
	request := httptest.NewRequest("GET", "/", buffer)
	dbCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	userCtrl := gomock.NewController(t)
	defer userCtrl.Finish()
	mockDb := mocks.NewMockDb(dbCtrl)
	mockUserManager := mocks.NewMockIUserManager(userCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockUserManager.EXPECT().FindByEmail(email).Times(1).Return(&mockUser)
	mockDb.EXPECT().Users().AnyTimes().Return(mockUserManager)
	crazy_nl_backend.Server{Db: mockDb}.UserAvailability()(responseRecorder, request)
	expect.Equal(http.StatusOK, responseRecorder.Code)
	response := new(crazy_nl_backend.UserAvailabilityResponse)
	json.NewDecoder(responseRecorder.Body).Decode(&response)
	expect.False(response.Available)
}

func TestServer_UserAvailabilityReturnsTrueIfAvailable(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	email := "admin@example.com"
	mockUser := db.User{
		Email: email,
	}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(mockUser)
	request := httptest.NewRequest("GET", "/", buffer)
	dbCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	userCtrl := gomock.NewController(t)
	defer userCtrl.Finish()
	mockDb := mocks.NewMockDb(dbCtrl)
	mockUserManager := mocks.NewMockIUserManager(userCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockUserManager.EXPECT().FindByEmail(email).Times(1).Return(nil)
	mockDb.EXPECT().Users().AnyTimes().Return(mockUserManager)
	crazy_nl_backend.Server{Db: mockDb}.UserAvailability()(responseRecorder, request)
	expect.Equal(http.StatusOK, responseRecorder.Code)
	response := new(crazy_nl_backend.UserAvailabilityResponse)
	json.NewDecoder(responseRecorder.Body).Decode(&response)
	expect.True(response.Available)
}

// endregion

// region User.Edit

func TestServer_EditUserErrorIfUserDoesNotExist(t *testing.T) {
	expect := assert.New(t)
	dbCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	mockDb := mocks.NewMockDb(dbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	userCtrl := gomock.NewController(t)
	defer userCtrl.Finish()
	mockUser := mocks.NewMockIUserManager(userCtrl)
	mockUser.EXPECT().FindById("id").Return(nil)
	mockDb.EXPECT().Users().AnyTimes().Return(mockUser)

	router := mux.NewRouter()
	server := crazy_nl_backend.Server{
		Db: mockDb,
	}
	router.HandleFunc("/users/{user}", server.EditUser()).Methods("PUT")
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("PUT", "/users/id", nil)
	router.ServeHTTP(responseRecorder, request)
	expect.Equal(http.StatusPreconditionFailed, responseRecorder.Code)
}

func TestServer_EditUserErrorIfUserIsInvalid(t *testing.T) {
	expect := assert.New(t)
	dbCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	mockDb := mocks.NewMockDb(dbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	userCtrl := gomock.NewController(t)
	defer userCtrl.Finish()
	mockUser := mocks.NewMockIUserManager(userCtrl)
	mockUser.EXPECT().FindById("id").Return(&db.User{
		Email:       "mail@example.com",
		Permissions: []string{"sudo"},
	})
	mockDb.EXPECT().Users().AnyTimes().Return(mockUser)

	router := mux.NewRouter()
	server := crazy_nl_backend.Server{
		Db: mockDb,
	}
	router.HandleFunc("/users/{user}", server.EditUser()).Methods("PUT")
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(crazy_nl_backend.NewUser{
		Email: "new_email.example.com",
	})
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("PUT", "/users/id", buffer)
	router.ServeHTTP(responseRecorder, request)
	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
}

func TestServer_EditUserHandlesDbNotFoundError(t *testing.T) {
	expect := assert.New(t)
	dbCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	mockDb := mocks.NewMockDb(dbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	userCtrl := gomock.NewController(t)
	defer userCtrl.Finish()
	mockUserManager := mocks.NewMockIUserManager(userCtrl)
	mockUserManager.EXPECT().FindById("id").Return(&db.User{
		Email:       "mail@example.com",
		Permissions: []string{"sudo"},
	})
	mockUserManager.EXPECT().Edit("id", gomock.Any()).Times(1).Return(mgo.ErrNotFound)
	mockDb.EXPECT().Users().AnyTimes().Return(mockUserManager)

	router := mux.NewRouter()
	server := crazy_nl_backend.Server{
		Db: mockDb,
	}
	router.HandleFunc("/users/{user}", server.EditUser()).Methods("PUT")
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(crazy_nl_backend.NewUser{
		Email: "new_email@example.com",
	})
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("PUT", "/users/id", buffer)
	router.ServeHTTP(responseRecorder, request)
	expect.Equal(http.StatusPreconditionFailed, responseRecorder.Code)
}

func TestServer_EditUserHandlesDbError(t *testing.T) {
	expect := assert.New(t)
	dbCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	mockDb := mocks.NewMockDb(dbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	userCtrl := gomock.NewController(t)
	defer userCtrl.Finish()
	mockUserManager := mocks.NewMockIUserManager(userCtrl)
	mockUserManager.EXPECT().FindById("id").Return(&db.User{
		Email:       "mail@example.com",
		Permissions: []string{"sudo"},
	})
	mockUserManager.EXPECT().Edit("id", gomock.Any()).Times(1).Return(errors.New("db error"))
	mockDb.EXPECT().Users().AnyTimes().Return(mockUserManager)

	router := mux.NewRouter()
	server := crazy_nl_backend.Server{
		Db: mockDb,
	}
	router.HandleFunc("/users/{user}", server.EditUser()).Methods("PUT")
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(crazy_nl_backend.NewUser{
		Email: "new_email@example.com",
	})
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("PUT", "/users/id", buffer)
	router.ServeHTTP(responseRecorder, request)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
}

func TestServer_EditUserReturnsOkWhenValid(t *testing.T) {
	expect := assert.New(t)
	dbCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	mockDb := mocks.NewMockDb(dbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	userCtrl := gomock.NewController(t)
	defer userCtrl.Finish()
	mockUserManager := mocks.NewMockIUserManager(userCtrl)
	mockUserManager.EXPECT().FindById("id").AnyTimes().Return(&db.User{
		Email:       "mail@example.com",
		Permissions: []string{"sudo"},
	})
	mockUserManager.EXPECT().Edit("id", gomock.Any()).Times(1).Return(nil)
	mockDb.EXPECT().Users().AnyTimes().Return(mockUserManager)

	router := mux.NewRouter()
	server := crazy_nl_backend.Server{
		Db: mockDb,
	}
	router.HandleFunc("/users/{user}", server.EditUser()).Methods("PUT")
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(crazy_nl_backend.NewUser{
		Email:    "new_email@example.com",
		Password: "",
	})
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("PUT", "/users/id", buffer)
	router.ServeHTTP(responseRecorder, request)
	expect.Equal(http.StatusOK, responseRecorder.Code)
}

// endregion

// region User.Get

func TestServer_UserGetReturnsStatusNotFoundWhenUserDoesNotExist(t *testing.T) {
	expect := assert.New(t)
	dbCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	mockDb := mocks.NewMockDb(dbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	userCtrl := gomock.NewController(t)
	defer userCtrl.Finish()
	mockUser := mocks.NewMockIUserManager(userCtrl)
	mockUser.EXPECT().FindById("id").Return(nil)
	mockDb.EXPECT().Users().AnyTimes().Return(mockUser)

	router := mux.NewRouter()
	server := crazy_nl_backend.Server{
		Db: mockDb,
	}
	router.HandleFunc("/users/{user}", server.GetUser()).Methods("GET")
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/users/id", nil)
	router.ServeHTTP(responseRecorder, request)
	expect.Equal(http.StatusNotFound, responseRecorder.Code)
}

func TestServer_GetUserReturnsUser(t *testing.T) {
	expect := assert.New(t)
	dbCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	mockDb := mocks.NewMockDb(dbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	userCtrl := gomock.NewController(t)
	defer userCtrl.Finish()
	mockUser := mocks.NewMockIUserManager(userCtrl)
	dbUser := db.User{
		ID: bson.NewObjectId(),
		Email: "example@admin.com",
		Permissions:[]string{"users.create"},
	}
	mockUser.EXPECT().FindById("id").Return(&dbUser)
	mockDb.EXPECT().Users().AnyTimes().Return(mockUser)

	router := mux.NewRouter()
	server := crazy_nl_backend.Server{
		Db: mockDb,
	}
	router.HandleFunc("/users/{user}", server.GetUser()).Methods("GET")
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/users/id", nil)
	router.ServeHTTP(responseRecorder, request)
	expect.Equal(http.StatusOK, responseRecorder.Code)
	var responseUser transformers.ResponseUser
	json.NewDecoder(responseRecorder.Body).Decode(&responseUser)
	expect.Equal("example@admin.com", responseUser.Email)
	expect.Equal(dbUser.ID.Hex(), responseUser.ID.Hex())
}
// endregion
