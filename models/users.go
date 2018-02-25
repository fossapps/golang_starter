package models

import (
	"crazy_nl_backend/helpers"

	"github.com/globalsign/mgo/bson"
)

type User struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Email       string        `json:"email"`
	Password    string        `json:"password"`
	Permissions []string      `json:"permissions"`
}

func (User) FindUserByEmail(email string, db helpers.IDatabase) *User {
	user := new(User)
	db.C("users").Find(bson.M{
		"email": email,
	}).One(&user)
	return user
}

func (User) FindUserById(id string, db helpers.IDatabase) *User {
	user := new(User)
	db.C("users").Find(bson.M{
		"_id": bson.ObjectIdHex(id),
	}).One(&user)
	if user.ID == "" {
		return nil
	}
	return user
}
