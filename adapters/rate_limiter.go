package adapters

import (
	"net/http"
	"time"
	"gopkg.in/matryer/respond.v1"
	"crazy_nl_backend/helpers"
	"strconv"
)

type IRequestHelper interface {
	GetJwtData(r *http.Request) (*Claims, error)
	GetIpAddress(r *http.Request) string
}

type IRateLimiter interface {
	Hit(key string, duration time.Duration) (int64, error)
	Count(key string) (int64, error)
}

type LimiterOptions struct {
	Namespace     string
	Decay         time.Duration
	RequestHelper IRequestHelper
	Limit         int
	RedisClient   helpers.RedisClient
	AddHeaders    bool
	Logger        ILogger
}

type ILogger interface {
	Warn(args ...interface{})
}

func Limit(options LimiterOptions) Adapter {
	return func(handler http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			limiter := helpers.Limiter{
				RedisClient: options.RedisClient,
				Decay:       options.Decay,
				Limit:       options.Limit,
			}
			key := getKeyFromRequest(options.Namespace, r, options.RequestHelper)

			card, err := limiter.Count(key)
			if err != nil {
				options.Logger.Warn("rate limiting counting error", err)
			}

			if card+1 > int64(options.Limit) {
				if options.AddHeaders {
					options.addHeaders(w, card)
				}
				respond.With(w, r, http.StatusTooManyRequests, "too many requests")
				return
			}
			card, err = limiter.Hit(key)
			if err != nil {
				options.Logger.Warn("rate limit error", err)
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

func getKeyFromRequest(namespace string, r *http.Request, requestHelper IRequestHelper) string {
	data, err := requestHelper.GetJwtData(r)
	var key string
	if err != nil {
		key = requestHelper.GetIpAddress(r)
	} else {
		key = data.ID
	}
	return namespace + "-" + key
}
