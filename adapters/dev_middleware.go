package adapters

import (
	"math/rand"
	"net/http"
	"time"
)

// DevMw simple middleware to delay a request
func DevMw(milliseconds int) Adapter {
	return func(handler http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			duration := time.Millisecond * time.Duration(rand.Intn(milliseconds))
			time.Sleep(duration)
			handler.ServeHTTP(w, r)
		}
	}
}
