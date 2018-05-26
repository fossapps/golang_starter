package adapter_test

import (
	"fmt"
	"github.com/fossapps/starter/adapter"
	"github.com/stretchr/testify/assert"
	"gopkg.in/matryer/respond.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ExampleAdapt() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handler function")
	})
	newHandler := adapter.Adapt(handler, getTestAdapter())
	newHandler(httptest.NewRecorder(), nil)
	// Output:
	// before
	// handler function
	// after
}

func TestAdaptLetsAdapterWrapHandlers(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond.With(w, r, http.StatusNotImplemented, nil)
	})
	newHandler := adapter.Adapt(handler, getTestAdapter())
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	newHandler(responseRecorder, request)
	assert.Equal(t, http.StatusNotImplemented, responseRecorder.Code)
	responseRecorder = httptest.NewRecorder()
	adapter.Adapt(handler, getTestAdapterWithBlockHandler())(responseRecorder, request)
	assert.NotEqual(t, http.StatusNotImplemented, responseRecorder.Code)
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func getTestAdapterWithBlockHandler() adapter.Adapter {
	return func(handler http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func getTestAdapter() adapter.Adapter {
	return func(handler http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("before")
			defer fmt.Println("after")
			handler.ServeHTTP(w, r)
		}
	}
}
