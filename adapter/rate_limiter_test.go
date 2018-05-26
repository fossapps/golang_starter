package adapter_test

import (
	"errors"
	"github.com/fossapps/starter/adapter"
	"github.com/fossapps/starter/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/matryer/respond.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

// region setup data

func getMockRequestHelper(t *testing.T) *mock.MockRequestHelper {
	ctrl := gomock.NewController(t)
	return mock.NewMockRequestHelper(ctrl)
}

func getMockLogger(t *testing.T) *mock.MockLogger {
	mockLogger := mock.NewMockLogger(gomock.NewController(t))
	mockLogger.EXPECT().Warn(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warning(gomock.Any()).AnyTimes()
	return mockLogger
}

func getMockRateLimiter(t *testing.T) *mock.MockRateLimiter {
	return mock.NewMockRateLimiter(gomock.NewController(t))
}

func getTestHandler(status int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond.With(w, r, status, "success")
	}
}

func getLimiterOptions(t *testing.T) adapter.LimiterOptions {
	return adapter.LimiterOptions{
		Limit:     5,
		Namespace: "my_key",
		Logger:    getMockLogger(t),
	}
}

// endregion

func TestLimitUsesUserIdIfAvailable(t *testing.T) {
	limiterOptions := getLimiterOptions(t)
	mockRequestHelper := getMockRequestHelper(t)
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).Times(1).Return(&adapter.Claims{ID: "my_id"}, nil)
	mockRateLimiter := getMockRateLimiter(t)
	mockRateLimiter.EXPECT().Count("my_key-my_id").Times(1).Return(int64(0), errors.New("error"))
	limiterOptions.RequestHelper = mockRequestHelper
	limiterOptions.Limiter = mockRateLimiter
	handler := adapter.Adapt(getTestHandler(http.StatusOK), adapter.Limit(limiterOptions))
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	handler(responseRecorder, request)
}

func TestLimitUsesIpAddrIfIdNotAvailable(t *testing.T) {
	limiterOptions := getLimiterOptions(t)
	mockRequestHelper := getMockRequestHelper(t)
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).Times(1).Return(nil, errors.New("error"))
	mockRequestHelper.EXPECT().GetIPAddress(gomock.Any()).Times(1).Return("ip_addr")
	mockRateLimiter := getMockRateLimiter(t)
	mockRateLimiter.EXPECT().Count("my_key-ip_addr").Times(1).Return(int64(0), errors.New("error"))
	limiterOptions.RequestHelper = mockRequestHelper
	limiterOptions.Limiter = mockRateLimiter
	handler := adapter.Adapt(getTestHandler(http.StatusOK), adapter.Limit(limiterOptions))
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	handler(responseRecorder, request)
}

func TestLimitReturnsInternalServerErrorIfCountNull(t *testing.T) {
	expect := assert.New(t)
	limiterOptions := getLimiterOptions(t)
	mockRequestHelper := getMockRequestHelper(t)
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).Times(1).Return(&adapter.Claims{ID: "my_id"}, nil)
	mockRateLimiter := getMockRateLimiter(t)
	mockRateLimiter.EXPECT().Count("my_key-my_id").Times(1).Return(int64(0), errors.New("error"))
	limiterOptions.RequestHelper = mockRequestHelper
	limiterOptions.Limiter = mockRateLimiter
	handler := adapter.Adapt(getTestHandler(http.StatusOK), adapter.Limit(limiterOptions))
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	handler(responseRecorder, request)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
}

func TestLimitReturnsTooManyRequestIfCountGreaterThanOrEqualToLimit(t *testing.T) {
	expect := assert.New(t)
	limiterOptions := getLimiterOptions(t)
	mockRequestHelper := getMockRequestHelper(t)
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).Times(1).Return(&adapter.Claims{ID: "my_id"}, nil)
	mockRateLimiter := getMockRateLimiter(t)
	mockRateLimiter.EXPECT().Count("my_key-my_id").Times(1).Return(int64(5), nil)
	limiterOptions.RequestHelper = mockRequestHelper
	limiterOptions.Limiter = mockRateLimiter
	limiterOptions.AddHeaders = false
	handler := adapter.Adapt(getTestHandler(http.StatusOK), adapter.Limit(limiterOptions))
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	handler(responseRecorder, request)
	expect.Equal(http.StatusTooManyRequests, responseRecorder.Code)
}

