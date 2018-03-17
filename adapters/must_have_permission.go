package adapters

import (
"errors"
"net/http"

"crazy_nl_backend/config"

"github.com/dgrijalva/jwt-go"
"github.com/dgrijalva/jwt-go/request"
"gopkg.in/matryer/respond.v1"
	"fmt"
)

type Claims struct {
	Email string `json:"email"`
	ID string `json:"id"`
	Permissions []string `json:"permissions"`
	jwt.StandardClaims
}

func signingFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
	}
	return []byte(config.GetApplicationConfig().JWTSecret), nil
}

func MustHavePermission(permission string) Adapter {
	return func(handler http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var claims Claims
			token, parseErr := request.ParseFromRequestWithClaims(r, request.AuthorizationHeaderExtractor, &claims, signingFunc)
			// if user has sudo, skip permission checking
			err := claims.Valid()
			if parseErr != nil || err != nil || !token.Valid {
				respond.With(w, r, http.StatusUnauthorized, struct {
					Message string `json:"message"`
				}{
					Message: "permission denied",
				})
				return
			}
			if claims.Permissions[0] != "sudo" && !contains(claims.Permissions, permission) {
				respond.With(w, r, http.StatusUnauthorized, struct {
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
