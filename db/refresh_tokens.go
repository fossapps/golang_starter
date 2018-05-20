package db

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// RefreshToken representation of a Refresh Token
type RefreshToken struct {
	Token string `json:"token"`
	User  string `json:"user"`
}

// IRefreshTokenManager deals with persistence of RefreshTokens
type IRefreshTokenManager interface {
	FindOne(token string) *RefreshToken
	Add(token string, user string)
	List(user string) ([]RefreshToken, error)
	Delete(token string) error
}

type refreshTokenLayer struct {
	db *mgo.Database
}

// Add a refresh token for a user
func (dbLayer refreshTokenLayer) Add(token string, user string) {
	dbLayer.db.C("refresh_tokens").Insert(RefreshToken{
		Token: token,
		User:  user,
	})
}

// FindOne using token
func (dbLayer refreshTokenLayer) FindOne(token string) *RefreshToken {
	refreshToken := new(RefreshToken)
	dbLayer.db.C("refresh_tokens").Find(bson.M{
		"token": token,
	}).One(&refreshToken)
	if refreshToken.Token == "" {
		return nil
	}
	return refreshToken
}

// List all tokens belonging to user
func (dbLayer refreshTokenLayer) List(user string) ([]RefreshToken, error) {
	var tokens []RefreshToken
	err := dbLayer.db.C("refresh_tokens").Find(bson.M{
		"user": user,
	}).All(&tokens)
	return tokens, err
}

func (dbLayer refreshTokenLayer) Delete(token string) error {
	return dbLayer.db.C("refresh_tokens").Remove(bson.M{
		"token": token,
	})
}

// GetRefreshTokenManager returns implementation of IRefreshTokenManager
func GetRefreshTokenManager(db *mgo.Database) IRefreshTokenManager {
	return refreshTokenLayer{
		db: db,
	}
}
