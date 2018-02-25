package adapters

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"time"
)

func ResponseTime(logger logrus.Logger) Adapter {
	return func(handler http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			defer logger.Info("response in ", time.Since(now).Nanoseconds())
			handler.ServeHTTP(w, r)
		}
	}
}
