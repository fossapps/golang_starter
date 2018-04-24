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

//
//func TestServer_RefreshTokenHandlerRespondsWithStatusBadRequestIfNoAuthToken(t *testing.T) {
//	expect := assert.New(t)
//	responseRecorder := httptest.NewRecorder()
//	request := httptest.NewRequest("POST", "/", nil)
//	server := crazy_nl_backend.Server{}
//	server.RefreshTokenHandler()(responseRecorder, request)
//	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
//}
//
//func TestServer_RefreshTokenHandlerRespondsWithStatusUnauthorizedIfRefreshTokenInvalid(t *testing.T) {
//	expect := assert.New(t)
//	responseRecorder := httptest.NewRecorder()
//	request := httptest.NewRequest("POST", "/", nil)
//	request.Header.Add("Authorization", "Bearer invalid_bearer_token")
//	server := crazy_nl_backend.Server{
//		Mongo:  getSession(),
//		Logger: *logrus.New(),
//	}
//	server.RefreshTokenHandler()(responseRecorder, request)
//	expect.Equal(http.StatusUnauthorized, responseRecorder.Code)
//}
//
//func TestServer_RefreshTokenHandlerRefreshTokenNotLinkedToUserRespondsWithStatusUnauthorized(t *testing.T) {
//	expect := assert.New(t)
//	responseRecorder := httptest.NewRecorder()
//	request := httptest.NewRequest("POST", "/", nil)
//	request.Header.Add("Authorization", "Bearer random_token")
//	session := getSession()
//	session.DB(config.GetMongoConfig().DbName).DropDatabase()
//	session.DB(config.GetMongoConfig().DbName).C("refresh_tokens").Insert(models.RefreshToken{
//		User:  "aaaaaaaaaaaaaaaaaaaaaaaa",
//		Token: "random_token",
//	})
//	server := crazy_nl_backend.Server{
//		Mongo:  session,
//		Logger: *logrus.New(),
//	}
//	server.RefreshTokenHandler()(responseRecorder, request)
//	expect.Equal(http.StatusUnauthorized, responseRecorder.Code)
//}
//
//func TestServer_RefreshTokenHandlerReturnsJWT(t *testing.T) {
//	expect := assert.New(t)
//	responseRecorder := httptest.NewRecorder()
//	request := httptest.NewRequest("POST", "/", nil)
//	request.Header.Add("Authorization", "Bearer random_token")
//	session := getSession()
//	session.DB(config.GetMongoConfig().DbName).DropDatabase()
//	session.DB(config.GetMongoConfig().DbName).C("refresh_tokens").Insert(models.RefreshToken{
//		User:  "bbbbbbbbbbbbbbbbbbbbbbbb",
//		Token: "random_token",
//	})
//	session.DB(config.GetMongoConfig().DbName).C("users").Insert(models.User{
//		ID: bson.ObjectIdHex("bbbbbbbbbbbbbbbbbbbbbbbb"),
//	})
//	server := crazy_nl_backend.Server{
//		Mongo:  session,
//		Logger: *logrus.New(),
//	}
//	response := crazy_nl_backend.RefreshTokenHandlerResponse{}
//	server.RefreshTokenHandler()(responseRecorder, request)
//	json.NewDecoder(responseRecorder.Body).Decode(&response)
//	expect.Equal(http.StatusOK, responseRecorder.Code)
//	expect.NotNil(response.Token)
//	assert.True(t, strings.Count(response.Token, ".") == 2)
//}
//
//func TestServer_RegisterHandlerRespondsWithStatusBadRequestIfBodyDoesNotContainToken(t *testing.T) {
//	expect := assert.New(t)
//	responseRecorder := httptest.NewRecorder()
//	request := httptest.NewRequest("POST", "/", nil)
//	server := crazy_nl_backend.Server{}
//	server.RegisterHandler()(responseRecorder, request)
//	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
//}
//
//func TestServer_RegisterHandlerRespondsWithStatusBadRequestIfTokenIsTooShort(t *testing.T) {
//	expect := assert.New(t)
//	responseRecorder := httptest.NewRecorder()
//	buffer := new(bytes.Buffer)
//	json.NewEncoder(buffer).Encode(crazy_nl_backend.NewRegistration{
//		Token: "token",
//	})
//	request := httptest.NewRequest("POST", "/", buffer)
//	server := crazy_nl_backend.Server{}
//	server.RegisterHandler()(responseRecorder, request)
//	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
//	expect.Contains(responseRecorder.Body.String(), "registration token invalid")
//}
//
//func TestServer_RegisterHandlerRespondsWithStatusBadRequestTokenAlreadyExists(t *testing.T) {
//	expect := assert.New(t)
//	responseRecorder := httptest.NewRecorder()
//	controller := gomock.NewController(t)
//	mockRedis := mocks.NewMockIRedisClient(controller)
//	buffer := new(bytes.Buffer)
//	token := "very_large_token_here"
//	mockRedis.EXPECT().SIsMember("registration", token).Times(1).Return(true, nil)
//	json.NewEncoder(buffer).Encode(crazy_nl_backend.NewRegistration{
//		Token: token,
//	})
//	request := httptest.NewRequest("POST", "/", buffer)
//	server := crazy_nl_backend.Server{
//		Redis: mockRedis,
//	}
//	server.RegisterHandler()(responseRecorder, request)
//	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
//	expect.Contains(responseRecorder.Body.String(), "token already exists")
//}
//
//func TestServer_RegisterHandlerQueuesRegistrationIfTokenOk(t *testing.T) {
//	expect := assert.New(t)
//	responseRecorder := httptest.NewRecorder()
//	controller := gomock.NewController(t)
//	mockRedis := mocks.NewMockIRedisClient(controller)
//	buffer := new(bytes.Buffer)
//	token := "very_large_token_here"
//	mockRedis.EXPECT().SIsMember("registration", token).Times(1).Return(false, nil)
//	mockRedis.EXPECT().SAdd("registration", token).Times(1).Return(int64(1), nil)
//	json.NewEncoder(buffer).Encode(crazy_nl_backend.NewRegistration{
//		Token: token,
//	})
//	request := httptest.NewRequest("POST", "/", buffer)
//	server := crazy_nl_backend.Server{
//		Redis: mockRedis,
//	}
//	server.RegisterHandler()(responseRecorder, request)
//	expect.Equal(http.StatusOK, responseRecorder.Code)
//	result := &crazy_nl_backend.RegistrationResponse{}
//	json.NewDecoder(responseRecorder.Body).Decode(&result)
//	expect.Equal(result.Status, "registration pending")
//}
