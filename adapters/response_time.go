package adapters

import (
	"net/http"
	"time"
)

// IResponseTimeLogger interface implementation needed for ResponseTime for logging
type IResponseTimeLogger interface {
	Info(args ...interface{})
}

// ResponseTime adapter to log response time for a request
func ResponseTime(logger IResponseTimeLogger) Adapter {
	return func(handler http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			defer logger.Info("response in ", time.Since(now).Nanoseconds())
			handler.ServeHTTP(w, r)
		}
	}
}
