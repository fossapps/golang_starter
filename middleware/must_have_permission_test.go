package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fossapps/starter/jwt"
	"github.com/fossapps/starter/middleware"
	"github.com/fossapps/starter/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/matryer/respond.v1"
)

func getFakeJWTWithPermission(permissions []string) string {
	return "token: " + strings.Join(permissions, ",")
}

func TestMustHavePermissionRespondsWithStatusUnauthorized(t *testing.T) {
	mockJwtCtrl := gomock.NewController(t)
	defer mockJwtCtrl.Finish()
	mockJwt := mock.NewMockJwtManager(mockJwtCtrl)
	mockJwt.EXPECT().GetJwtDataFromRequest(gomock.Any()).AnyTimes().Return(nil, errors.New("err"))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond.With(w, r, http.StatusNotImplemented, nil)
	})
	responseRecorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	token := getFakeJWTWithPermission([]string{"user.create"})
	assert.NotNil(t, token)
	req.Header.Add("Authorization", "Bearer "+token)
	middleware.MustHavePermission("user.destroy", mockJwt)(handler)(responseRecorder, req)
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
}

func TestMustHavePermissionLetsHttpHandlerSetStatusCodeIfHavePermission(t *testing.T) {
	mockJwtCtrl := gomock.NewController(t)
	defer mockJwtCtrl.Finish()
	mockJwt := mock.NewMockJwtManager(mockJwtCtrl)
	mockJwt.EXPECT().GetJwtDataFromRequest(gomock.Any()).AnyTimes().Return(&jwt.Claims{Permissions: []string{"user.create"}}, nil)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond.With(w, r, http.StatusNotImplemented, nil)
	})
	responseRecorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	token := getFakeJWTWithPermission([]string{"user.create"})
	assert.NotNil(t, token)
	req.Header.Add("Authorization", "Bearer "+token)
	middleware.MustHavePermission("user.create", mockJwt)(handler)(responseRecorder, req)
	assert.Equal(t, http.StatusNotImplemented, responseRecorder.Code)
}

func TestMustHavePermissionStopsInvalidJWT(t *testing.T) {
	mockJwtCtrl := gomock.NewController(t)
	defer mockJwtCtrl.Finish()
	mockJwt := mock.NewMockJwtManager(mockJwtCtrl)
	mockJwt.EXPECT().GetJwtDataFromRequest(gomock.Any()).AnyTimes().Return(nil, errors.New("error"))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond.With(w, r, http.StatusNotImplemented, nil)
	})
	responseRecorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Bearer "+"some.random.string")
	middleware.MustHavePermission("user.create", mockJwt)(handler)(responseRecorder, req)
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
}

func TestMustHavePermissionLetsSudoPermissionThrough(t *testing.T) {
	mockJwtCtrl := gomock.NewController(t)
	defer mockJwtCtrl.Finish()
	mockJwt := mock.NewMockJwtManager(mockJwtCtrl)
	mockJwt.EXPECT().GetJwtDataFromRequest(gomock.Any()).AnyTimes().Return(&jwt.Claims{Permissions: []string{"sudo"}}, nil)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond.With(w, r, http.StatusNotImplemented, nil)
	})
	responseRecorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	token := getFakeJWTWithPermission([]string{"sudo"})
	assert.NotNil(t, token)
	req.Header.Add("Authorization", "Bearer "+token)
	middleware.MustHavePermission("user.create", mockJwt)(handler)(responseRecorder, req)
	assert.Equal(t, http.StatusNotImplemented, responseRecorder.Code)
}
