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
