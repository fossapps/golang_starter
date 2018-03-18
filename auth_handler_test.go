//go:generate mockgen -destination=./mocks/mock_redis.go -package=mocks crazy_nl_backend/helpers IRedisClient

package crazy_nl_backend_test

import (
	"testing"
	"net/http/httptest"
	"crazy_nl_backend"
	"github.com/stretchr/testify/assert"
	"net/http"
	"crazy_nl_backend/migrations"
	"github.com/globalsign/mgo"
	"crazy_nl_backend/config"
	"os"
	"encoding/json"
	"strings"
	"crazy_nl_backend/models"
	"github.com/sirupsen/logrus"
	"github.com/globalsign/mgo/bson"
	"bytes"
	"crazy_nl_backend/mocks"
	"github.com/golang/mock/gomock"
)

func TestMain(m *testing.M) {
	session, _ := mgo.Dial(config.GetTestingDbConnection())
	migrations.ApplyAll(config.GetMongoConfig().DbName, session)
	defer session.DB(config.GetMongoConfig().DbName).DropDatabase()
	result := m.Run()
	os.Exit(result)
}

func getSession() *mgo.Session {
	session, err := mgo.Dial(config.GetTestingDbConnection())
	if err != nil {
		panic(err)
	}
	return session
}

func TestServer_LoginHandlerRespondsWithUnauthorizedIfNoHeader(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	server := crazy_nl_backend.Server{}
	server.LoginHandler()(responseRecorder, request)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
}

func TestServer_LoginHandlerRespondsWithUnauthorizedIfWrongPassword(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.SetBasicAuth("admin@example.com", "pass")
	session := getSession()
	assert.NotNil(t, session)
	server := crazy_nl_backend.Server{
		Mongo:session,
	}
	server.LoginHandler()(responseRecorder, request)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
}

