//go:generate mockgen -destination=./mocks/mock_redis.go -package=mocks crazy_nl_backend/helpers IRedisClient

package crazy_nl_backend_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"crazy_nl_backend"
	"crazy_nl_backend/db"
	"crazy_nl_backend/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"github.com/sirupsen/logrus"
)

func getLogger() crazy_nl_backend.ILogger {
	logger := logrus.New()
	logger.Out = httptest.NewRecorder()
	return logger
}
func TestServer_LoginHandlerRespondsWithUnauthorizedIfNoHeader(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	server := crazy_nl_backend.Server{
		Logger: getLogger(),
	}
	server.LoginHandler()(responseRecorder, request)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
}

func TestServer_LoginHandlerRespondsWithUnauthorizedIfWrongPassword(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	mockUsersCtrl := gomock.NewController(t)
	mockDb := mocks.NewMockDb(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockUserManager := mocks.NewMockIUserManager(mockUsersCtrl)
	email := "admin@example.com"
	pass := "pass"
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	user := db.User{Email: email, Password: string(hash)}
	mockUserManager.EXPECT().FindByEmail("admin@example.com").Return(&user)
	mockDb.EXPECT().Users().Times(1).Return(mockUserManager)
	server := crazy_nl_backend.Server{
		Logger: getLogger(),
		Db: mockDb,
	}

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.SetBasicAuth(email, "wrong_password")
	server.LoginHandler()(responseRecorder, request)
	expect.Equal(http.StatusUnauthorized, responseRecorder.Code)
}

func TestServer_LoginHandlerRespondsWithOkOnCorrectCredentials(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	mockUsersCtrl := gomock.NewController(t)
	refreshTokenCtrl := gomock.NewController(t)
	mockRefreshTokenManager := mocks.NewMockIRefreshTokenManager(refreshTokenCtrl)
	mockRefreshTokenManager.EXPECT().Add(gomock.Any(), gomock.Any())
	mockDb := mocks.NewMockDb(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockDb.EXPECT().RefreshTokens().Times(1).Return(mockRefreshTokenManager)
	mockUserManager := mocks.NewMockIUserManager(mockUsersCtrl)
	email := "admin@example.com"
	pass := "pass"
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	user := db.User{Email: email, Password: string(hash)}
	mockUserManager.EXPECT().FindByEmail("admin@example.com").Return(&user)
	mockDb.EXPECT().Users().Times(1).Return(mockUserManager)
	server := crazy_nl_backend.Server{
		Db: mockDb,
	}

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.SetBasicAuth(email, pass)
	server.LoginHandler()(responseRecorder, request)
	expect.Equal(http.StatusOK, responseRecorder.Code)
	res := new(crazy_nl_backend.LoginResponse)
	json.NewDecoder(responseRecorder.Body).Decode(&res)
	expect.NotNil(res.RefreshToken)
	expect.NotNil(res.JWT)
	expect.True(len(res.RefreshToken) >= 128)
	expect.True(strings.Count(res.JWT, ".") == 2)
}

func TestServer_RefreshTokenHandlerStoresRefreshTokenInDb(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	refreshTokenCtrl := gomock.NewController(t)
	mockUsersCtrl := gomock.NewController(t)
	mockRefreshTokenManager := mocks.NewMockIRefreshTokenManager(refreshTokenCtrl)
	mockDb := mocks.NewMockDb(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockUserManager := mocks.NewMockIUserManager(mockUsersCtrl)
	email := "admin@example.com"
	pass := "pass"
	mockRefreshTokenManager.EXPECT().Add(gomock.Any(), gomock.Any())
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	user := db.User{Email: email, Password: string(hash)}
	mockUserManager.EXPECT().FindByEmail(email).Return(&user)
	mockDb.EXPECT().Users().Times(1).Return(mockUserManager)
	mockDb.EXPECT().RefreshTokens().Times(1).Return(mockRefreshTokenManager)
	server := crazy_nl_backend.Server{
		Db: mockDb,
	}

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.SetBasicAuth(email, pass)
	server.LoginHandler()(responseRecorder, request)
	expect.Equal(http.StatusOK, responseRecorder.Code)
	res := new(crazy_nl_backend.LoginResponse)
	json.NewDecoder(responseRecorder.Body).Decode(&res)
	expect.NotNil(res.RefreshToken)
	expect.NotNil(res.JWT)
	expect.True(len(res.RefreshToken) >= 128)
	expect.True(strings.Count(res.JWT, ".") == 2)
}

func TestServer_RefreshTokenHandlerRespondsWithStatusBadRequestIfNoAuthToken(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	server := crazy_nl_backend.Server{}
	server.RefreshTokenHandler()(responseRecorder, request)
	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
}

func TestServer_RefreshTokenHandlerRespondsWithStatusUnauthorizedIfRefreshTokenInvalid(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	refreshTokenCtrl := gomock.NewController(t)
	mockUsersCtrl := gomock.NewController(t)
	mockUserManager := mocks.NewMockIUserManager(mockUsersCtrl)
	mockRefreshTokenManager := mocks.NewMockIRefreshTokenManager(refreshTokenCtrl)
	mockDb := mocks.NewMockDb(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockDb.EXPECT().Users().Times(1).Return(mockUserManager)
	mockRefreshTokenManager.EXPECT().FindOne("auth_token").Times(1).Return(nil)
	mockDb.EXPECT().RefreshTokens().Times(1).Return(mockRefreshTokenManager)

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.Header.Add("Authorization", "Bearer auth_token")
	server := crazy_nl_backend.Server{
		Db: mockDb,
		Logger: getLogger(),
	}
	server.RefreshTokenHandler()(responseRecorder, request)
	expect.Equal(http.StatusUnauthorized, responseRecorder.Code)
}

func TestServer_RefreshTokenHandlerRefreshTokenNotLinkedToUserRespondsWithStatusUnauthorized(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	refreshTokenCtrl := gomock.NewController(t)
	mockUsersCtrl := gomock.NewController(t)
	mockUserManager := mocks.NewMockIUserManager(mockUsersCtrl)
	mockRefreshTokenManager := mocks.NewMockIRefreshTokenManager(refreshTokenCtrl)
	mockDb := mocks.NewMockDb(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockUserManager.EXPECT().FindById("some_user").Times(1).Return(nil)
	mockDb.EXPECT().Users().Times(1).Return(mockUserManager)
	mockToken := &db.RefreshToken{
		Token:"auth_token",
		User: "some_user",
	}
	mockRefreshTokenManager.EXPECT().FindOne("auth_token").Times(1).Return(mockToken)
	mockDb.EXPECT().RefreshTokens().Times(1).Return(mockRefreshTokenManager)

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.Header.Add("Authorization", "Bearer auth_token")

	server := crazy_nl_backend.Server{
		Db: mockDb,
		Logger: getLogger(),
	}
	server.RefreshTokenHandler()(responseRecorder, request)
	expect.Equal(http.StatusUnauthorized, responseRecorder.Code)
}

func TestServer_RefreshTokenHandlerReturnsJWT(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	refreshTokenCtrl := gomock.NewController(t)
	mockUsersCtrl := gomock.NewController(t)
	mockUserManager := mocks.NewMockIUserManager(mockUsersCtrl)
	mockRefreshTokenManager := mocks.NewMockIRefreshTokenManager(refreshTokenCtrl)
	mockDb := mocks.NewMockDb(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockUser := &db.User{
		ID: "random",
		Email: "random",
		Permissions: []string{"sudo"},
	}
	mockUserManager.EXPECT().FindById("some_user").Times(1).Return(mockUser)
	mockDb.EXPECT().Users().Times(1).Return(mockUserManager)
	mockToken := &db.RefreshToken{
		Token: "auth_token",
		User: "some_user",
	}
	mockRefreshTokenManager.EXPECT().FindOne("auth_token").Times(1).Return(mockToken)
	mockDb.EXPECT().RefreshTokens().Times(1).Return(mockRefreshTokenManager)

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.Header.Add("Authorization", "Bearer auth_token")

	server := crazy_nl_backend.Server{
		Db: mockDb,
		Logger: getLogger(),
	}
	server.RefreshTokenHandler()(responseRecorder, request)
	expect.Equal(http.StatusOK, responseRecorder.Code)
}
