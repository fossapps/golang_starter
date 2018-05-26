package middleware

import "net/http"

// Middleware should be returned from a func to be used as handler (a.k.a middlewares)
type Middleware func(http.Handler) http.HandlerFunc

// Adapt is used to wrap a http handler with Middleware
func Adapt(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, mw := range middlewares {
		h = mw(h)
	}
	return h
}
