package crazy_nl_backend

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"time"

	"crazy_nl_backend/config"
	"crazy_nl_backend/models"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/matryer/respond.v1"
	"errors"
	"crazy_nl_backend/helpers"
	"encoding/json"
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
		session := s.Mongo.Clone()
		defer session.Close()
		db := session.DB(config.GetMongoConfig().DbName)
		user := models.User{}.FindUserByEmail(email, db)
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
		db.C("refresh_tokens").Insert(models.RefreshToken{
			Token: res.RefreshToken,
			User:  user.ID.Hex()})
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
		session := s.Mongo.Clone()
		db := session.DB(config.GetMongoConfig().DbName)
		refreshToken := models.RefreshToken{}.FindOne(token, db)
		if refreshToken == nil {
			s.ErrorResponse(w, r, http.StatusUnauthorized, "refresh token invalid")
			return
		}
		user := models.User{}.FindUserById(refreshToken.User, db)
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

func getJwtForUser(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":          user.ID,
		"email":       user.Email,
		"permissions": user.Permissions,
		"exp": time.Now().Add(config.GetApplicationConfig().JWTExpiryTime).Unix(),
	})
	return token.SignedString([]byte(config.GetApplicationConfig().JWTSecret))
}

func (registration *NewRegistration) OK () error {
	if len(registration.Token) < 20 {
		return errors.New("registration token invalid")
	}
	return nil
}

func (s *Server) RegisterHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		registration, err := decodeRegistrationToken(r)
		if err != nil {
			s.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}

		err = helpers.QueueDeviceRegistration(registration.Token, s.Redis)

		if err != nil {
			s.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}
		respond.With(w, r, 200, struct {
			Status string `json:"status"`
		}{
			Status: "registration pending",
		})
	})
}

func decodeRegistrationToken(r *http.Request) (*NewRegistration, error) {
	registration := new(NewRegistration)
	json.NewDecoder(r.Body).Decode(&registration)
	err := registration.OK()

	if err != nil {
		return nil, err
	}

	return registration, nil
}
