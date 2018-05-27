package jwt

import (
	"fmt"
	"net/http"
	"errors"
	"time"

	"github.com/fossapps/starter/db"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

// Client Creates and Parses Jwt
type Client struct {
	Config Config
}

// Manager for creating and parsing JWT
type Manager interface {
	GetJwtDataFromRequest(r *http.Request) (*Claims, error)
	CreateForUser(user *db.User) (string, error)
}

// Config for handling JWT
type Config struct {
	Secret string
	Expiry time.Duration
}

// Claims stored or to be stored in JWT
type Claims struct {
	Email       string   `json:"email"`
	ID          string   `json:"id"`
	Permissions []string `json:"permissions"`
	jwt.StandardClaims
}

func (j Client) signingFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(j.Config.Secret), nil
}

// GetJwtDataFromRequest takes in request and returns claims associated with that request
func (j Client) GetJwtDataFromRequest(r *http.Request) (*Claims, error) {
	var claims Claims
	token, parseErr := request.ParseFromRequestWithClaims(r, request.AuthorizationHeaderExtractor, &claims, j.signingFunc)
	err := claims.Valid()
	if parseErr != nil {
		return nil, parseErr
	}
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return &claims, nil
}

// CreateForUser returns signed JWT token
func (j Client) CreateForUser(user *db.User) (string, error) {
	claims := jwt.MapClaims{
		"id":          user.ID,
		"email":       user.Email,
		"permissions": user.Permissions,
		"exp":         time.Now().Add(j.Config.Expiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(j.Config.Secret))
}
