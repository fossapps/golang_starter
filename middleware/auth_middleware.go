package middleware

import (
	"net/http"

	"gopkg.in/matryer/respond.v1"
)

type jwtHelper interface {
	GetJwtData(r *http.Request) (*Claims, error)
}

// AuthMw ensures request is authenticated
func AuthMw(helper jwtHelper) Middleware {
	return func(handler http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			data, err := helper.GetJwtData(r)
			if err != nil || data == nil {
				respond.With(w, r, http.StatusUnauthorized, nil)
				return
			}
			handler.ServeHTTP(w, r)
		}
	}
}
