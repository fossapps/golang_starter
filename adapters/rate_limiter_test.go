package adapters_test

import (
	"testing"
	"github.com/golang/mock/gomock"
	"crazy_nl_backend/mocks"
	"crazy_nl_backend/adapters"
	"net/http"
	"gopkg.in/matryer/respond.v1"
	"errors"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
)

// region setup data

func getMockRequestHelper(t *testing.T) *mocks.MockIRequestHelper {
	ctrl := gomock.NewController(t)
	return mocks.NewMockIRequestHelper(ctrl)
}

func getMockLogger(t *testing.T) *mocks.MockILogger {
	mockLogger := mocks.NewMockILogger(gomock.NewController(t))
	mockLogger.EXPECT().Warn(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warning(gomock.Any()).AnyTimes()
	return mockLogger
}

func getMockRateLimiter(t *testing.T) *mocks.MockIRateLimiter {
	return mocks.NewMockIRateLimiter(gomock.NewController(t))
}

func getTestHandler(status int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond.With(w, r, status, "success")
	}
}

func getLimiterOptions(t *testing.T) adapters.LimiterOptions {
	return adapters.LimiterOptions{
		Limit:     5,
		Namespace: "my_key",
		Logger:    getMockLogger(t),
	}
}

// endregion

func TestLimitUsesUserIdIfAvailable(t *testing.T) {
	limiterOptions := getLimiterOptions(t)
	mockRequestHelper := getMockRequestHelper(t)
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).Times(1).Return(&adapters.Claims{ID: "my_id"}, nil)
	mockRateLimiter := getMockRateLimiter(t)
	mockRateLimiter.EXPECT().Count("my_key-my_id").Times(1).Return(int64(0), errors.New("error"))
	limiterOptions.RequestHelper = mockRequestHelper
	limiterOptions.Limiter = mockRateLimiter
	handler := adapters.Adapt(getTestHandler(http.StatusOK), adapters.Limit(limiterOptions))
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	handler(responseRecorder, request)
}

func TestLimitUsesIpAddrIfIdNotAvailable(t *testing.T) {
	limiterOptions := getLimiterOptions(t)
	mockRequestHelper := getMockRequestHelper(t)
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).Times(1).Return(nil, errors.New("error"))
	mockRequestHelper.EXPECT().GetIpAddress(gomock.Any()).Times(1).Return("ip_addr")
	mockRateLimiter := getMockRateLimiter(t)
	mockRateLimiter.EXPECT().Count("my_key-ip_addr").Times(1).Return(int64(0), errors.New("error"))
	limiterOptions.RequestHelper = mockRequestHelper
	limiterOptions.Limiter = mockRateLimiter
	handler := adapters.Adapt(getTestHandler(http.StatusOK), adapters.Limit(limiterOptions))
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	handler(responseRecorder, request)
}

func TestLimitReturnsInternalServerErrorIfCountNull(t *testing.T) {
	expect := assert.New(t)
	limiterOptions := getLimiterOptions(t)
	mockRequestHelper := getMockRequestHelper(t)
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).Times(1).Return(&adapters.Claims{ID: "my_id"}, nil)
	mockRateLimiter := getMockRateLimiter(t)
	mockRateLimiter.EXPECT().Count("my_key-my_id").Times(1).Return(int64(0), errors.New("error"))
	limiterOptions.RequestHelper = mockRequestHelper
	limiterOptions.Limiter = mockRateLimiter
	handler := adapters.Adapt(getTestHandler(http.StatusOK), adapters.Limit(limiterOptions))
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	handler(responseRecorder, request)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
}

func TestLimitReturnsTooManyRequestIfCountGreaterThanOrEqualToLimit(t *testing.T) {
	expect := assert.New(t)
	limiterOptions := getLimiterOptions(t)
	mockRequestHelper := getMockRequestHelper(t)
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).Times(1).Return(&adapters.Claims{ID: "my_id"}, nil)
	mockRateLimiter := getMockRateLimiter(t)
	mockRateLimiter.EXPECT().Count("my_key-my_id").Times(1).Return(int64(5), nil)
	limiterOptions.RequestHelper = mockRequestHelper
	limiterOptions.Limiter = mockRateLimiter
	limiterOptions.AddHeaders = false
	handler := adapters.Adapt(getTestHandler(http.StatusOK), adapters.Limit(limiterOptions))
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	handler(responseRecorder, request)
	expect.Equal(http.StatusTooManyRequests, responseRecorder.Code)
}

