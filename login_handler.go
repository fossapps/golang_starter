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
)

type LoginResponse struct {
	JWT          string `json:"jwt"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Server) LoginHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, password, ok := r.BasicAuth()
		if !ok {
			s.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
			return
		}
		session := s.Mongo.Clone()
		defer session.Close()
		db := session.DB(config.GetMongoConfig().DbName)
		user := models.User{}.FindUserByEmail(email, db)
		// s.Logger.Info(user.ID.String(), user.Email)
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
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
		db.C("refresh_tokens").Insert(struct {
			Token string `json:"token"`
			User  string `json:"user"`
		}{
			Token: res.RefreshToken,
			User:  user.ID.Hex()})
		respond.With(w, r, http.StatusOK, res)
	})
}

func (s *Server) RefreshTokenHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == ""{
			s.ErrorResponse(w, r, http.StatusBadRequest, "token missing")
			return
		}
		session := s.Mongo.Clone()
		db := session.DB(config.GetMongoConfig().DbName)
		refreshToken := models.RefreshToken{}.FindOne(token, db)
		s.Logger.Info(refreshToken)
		if refreshToken == nil {
			s.ErrorResponse(w, r, http.StatusUnauthorized, "refresh token invalid")
			return
		}
		user := models.User{}.FindUserById(refreshToken.User, db)
		if user == nil {
			s.ErrorResponse(w, r, http.StatusUnauthorized, "invalid refresh token")
			s.Logger.Error("user should not have been nil, refreshToken: ", refreshToken)
			return
		}
		token, err := getJwtForUser(user)
		if err != nil {
			s.Logger.Error("token generation error: ", err)
			s.ErrorResponse(w, r, http.StatusInternalServerError, "error generating token")
			return
		}
		respond.With(w, r, http.StatusOK, struct {
			Token string `json:"token"`
		}{Token:token})
	})
}

func getRefreshToken(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func getJwtForUser(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id": user.ID,
		"email":       user.Email,
		"permissions": user.Permissions,
	})
	token.Header["exp"] = time.Now().Add(config.GetApplicationConfig().JWTExpiryTime).UTC().Unix()
	return token.SignedString([]byte(config.GetApplicationConfig().JWTSecret))
}
