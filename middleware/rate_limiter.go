package middleware

import (
	"net/http"
	"strconv"

	"github.com/fossapps/starter/jwt"
	"gopkg.in/matryer/respond.v1"
)

// RequestHelper see RequestHelper in server.go
type RequestHelper interface {
	GetIPAddress(r *http.Request) string
}

// RateLimiter implementation needed to work
type RateLimiter interface {
	Hit(key string) (int64, error)
	Count(key string) (int64, error)
}

// LimiterOptions dependencies of Limit middleware
type LimiterOptions struct {
	Namespace     string
	RequestHelper RequestHelper
	Limit         int
	AddHeaders    bool
	Logger        Logger
	Limiter       RateLimiter
	Jwt           jwt.Manager
}

// Logger implementation needed for Limit middleware to log
type Logger interface {
	Warn(args ...interface{})
}

// Limit middleware ensures a handler isn't called more than N times in D duraton
func Limit(options LimiterOptions) Middleware {
	return func(handler http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			key := getKeyFromRequest(options.Namespace, r, options.RequestHelper, options.Jwt)

			card, err := options.Limiter.Count(key)
			if err != nil {
				options.Logger.Warn("rate limiting counting error", err)
				respond.With(w, r, http.StatusInternalServerError, "server error")
				return
			}

			if card >= int64(options.Limit) {
				if options.AddHeaders {
					options.addHeaders(w, card)
				}
				respond.With(w, r, http.StatusTooManyRequests, "too many requests")
				return
			}
			card, err = options.Limiter.Hit(key)
			if err != nil {
				options.Logger.Warn("rate limit error", err)
				respond.With(w, r, http.StatusInternalServerError, "server error")
				return
			}
			if options.AddHeaders {
				options.addHeaders(w, card)
			}
			handler.ServeHTTP(w, r)
		}
	}
}

func (limiter LimiterOptions) addHeaders(w http.ResponseWriter, currentCount int64) {
	remaining := int64(limiter.Limit) - currentCount
	w.Header().Add("X-RateLimit-Limit", strconv.Itoa(limiter.Limit))
	w.Header().Add("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
}

func getKeyFromRequest(namespace string, r *http.Request, requestHelper RequestHelper, jwtClient jwt.Manager) string {
	data, err := jwtClient.GetJwtDataFromRequest(r)
	if err != nil {
		return namespace + "-" + requestHelper.GetIPAddress(r)
	}
	return namespace + "-" + data.ID
}
