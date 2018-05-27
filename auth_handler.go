package starter

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
		"github.com/fossapps/starter/config"
			"github.com/fossapps/starter/transformer"
	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/matryer/respond.v1"
)

// LoginResponse responds with this type when login is successful
type LoginResponse struct {
	JWT          string `json:"jwt"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenHandlerResponse is used when a refresh token is requested
type RefreshTokenHandlerResponse struct {
	JWT string `json:"jwt"`
}

// LoginHandler handles login requests
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
		user, err := dbLayer.Users().FindByEmail(email)
		if err != nil {
			s.ErrorResponse(w, r, http.StatusInternalServerError, "server error")
			return
		}
		if user == nil {
			s.ErrorResponse(w, r, http.StatusBadRequest, "invalid credentials")
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			s.ErrorResponse(w, r, http.StatusUnauthorized, "invalid credentials")
			return
		}
		jwtToken, err := s.Jwt.CreateForUser(user)
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

// RefreshTokenHandler is used to refresh token
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
		refreshToken, err := dbLayer.RefreshTokens().FindOne(token)
		if err != nil {
			s.ErrorResponse(w, r, http.StatusInternalServerError, "server error")
			return
		}
		if refreshToken == nil {
			s.ErrorResponse(w, r, http.StatusUnauthorized, "invalid refresh token")
			return
		}
		user, err := dbLayer.Users().FindByID(refreshToken.User)
		if err != nil {
			s.ErrorResponse(w, r, http.StatusInternalServerError, "server error")
			return
		}
		if user == nil {
			// todo shouldn't have happened, log
			s.Logger.Warn("refresh token found in collection, but user is nil: ", refreshToken)
			s.ErrorResponse(w, r, http.StatusUnauthorized, "invalid refresh token")
			return
		}
		token, err = s.Jwt.CreateForUser(user)
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

// RefreshTokensList returns list of refresh token associated with a user
func (s *Server) RefreshTokensList() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := s.Jwt.GetJwtDataFromRequest(r)
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
		respond.With(w, r, http.StatusOK, transformer.TransformRefreshTokens(tokens))
	})
}

// DeleteSession deletes a session, used for logging out
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
