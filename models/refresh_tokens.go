package models

import (

"crazy_nl_backend/helpers"

"github.com/globalsign/mgo/bson"

)

type RefreshToken struct {
	Token string `json:"token"`
	User string `json:"user"`
}
func (RefreshToken) FindOne(token string, db helpers.IDatabase) *RefreshToken {
	refreshToken := new(RefreshToken)
	db.C("refresh_tokens").Find(bson.M{
		"token": token,
	}).One(&refreshToken)
	if refreshToken.Token == "" {
		return nil
	}
	return refreshToken
}
