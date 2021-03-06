package middleware

import (
	"net/http"
	"time"
)

// ResponseTimeLogger interface implementation needed for ResponseTime for logging
type ResponseTimeLogger interface {
	Info(args ...interface{})
}

// ResponseTime middleware to log response time for a request
func ResponseTime(logger ResponseTimeLogger) Middleware {
	return func(handler http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			defer logger.Info("response in ", time.Since(now).Nanoseconds())
			handler.ServeHTTP(w, r)
		}
	}
}
