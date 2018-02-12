package adapters

import "net/http"

type Adapter func(http.Handler) http.HandlerFunc

func Adapt(h http.HandlerFunc, adapters ...Adapter) http.HandlerFunc {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}
