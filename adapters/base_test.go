package adapters_test

import (
	"net/http"
	"crazy_nl_backend/adapters"
	"net/http/httptest"
	"fmt"
)

func ExampleAdapt() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		fmt.Println("handler function")
	})
	newHandler := adapters.Adapt(handler, getTestAdapter())
	newHandler(httptest.NewRecorder(), nil)
	// Output:
	// before
	// handler function
	// after
}

func getTestAdapter() adapters.Adapter {
	return func(handler http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("before")
			defer fmt.Println("after")
			handler.ServeHTTP(w, r)
		}
	}
}
