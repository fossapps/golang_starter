package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fossapps/starter/middleware"
	"github.com/fossapps/starter/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/matryer/respond.v1"
)

func TestAuthMwBlocksUnauthorizedUsers(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond.With(w, r, http.StatusNotImplemented, nil)
	})
	responseRecorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	ctrl := gomock.NewController(t)
	mockRequestHelper := mock.NewMockRequestHelper(ctrl)
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).AnyTimes().Return(nil, nil)
	middleware.AuthMw(mockRequestHelper)(handler)(responseRecorder, req)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
}

func TestAuthMwLetsAuthorizedRequestPass(t *testing.T) {
	responseStatus := http.StatusNotImplemented
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond.With(w, r, responseStatus, nil)
	})
	responseRecorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	token, _ := getFakeJWTWithPermission([]string{"user.create"})
	req.Header.Add("Authorization", "Bearer "+token)
	ctrl := gomock.NewController(t)
	mockRequestHelper := mock.NewMockRequestHelper(ctrl)
	claims := middleware.Claims{
		Email: "test@example.com",
	}
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).AnyTimes().Return(&claims, nil)
	middleware.AuthMw(mockRequestHelper)(handler)(responseRecorder, req)
	assert.Equal(t, responseStatus, responseRecorder.Code)
}