func TestLimitReturnsTooManyRequestIfCountGreaterThanOrEqualToLimitAndAddsHeaderIfNeeded(t *testing.T) {
	expect := assert.New(t)
	limiterOptions := getLimiterOptions(t)
	mockRequestHelper := getMockRequestHelper(t)
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).Times(1).Return(&adapters.Claims{ID: "my_id"}, nil)
	mockRateLimiter := getMockRateLimiter(t)
	mockRateLimiter.EXPECT().Count("my_key-my_id").Times(1).Return(int64(5), nil)
	limiterOptions.RequestHelper = mockRequestHelper
	limiterOptions.Limiter = mockRateLimiter
	limiterOptions.AddHeaders = true
	handler := adapters.Adapt(getTestHandler(http.StatusOK), adapters.Limit(limiterOptions))
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
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).Times(1).Return(&adapters.Claims{ID: "my_id"}, nil)
	mockRateLimiter := getMockRateLimiter(t)
	mockRateLimiter.EXPECT().Count("my_key-my_id").Times(1).Return(int64(0), nil)
	mockRateLimiter.EXPECT().Hit("my_key-my_id").Times(1).Return(int64(0), errors.New("error"))
	limiterOptions.RequestHelper = mockRequestHelper
	limiterOptions.Limiter = mockRateLimiter
	limiterOptions.AddHeaders = false
	handler := adapters.Adapt(getTestHandler(http.StatusOK), adapters.Limit(limiterOptions))
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	handler(responseRecorder, request)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
}
func TestLimitCallsHandler(t *testing.T) {
	expect := assert.New(t)
	limiterOptions := getLimiterOptions(t)
	mockRequestHelper := getMockRequestHelper(t)
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).Times(1).Return(&adapters.Claims{ID: "my_id"}, nil)
	mockRateLimiter := getMockRateLimiter(t)
	mockRateLimiter.EXPECT().Count("my_key-my_id").Times(1).Return(int64(0), nil)
	mockRateLimiter.EXPECT().Hit("my_key-my_id").Times(1).Return(int64(0), nil)
	limiterOptions.RequestHelper = mockRequestHelper
	limiterOptions.Limiter = mockRateLimiter
	limiterOptions.AddHeaders = false
	handler := adapters.Adapt(getTestHandler(http.StatusAccepted), adapters.Limit(limiterOptions))
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	handler(responseRecorder, request)
	expect.Equal(http.StatusAccepted, responseRecorder.Code)
}

func TestLimitCallsHandlerAndSetsHeaderIfRequested(t *testing.T) {
	expect := assert.New(t)
	limiterOptions := getLimiterOptions(t)
	mockRequestHelper := getMockRequestHelper(t)
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).Times(1).Return(&adapters.Claims{ID: "my_id"}, nil)
	mockRateLimiter := getMockRateLimiter(t)
	mockRateLimiter.EXPECT().Count("my_key-my_id").Times(1).Return(int64(0), nil)
	mockRateLimiter.EXPECT().Hit("my_key-my_id").Times(1).Return(int64(0), nil)
	limiterOptions.RequestHelper = mockRequestHelper
	limiterOptions.Limiter = mockRateLimiter
	limiterOptions.AddHeaders = true
	handler := adapters.Adapt(getTestHandler(http.StatusAccepted), adapters.Limit(limiterOptions))
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	handler(responseRecorder, request)
	expect.Equal(http.StatusAccepted, responseRecorder.Code)
	expect.NotNil(responseRecorder.Header().Get("X-RateLimit-Limit"))
	expect.NotNil(responseRecorder.Header().Get("X-RateLimit-Remaining"))
}
