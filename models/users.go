package models

import (
	"crazy_nl_backend/helpers"
	"github.com/globalsign/mgo/bson"
)

type User struct {
	Email string `json:"email"`
	Password string `json:"password"`
	Permissions []string `json:"permissions"`
}

func (User) FindUserByEmail(email string, db helpers.IDatabase) *User {
	user := new(User)
	db.C("users").Find(bson.M{
		"email": email,
	}).One(&user)
	return user
}