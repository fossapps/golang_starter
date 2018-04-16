package crazy_nl_backend

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"crazy_nl_backend/config"
	"crazy_nl_backend/db"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/matryer/respond.v1"
)

type LoginResponse struct {
	JWT          string `json:"jwt"`
	RefreshToken string `json:"refresh_token"`
}

type NewRegistration struct {
	Token string `json:"token"`
}

type RefreshTokenHandlerResponse struct {
	Token string `json:"token"`
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
			s.Logger.Error("token generation error: ", err)
			s.ErrorResponse(w, r, http.StatusInternalServerError, "error generating token")
			return
		}
		respond.With(w, r, http.StatusOK, RefreshTokenHandlerResponse{Token: token})
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

func (registration *NewRegistration) OK() bool {
	if len(registration.Token) < 20 {
		return false
	}
	return true
}

type RegistrationResponse struct {
	Status string `json:"status"`
}

func (s *Server) RegisterHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		registration := new(NewRegistration)
		err := json.NewDecoder(r.Body).Decode(&registration)
		if err != nil {
			s.ErrorResponse(w, r, http.StatusBadRequest, "invalid request")
			return
		}
		if !registration.OK() {
			s.ErrorResponse(w, r, http.StatusBadRequest, "invalid token")
			return
		}

		// save to mongodb
		if s.Db.Devices().Exists(registration.Token) {
			s.ErrorResponse(w, r, http.StatusBadRequest, "already registered")
			return
		}
		err = s.Db.Devices().Register(registration.Token)

		if err != nil {
			s.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}
		respond.With(w, r, 200, RegistrationResponse{
			Status: "success",
		})
	})
}
