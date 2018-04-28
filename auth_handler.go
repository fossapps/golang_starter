package crazy_nl_backend

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"time"

	"crazy_nl_backend/config"
	"crazy_nl_backend/db"

	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/matryer/respond.v1"
)

type LoginResponse struct {
	JWT          string `json:"jwt"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenHandlerResponse struct {
	JWT string `json:"jwt"`
}

func (s *Server) LoginHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, password, ok := r.BasicAuth()
		if !ok {
			s.Logger.Warn("basic auth not present")
			s.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
			return
		}
		dbLayer := s.Db.Clone()
		defer dbLayer.Close()
		user := dbLayer.Users().FindByEmail(email)
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			s.Logger.Warn("username and password wrong")
			s.ErrorResponse(w, r, http.StatusUnauthorized, "username/password mismatch")
			return
		}
		jwtToken, err := getJwtForUser(user)
		if err != nil {
			// since it's a separate thing, maybe application init test should take care of it
			// then we don't need to check if it had an error?
			// todo can't test as of now, let's see.
			s.Logger.Error("jwt error", err)
			s.ErrorResponse(w, r, http.StatusInternalServerError, "error generating token")
			return
		}
		refreshToken := getRefreshToken(config.GetApplicationConfig().RefreshTokenSize)
		res := LoginResponse{
			JWT:          jwtToken,
			RefreshToken: refreshToken,
		}
		dbLayer.RefreshTokens().Add(res.RefreshToken, user.ID.Hex())
		respond.With(w, r, http.StatusOK, res)
	})
}

func (s *Server) RefreshTokenHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			s.ErrorResponse(w, r, http.StatusBadRequest, "token missing")
			return
		}
		dbLayer := s.Db.Clone()
		defer dbLayer.Close()
		refreshToken := dbLayer.RefreshTokens().FindOne(token)
		if refreshToken == nil {
			s.ErrorResponse(w, r, http.StatusUnauthorized, "refresh token invalid")
			return
		}
		user := dbLayer.Users().FindById(refreshToken.User)
		if user == nil {
			s.ErrorResponse(w, r, http.StatusUnauthorized, "invalid refresh token")
			return
		}
		token, err := getJwtForUser(user)
		if err != nil {
			// since it's a separate thing, maybe application init test should take care of it
			// then we don't need to check if it had an error?
			// todo can't test as of now, let's see. (same comment in above handler, good idea to extract)
			s.Logger.Error("token generation error: ", err)
			s.ErrorResponse(w, r, http.StatusInternalServerError, "error generating token")
			return
		}
		respond.With(w, r, http.StatusOK, RefreshTokenHandlerResponse{JWT: token})
	})
}

func (s *Server) RefreshTokensList() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := s.ReqHelper.GetJwtData(r)
		if err != nil {
			s.ErrorResponse(w, r, http.StatusBadRequest, "error parsing token")
			return
		}
		user := claims.ID
		database := s.Db.Clone()
		tokens, err := database.RefreshTokens().List(user)
		if err != nil {
			s.ErrorResponse(w, r, http.StatusInternalServerError, "database error")
			return
		}
		respond.With(w, r, http.StatusOK, tokens)
	})
}

func (s *Server) DeleteSession() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := mux.Vars(r)["token"]
		err := s.Db.RefreshTokens().Delete(token)
		if err != nil && err.Error() == mgo.ErrNotFound.Error() {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			s.ErrorResponse(w, r, http.StatusInternalServerError, "db error")
			return
		}
		respond.With(w, r, http.StatusNoContent, nil)
	})
}

func getRefreshToken(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func getJwtForUser(user *db.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":          user.ID,
		"email":       user.Email,
		"permissions": user.Permissions,
		"exp":         time.Now().Add(config.GetApplicationConfig().JWTExpiryTime).Unix(),
	})
	return token.SignedString([]byte(config.GetApplicationConfig().JWTSecret))
}
