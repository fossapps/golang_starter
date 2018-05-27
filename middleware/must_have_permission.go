package middleware

import (
	"net/http"
	"gopkg.in/matryer/respond.v1"
	"github.com/fossapps/starter/jwt"
)

// MustHavePermission is a middleware which ensures a request has permission before handler is invoked
func MustHavePermission(permission string, manager jwt.Manager) Middleware {
	return func(handler http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			claims, err := manager.GetJwtDataFromRequest(r)
			// if user has sudo, skip permission checking
			if err != nil {
				respond.With(w, r, http.StatusForbidden, struct {
					Message string `json:"message"`
				}{
					Message: "permission denied",
				})
				return
			}
			if claims.Permissions[0] != "sudo" && !contains(claims.Permissions, permission) {
				respond.With(w, r, http.StatusForbidden, struct {
					Message string `json:"message"`
				}{
					Message: "permission denied",
				})
				return
			}
			handler.ServeHTTP(w, r)
		}
	}
}

func contains(collection []string, item string) bool {
	for _, value := range collection {
		if value == item {
			return true
		}
	}
	return false
}
