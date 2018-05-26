package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fossapps/starter/middleware"
	"github.com/stretchr/testify/assert"
	"gopkg.in/matryer/respond.v1"
)

func ExampleAdapt() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handler function")
	})
	newHandler := middleware.Adapt(handler, getTestAdapter())
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
	newHandler := middleware.Adapt(handler, getTestAdapter())
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	newHandler(responseRecorder, request)
	assert.Equal(t, http.StatusNotImplemented, responseRecorder.Code)
	responseRecorder = httptest.NewRecorder()
	middleware.Adapt(handler, getTestAdapterWithBlockHandler())(responseRecorder, request)
	assert.NotEqual(t, http.StatusNotImplemented, responseRecorder.Code)
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func getTestAdapterWithBlockHandler() middleware.Middleware {
	return func(handler http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func getTestAdapter() middleware.Middleware {
	return func(handler http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("before")
			defer fmt.Println("after")
			handler.ServeHTTP(w, r)
		}
	}
}
