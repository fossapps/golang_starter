package adapters

import (
	"net/http"
	"time"
)
type IResponseTimeLogger interface {
	Info(args ...interface{})
}

func ResponseTime(logger IResponseTimeLogger) Adapter {
	return func(handler http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			defer logger.Info("response in ", time.Since(now).Nanoseconds())
			handler.ServeHTTP(w, r)
		}
	}
}
