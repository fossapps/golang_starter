package adapters

import "net/http"

// Adapter should be returned from a func to be used as handler (a.k.a middlewares)
type Adapter func(http.Handler) http.HandlerFunc

// Adapt is used to wrap a http handler with Adapter
func Adapt(h http.HandlerFunc, adapters ...Adapter) http.HandlerFunc {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}
