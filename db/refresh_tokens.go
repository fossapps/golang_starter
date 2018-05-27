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

// RefreshTokenManager deals with persistence of RefreshTokens
type RefreshTokenManager interface {
	FindOne(token string) (*RefreshToken, error)
	Add(token string, user string) error
	List(user string) ([]RefreshToken, error)
	Delete(token string) error
}

type refreshTokenManager struct {
	db *mgo.Database
}

// Add a refresh token for a user
func (dbLayer refreshTokenManager) Add(token string, user string) error {
	return dbLayer.db.C("refresh_tokens").Insert(RefreshToken{
		Token: token,
		User:  user,
	})
}

// FindOne using token
func (dbLayer refreshTokenManager) FindOne(token string) (*RefreshToken, error) {
	refreshToken := new(RefreshToken)
	err := dbLayer.db.C("refresh_tokens").Find(bson.M{
		"token": token,
	}).One(&refreshToken)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if refreshToken.Token == "" {
		return nil, nil
	}
	return refreshToken, nil
}

// List all tokens belonging to user
func (dbLayer refreshTokenManager) List(user string) ([]RefreshToken, error) {
	var tokens []RefreshToken
	err := dbLayer.db.C("refresh_tokens").Find(bson.M{
		"user": user,
	}).All(&tokens)
	return tokens, err
}

func (dbLayer refreshTokenManager) Delete(token string) error {
	return dbLayer.db.C("refresh_tokens").Remove(bson.M{
		"token": token,
	})
}

// GetRefreshTokenManager returns implementation of RefreshTokenManager
func GetRefreshTokenManager(db *mgo.Database) RefreshTokenManager {
	return refreshTokenManager{
		db: db,
	}
}
