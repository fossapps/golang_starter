// +build integration

package crazy_nl_backend_test

import (
	"testing"
	"net/http/httptest"
	"crazy_nl_backend"
	"github.com/stretchr/testify/assert"
	"net/http"
	"crazy_nl_backend/helpers"
	"crazy_nl_backend/config"
	"encoding/json"
	"strings"
)

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
	session, err := helpers.GetMongo(config.GetMongoConfig())
	assert.Nil(t, err)
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
	session, err := helpers.GetMongo(config.GetMongoConfig())
	assert.Nil(t, err)
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
	session, err := helpers.GetMongo(config.GetMongoConfig())
	assert.Nil(t, err)
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
