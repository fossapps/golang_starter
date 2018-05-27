package middleware

import (
	"net/http"

	"github.com/fossapps/starter/jwt"
	"gopkg.in/matryer/respond.v1"
)

// AuthMw ensures request is authenticated
func AuthMw(jwtManager jwt.Manager) Middleware {
	return func(handler http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			data, err := jwtManager.GetJwtDataFromRequest(r)
			if err != nil || data == nil {
				respond.With(w, r, http.StatusUnauthorized, nil)
				return
			}
			handler.ServeHTTP(w, r)
		}
	}
}