func TestLimitReturnsTooManyRequestIfCountGreaterThanOrEqualToLimitAndAddsHeaderIfNeeded(t *testing.T) {
	expect := assert.New(t)
	limiterOptions := getLimiterOptions(t)
	mockRequestHelper := getMockRequestHelper(t)
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).Times(1).Return(&adapter.Claims{ID: "my_id"}, nil)
	mockRateLimiter := getMockRateLimiter(t)
	mockRateLimiter.EXPECT().Count("my_key-my_id").Times(1).Return(int64(5), nil)
	limiterOptions.RequestHelper = mockRequestHelper
	limiterOptions.Limiter = mockRateLimiter
	limiterOptions.AddHeaders = true
	handler := adapter.Adapt(getTestHandler(http.StatusOK), adapter.Limit(limiterOptions))
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	handler(responseRecorder, request)
	expect.Equal(http.StatusTooManyRequests, responseRecorder.Code)
	expect.NotNil(responseRecorder.Header().Get("X-RateLimit-Limit"))
	expect.NotNil(responseRecorder.Header().Get("X-RateLimit-Remaining"))
}

func TestLimitHandlesHitError(t *testing.T) {
	expect := assert.New(t)
	limiterOptions := getLimiterOptions(t)
	mockRequestHelper := getMockRequestHelper(t)
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).Times(1).Return(&adapter.Claims{ID: "my_id"}, nil)
	mockRateLimiter := getMockRateLimiter(t)
	mockRateLimiter.EXPECT().Count("my_key-my_id").Times(1).Return(int64(0), nil)
	mockRateLimiter.EXPECT().Hit("my_key-my_id").Times(1).Return(int64(0), errors.New("error"))
	limiterOptions.RequestHelper = mockRequestHelper
	limiterOptions.Limiter = mockRateLimiter
	limiterOptions.AddHeaders = false
	handler := adapter.Adapt(getTestHandler(http.StatusOK), adapter.Limit(limiterOptions))
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	handler(responseRecorder, request)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
}
func TestLimitCallsHandler(t *testing.T) {
	expect := assert.New(t)
	limiterOptions := getLimiterOptions(t)
	mockRequestHelper := getMockRequestHelper(t)
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).Times(1).Return(&adapter.Claims{ID: "my_id"}, nil)
	mockRateLimiter := getMockRateLimiter(t)
	mockRateLimiter.EXPECT().Count("my_key-my_id").Times(1).Return(int64(0), nil)
	mockRateLimiter.EXPECT().Hit("my_key-my_id").Times(1).Return(int64(0), nil)
	limiterOptions.RequestHelper = mockRequestHelper
	limiterOptions.Limiter = mockRateLimiter
	limiterOptions.AddHeaders = false
	handler := adapter.Adapt(getTestHandler(http.StatusAccepted), adapter.Limit(limiterOptions))
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	handler(responseRecorder, request)
	expect.Equal(http.StatusAccepted, responseRecorder.Code)
}

func TestLimitCallsHandlerAndSetsHeaderIfRequested(t *testing.T) {
	expect := assert.New(t)
	limiterOptions := getLimiterOptions(t)
	mockRequestHelper := getMockRequestHelper(t)
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).Times(1).Return(&adapter.Claims{ID: "my_id"}, nil)
	mockRateLimiter := getMockRateLimiter(t)
	mockRateLimiter.EXPECT().Count("my_key-my_id").Times(1).Return(int64(0), nil)
	mockRateLimiter.EXPECT().Hit("my_key-my_id").Times(1).Return(int64(0), nil)
	limiterOptions.RequestHelper = mockRequestHelper
	limiterOptions.Limiter = mockRateLimiter
	limiterOptions.AddHeaders = true
	handler := adapter.Adapt(getTestHandler(http.StatusAccepted), adapter.Limit(limiterOptions))
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	handler(responseRecorder, request)
	expect.Equal(http.StatusAccepted, responseRecorder.Code)
	expect.NotNil(responseRecorder.Header().Get("X-RateLimit-Limit"))
	expect.NotNil(responseRecorder.Header().Get("X-RateLimit-Remaining"))
}
