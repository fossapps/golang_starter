package crazy_nl_backend

import (
	"net/http"
	"gopkg.in/matryer/respond.v1"
	"crazy_nl_backend/config"
	"crazy_nl_backend/models"
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"crypto/rand"
	"github.com/dgrijalva/jwt-go"
	"time"
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
			s.Logger.Warn("jwt error", err)
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
			User:  user.ID.String()})
		respond.With(w, r, http.StatusOK, res)
	})
}

func getRefreshToken(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func getJwtForUser(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email":       user.Email,
		"permissions": user.Permissions,
	})
	token.Header["exp"] = time.Now().Add(config.GetApplicationConfig().JWTExpiryTime).UTC().Unix()
	return token.SignedString([]byte(config.GetApplicationConfig().JWTSecret))
}
