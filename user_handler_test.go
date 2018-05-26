package starter_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fossapps/starter"
	"github.com/fossapps/starter/db"
	"github.com/fossapps/starter/mock"

	"github.com/fossapps/starter/transformer"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// region User.Create
func TestServer_CreateUserReturnsBadRequestIfNoBody(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	starter.Server{}.CreateUser()(responseRecorder, request)
	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
}

func TestServer_CreateUserReturnsBadRequestIfUserIsInvalid(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(starter.NewUser{
		Email:    "invalid",
		Password: "pass",
	})
	request := httptest.NewRequest("POST", "/", buffer)
	starter.Server{}.CreateUser()(responseRecorder, request)
	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
}

func TestServer_CreateUserReturnsConflictStatusIfUserAlreadyPresent(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	buffer := new(bytes.Buffer)
	mockUser := starter.NewUser{
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
	userManager := mock.NewMockUserManager(userCtrl)
	userManager.EXPECT().FindByEmail(mockUser.Email).AnyTimes().Return(&db.User{
		Email:    mockUser.Email,
		Password: mockUser.Password,
	})
	dbManager := mock.NewMockDB(dbCtrl)
	dbManager.EXPECT().Clone().Times(1).Return(dbManager)
	dbManager.EXPECT().Close().Times(1)
	dbManager.EXPECT().Users().AnyTimes().Return(userManager)
	starter.Server{Db: dbManager}.CreateUser()(responseRecorder, request)
	expect.Equal(http.StatusConflict, responseRecorder.Code)
}
func TestServer_CreateUserRespondsWithInternalServerErrorIfDbError(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	buffer := new(bytes.Buffer)
	mockUser := starter.NewUser{
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
	userManager := mock.NewMockUserManager(userCtrl)
	userManager.EXPECT().FindByEmail(mockUser.Email).AnyTimes().Return(nil)
	userManager.EXPECT().Create(gomock.Any()).Return(errors.New("db error"))
	dbManager := mock.NewMockDB(dbCtrl)
	dbManager.EXPECT().Clone().Times(1).Return(dbManager)
	dbManager.EXPECT().Close().Times(1)
	dbManager.EXPECT().Users().AnyTimes().Return(userManager)
	starter.Server{Db: dbManager}.CreateUser()(responseRecorder, request)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
}

func TestServer_CreateUserRespondsWithStatusCreated(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	buffer := new(bytes.Buffer)
	mockUser := starter.NewUser{
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
	userManager := mock.NewMockUserManager(userCtrl)
	userManager.EXPECT().FindByEmail(mockUser.Email).AnyTimes().Return(nil)
	userManager.EXPECT().Create(gomock.Any()).Return(nil)
	dbManager := mock.NewMockDB(dbCtrl)
	dbManager.EXPECT().Clone().Times(1).Return(dbManager)
	dbManager.EXPECT().Close().Times(1)
	dbManager.EXPECT().Users().AnyTimes().Return(userManager)
	starter.Server{Db: dbManager}.CreateUser()(responseRecorder, request)
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
	mockDb := mock.NewMockDB(dbCtrl)
	mockUserManager := mock.NewMockUserManager(userCtrl)
	mockUserManager.EXPECT().List().AnyTimes().Return(nil, errors.New("db error"))
	mockDb.EXPECT().Users().AnyTimes().Return(mockUserManager)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	starter.Server{
		Db: mockDb,
	}.ListUsers()(responseRecorder, request)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
}

func TestServer_ListUsersReturnsListOfUsers(t *testing.T) {
	mockUsers := []db.User{
		{Email: "mail@example.com", Permissions: []string{"sudo"}, Password: ""},
		{Email: "mail2@example.com", Permissions: []string{"sudo"}, Password: ""},
	}
	expectedUsers := []transformer.ResponseUser{
		{Email: "mail@example.com", Permissions: []string{"sudo"}},
		{Email: "mail2@example.com", Permissions: []string{"sudo"}},
	}
	expect := assert.New(t)
	dbCtrl := gomock.NewController(t)
	userCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	defer userCtrl.Finish()
	mockDb := mock.NewMockDB(dbCtrl)
	mockUserManager := mock.NewMockUserManager(userCtrl)
	mockUserManager.EXPECT().List().AnyTimes().Return(mockUsers, nil)
	mockDb.EXPECT().Users().AnyTimes().Return(mockUserManager)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	starter.Server{
		Db: mockDb,
	}.ListUsers()(responseRecorder, request)
	expect.Equal(http.StatusOK, responseRecorder.Code)
	var resUsers []transformer.ResponseUser = nil
	json.NewDecoder(responseRecorder.Body).Decode(&resUsers)
	expect.Equal(expectedUsers, resUsers)
}

// endregion

// region User.Availability
func TestServer_UserAvailabilityRespondsWithBadRequestIfRequestInvalid(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	starter.Server{}.UserAvailability()(responseRecorder, request)
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
	mockDb := mock.NewMockDB(dbCtrl)
	mockUserManager := mock.NewMockUserManager(userCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockUserManager.EXPECT().FindByEmail(email).Times(1).Return(&mockUser)
	mockDb.EXPECT().Users().AnyTimes().Return(mockUserManager)
	starter.Server{Db: mockDb}.UserAvailability()(responseRecorder, request)
	expect.Equal(http.StatusOK, responseRecorder.Code)
	response := new(starter.UserAvailabilityResponse)
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
	mockDb := mock.NewMockDB(dbCtrl)
	mockUserManager := mock.NewMockUserManager(userCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockUserManager.EXPECT().FindByEmail(email).Times(1).Return(nil)
	mockDb.EXPECT().Users().AnyTimes().Return(mockUserManager)
	starter.Server{Db: mockDb}.UserAvailability()(responseRecorder, request)
	expect.Equal(http.StatusOK, responseRecorder.Code)
	response := new(starter.UserAvailabilityResponse)
	json.NewDecoder(responseRecorder.Body).Decode(&response)
	expect.True(response.Available)
}

// endregion

// region User.Edit

func TestServer_EditUserErrorIfUserDoesNotExist(t *testing.T) {
	expect := assert.New(t)
	dbCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	mockDb := mock.NewMockDB(dbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	userCtrl := gomock.NewController(t)
	defer userCtrl.Finish()
	mockUser := mock.NewMockUserManager(userCtrl)
	mockUser.EXPECT().FindByID("id").Return(nil)
	mockDb.EXPECT().Users().AnyTimes().Return(mockUser)

	router := mux.NewRouter()
	server := starter.Server{
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
	mockDb := mock.NewMockDB(dbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	userCtrl := gomock.NewController(t)
	defer userCtrl.Finish()
	mockUser := mock.NewMockUserManager(userCtrl)
	mockUser.EXPECT().FindByID("id").Return(&db.User{
		Email:       "mail@example.com",
		Permissions: []string{"sudo"},
	})
	mockDb.EXPECT().Users().AnyTimes().Return(mockUser)

	router := mux.NewRouter()
	server := starter.Server{
		Db: mockDb,
	}
	router.HandleFunc("/users/{user}", server.EditUser()).Methods("PUT")
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(starter.NewUser{
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
	mockDb := mock.NewMockDB(dbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	userCtrl := gomock.NewController(t)
	defer userCtrl.Finish()
	mockUserManager := mock.NewMockUserManager(userCtrl)
	mockUserManager.EXPECT().FindByID("id").Return(&db.User{
		Email:       "mail@example.com",
		Permissions: []string{"sudo"},
	})
	mockUserManager.EXPECT().Edit("id", gomock.Any()).Times(1).Return(mgo.ErrNotFound)
	mockDb.EXPECT().Users().AnyTimes().Return(mockUserManager)

	router := mux.NewRouter()
	server := starter.Server{
		Db: mockDb,
	}
	router.HandleFunc("/users/{user}", server.EditUser()).Methods("PUT")
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(starter.NewUser{
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
	mockDb := mock.NewMockDB(dbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	userCtrl := gomock.NewController(t)
	defer userCtrl.Finish()
	mockUserManager := mock.NewMockUserManager(userCtrl)
	mockUserManager.EXPECT().FindByID("id").Return(&db.User{
		Email:       "mail@example.com",
		Permissions: []string{"sudo"},
	})
	mockUserManager.EXPECT().Edit("id", gomock.Any()).Times(1).Return(errors.New("db error"))
	mockDb.EXPECT().Users().AnyTimes().Return(mockUserManager)

	router := mux.NewRouter()
	server := starter.Server{
		Db: mockDb,
	}
	router.HandleFunc("/users/{user}", server.EditUser()).Methods("PUT")
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(starter.NewUser{
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
	mockDb := mock.NewMockDB(dbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	userCtrl := gomock.NewController(t)
	defer userCtrl.Finish()
	mockUserManager := mock.NewMockUserManager(userCtrl)
	mockUserManager.EXPECT().FindByID("id").AnyTimes().Return(&db.User{
		Email:       "mail@example.com",
		Permissions: []string{"sudo"},
	})
	mockUserManager.EXPECT().Edit("id", gomock.Any()).Times(1).Return(nil)
	mockDb.EXPECT().Users().AnyTimes().Return(mockUserManager)

	router := mux.NewRouter()
	server := starter.Server{
		Db: mockDb,
	}
	router.HandleFunc("/users/{user}", server.EditUser()).Methods("PUT")
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(starter.NewUser{
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
	mockDb := mock.NewMockDB(dbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	userCtrl := gomock.NewController(t)
	defer userCtrl.Finish()
	mockUser := mock.NewMockUserManager(userCtrl)
	mockUser.EXPECT().FindByID("id").Return(nil)
	mockDb.EXPECT().Users().AnyTimes().Return(mockUser)

	router := mux.NewRouter()
	server := starter.Server{
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
	mockDb := mock.NewMockDB(dbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	userCtrl := gomock.NewController(t)
	defer userCtrl.Finish()
	mockUser := mock.NewMockUserManager(userCtrl)
	dbUser := db.User{
		ID:          bson.NewObjectId(),
		Email:       "example@admin.com",
		Permissions: []string{"users.create"},
	}
	mockUser.EXPECT().FindByID("id").Return(&dbUser)
	mockDb.EXPECT().Users().AnyTimes().Return(mockUser)

	router := mux.NewRouter()
	server := starter.Server{
		Db: mockDb,
	}
	router.HandleFunc("/users/{user}", server.GetUser()).Methods("GET")
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/users/id", nil)
	router.ServeHTTP(responseRecorder, request)
	expect.Equal(http.StatusOK, responseRecorder.Code)
	var responseUser transformer.ResponseUser
	json.NewDecoder(responseRecorder.Body).Decode(&responseUser)
	expect.Equal("example@admin.com", responseUser.Email)
	expect.Equal(dbUser.ID.Hex(), responseUser.ID.Hex())
}

// endregion
