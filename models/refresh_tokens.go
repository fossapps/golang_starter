package models

import (
	"github.com/globalsign/mgo/bson"

	"github.com/globalsign/mgo"
)

type RefreshToken struct {
	Token string `json:"token"`
	User  string `json:"user"`
}

func (RefreshToken) FindOne(token string, db *mgo.Database) *RefreshToken {
	refreshToken := new(RefreshToken)
	db.C("refresh_tokens").Find(bson.M{
		"token": token,
	}).One(&refreshToken)
	if refreshToken.Token == "" {
		return nil
	}
	return refreshToken
}
