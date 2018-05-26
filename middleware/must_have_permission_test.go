package middleware_test

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/fossapps/starter/middleware"
	"github.com/fossapps/starter/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/matryer/respond.v1"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func getFakeJWTWithPermission(permissions []string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"permissions": permissions,
		"exp":         time.Now().Add(10 * time.Minute).Unix(),
	})
	return token.SignedString([]byte(config.GetApplicationConfig().JWTSecret))
}

func TestMustHavePermissionRespondsWithStatusUnauthorized(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond.With(w, r, http.StatusNotImplemented, nil)
	})
	responseRecorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	token, err := getFakeJWTWithPermission([]string{"user.create"})
	assert.Nil(t, err)
	assert.NotNil(t, token)
	req.Header.Add("Authorization", "Bearer "+token)
	middleware.MustHavePermission("user.destroy")(handler)(responseRecorder, req)
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
}

func TestMustHavePermissionLetsHttpHandlerSetStatusCodeIfHavePermission(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond.With(w, r, http.StatusNotImplemented, nil)
	})
	responseRecorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	token, err := getFakeJWTWithPermission([]string{"user.create"})
	assert.Nil(t, err)
	assert.NotNil(t, token)
	req.Header.Add("Authorization", "Bearer "+token)
	middleware.MustHavePermission("user.create")(handler)(responseRecorder, req)
	assert.Equal(t, http.StatusNotImplemented, responseRecorder.Code)
}

func TestMustHavePermissionStopsInvalidJWT(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond.With(w, r, http.StatusNotImplemented, nil)
	})
	responseRecorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Bearer "+"some.random.string")
	middleware.MustHavePermission("user.create")(handler)(responseRecorder, req)
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
}

func TestMustHavePermissionLetsSudoPermissionThrough(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond.With(w, r, http.StatusNotImplemented, nil)
	})
	responseRecorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	token, err := getFakeJWTWithPermission([]string{"sudo"})
	assert.Nil(t, err)
	assert.NotNil(t, token)
	req.Header.Add("Authorization", "Bearer "+token)
	middleware.MustHavePermission("user.create")(handler)(responseRecorder, req)
	assert.Equal(t, http.StatusNotImplemented, responseRecorder.Code)
}
