package adapters_test

import (
	"golang_starter/adapters"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseTimeTakesInLoggerReturnsAnAdapter(t *testing.T) {
	logger := logrus.New()
	adapter := adapters.ResponseTime(logger)
	Assert := assert.New(t)
	Assert.NotNil(adapter)
}

func TestResponseTimeAdapterTakesInHandlerReturnsHandler(t *testing.T) {
	logger := logrus.New()
	adapter := adapters.ResponseTime(logger)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	newHandler := adapter(handler)
	assert.IsType(t, handler, newHandler)
}

func TestResponseTimeLogsResponseTime(t *testing.T) {
	logger := logrus.New()
	loggerOutput := httptest.NewRecorder()
	logger.Out = loggerOutput
	adapter := adapters.ResponseTime(logger)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	newHandler := adapter(handler)
	newHandler(httptest.NewRecorder(), nil)
	assert.Contains(t, loggerOutput.Body.String(), "response in")
	assert.Contains(t, loggerOutput.Body.String(), "level=info")
}
