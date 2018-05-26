package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fossapps/starter/middleware"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestResponseTimeTakesInLoggerReturnsAMiddleware(t *testing.T) {
	logger := logrus.New()
	mw := middleware.ResponseTime(logger)
	Assert := assert.New(t)
	Assert.NotNil(mw)
}

func TestResponseTimeAdapterTakesInHandlerReturnsHandler(t *testing.T) {
	logger := logrus.New()
	mw := middleware.ResponseTime(logger)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	newHandler := mw(handler)
	assert.IsType(t, handler, newHandler)
}

func TestResponseTimeLogsResponseTime(t *testing.T) {
	logger := logrus.New()
	loggerOutput := httptest.NewRecorder()
	logger.Out = loggerOutput
	mw := middleware.ResponseTime(logger)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	newHandler := mw(handler)
	newHandler(httptest.NewRecorder(), nil)
	assert.Contains(t, loggerOutput.Body.String(), "response in")
	assert.Contains(t, loggerOutput.Body.String(), "level=info")
}