func TestServer_LoginHandlerRespondsWithOkOnCorrectCredentials(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.SetBasicAuth("admin@example.com", "admin1234")
	session := getSession()
	assert.NotNil(t, session)
	server := crazy_nl_backend.Server{
		Mongo:session,
	}
	server.LoginHandler()(responseRecorder, request)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestServer_LoginHandlerRespondsWithTwoTokensOnCorrectCredentials(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.SetBasicAuth("admin@example.com", "admin1234")
	session := getSession()
	assert.NotNil(t, session)
	server := crazy_nl_backend.Server{
		Mongo:session,
	}
	server.LoginHandler()(responseRecorder, request)
	res := new(crazy_nl_backend.LoginResponse)
	json.NewDecoder(responseRecorder.Body).Decode(&res)
	assert.NotNil(t, res.RefreshToken)
	assert.NotNil(t, res.JWT)
	// we always use at least 64 bit, it's hexadecimal, so 64 bits give 128 chars
	assert.True(t, len(res.RefreshToken) >= 128)
	assert.True(t, strings.Count(res.JWT, ".") ==2)
}

func TestServer_RefreshTokenHandlerStoresRefreshTokenInDb(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.SetBasicAuth("admin@example.com", "admin1234")
	session := getSession()
	assert.NotNil(t, session)
	server := crazy_nl_backend.Server{
		Mongo:session,
	}
	server.LoginHandler()(responseRecorder, request)
	res := new(crazy_nl_backend.LoginResponse)
	json.NewDecoder(responseRecorder.Body).Decode(&res)
	refreshToken := models.RefreshToken{}.FindOne(res.RefreshToken, session.DB(config.GetMongoConfig().DbName))
	assert.NotNil(t, refreshToken.Token)
	assert.Equal(t, res.RefreshToken, refreshToken.Token)
	user := models.User{}.FindUserById(refreshToken.User, session.DB(config.GetMongoConfig().DbName))
	assert.Equal(t, "admin@example.com", user.Email)
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
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.Header.Add("Authorization", "Bearer invalid_bearer_token")
	server := crazy_nl_backend.Server{
		Mongo:getSession(),
		Logger:*logrus.New(),
	}
	server.RefreshTokenHandler()(responseRecorder, request)
	expect.Equal(http.StatusUnauthorized, responseRecorder.Code)
}

func TestServer_RefreshTokenHandlerRefreshTokenNotLinkedToUserRespondsWithStatusUnauthorized(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.Header.Add("Authorization", "Bearer random_token")
	session := getSession()
	session.DB(config.GetMongoConfig().DbName).DropDatabase()
	session.DB(config.GetMongoConfig().DbName).C("refresh_tokens").Insert(models.RefreshToken{
		User:"aaaaaaaaaaaaaaaaaaaaaaaa",
		Token:"random_token",
	})
	server := crazy_nl_backend.Server{
		Mongo:session,
		Logger:*logrus.New(),
	}
	server.RefreshTokenHandler()(responseRecorder, request)
	expect.Equal(http.StatusUnauthorized, responseRecorder.Code)
}

func TestServer_RefreshTokenHandlerReturnsJWT(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.Header.Add("Authorization", "Bearer random_token")
	session := getSession()
	session.DB(config.GetMongoConfig().DbName).DropDatabase()
	session.DB(config.GetMongoConfig().DbName).C("refresh_tokens").Insert(models.RefreshToken{
		User:"bbbbbbbbbbbbbbbbbbbbbbbb",
		Token:"random_token",
	})
	session.DB(config.GetMongoConfig().DbName).C("users").Insert(models.User{
		ID:bson.ObjectIdHex("bbbbbbbbbbbbbbbbbbbbbbbb"),
	})
	server := crazy_nl_backend.Server{
		Mongo:session,
		Logger:*logrus.New(),
	}
	response := crazy_nl_backend.RefreshTokenHandlerResponse{}
	server.RefreshTokenHandler()(responseRecorder, request)
	json.NewDecoder(responseRecorder.Body).Decode(&response)
	expect.Equal(http.StatusOK, responseRecorder.Code)
	expect.NotNil(response.Token)
	assert.True(t, strings.Count(response.Token, ".") == 2)
}

func TestServer_RegisterHandlerRespondsWithStatusBadRequestIfBodyDoesNotContainToken(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	server := crazy_nl_backend.Server{}
	server.RegisterHandler()(responseRecorder, request)
	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
}

func TestServer_RegisterHandlerRespondsWithStatusBadRequestIfTokenIsTooShort(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(crazy_nl_backend.NewRegistration{
		Token:"token",
	})
	request := httptest.NewRequest("POST", "/", buffer)
	server := crazy_nl_backend.Server{}
	server.RegisterHandler()(responseRecorder, request)
	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
	expect.Contains(responseRecorder.Body.String(), "registration token invalid")
}

func TestServer_RegisterHandlerRespondsWithStatusBadRequestTokenAlreadyExists(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	controller := gomock.NewController(t)
	mockRedis := mocks.NewMockIRedisClient(controller)
	buffer := new(bytes.Buffer)
	token := "very_large_token_here"
	mockRedis.EXPECT().SIsMember("registration", token).Times(1).Return(true, nil)
	json.NewEncoder(buffer).Encode(crazy_nl_backend.NewRegistration{
		Token:token,
	})
	request := httptest.NewRequest("POST", "/", buffer)
	server := crazy_nl_backend.Server{
		Redis:mockRedis,
	}
	server.RegisterHandler()(responseRecorder, request)
	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
	expect.Contains(responseRecorder.Body.String(), "token already exists")
}

func TestServer_RegisterHandlerQueuesRegistrationIfTokenOk(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	controller := gomock.NewController(t)
	mockRedis := mocks.NewMockIRedisClient(controller)
	buffer := new(bytes.Buffer)
	token := "very_large_token_here"
	mockRedis.EXPECT().SIsMember("registration", token).Times(1).Return(false, nil)
	mockRedis.EXPECT().SAdd("registration", token).Times(1).Return(int64(1), nil)
	json.NewEncoder(buffer).Encode(crazy_nl_backend.NewRegistration{
		Token:token,
	})
	request := httptest.NewRequest("POST", "/", buffer)
	server := crazy_nl_backend.Server{
		Redis:mockRedis,
	}
	server.RegisterHandler()(responseRecorder, request)
	expect.Equal(http.StatusOK, responseRecorder.Code)
	expect.Contains(responseRecorder.Body.String(), "registration pending")
}

