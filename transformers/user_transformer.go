package transformers

import (
	"crazy_nl_backend/db"
	"github.com/globalsign/mgo/bson"
)

type ResponseUser struct {
	ID          bson.ObjectId   `json:"id"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
}

func TransformUser(user db.User) ResponseUser {
	return ResponseUser{
		ID: user.ID,
		Permissions: user.Permissions,
		Email: user.Email,
	}
}

func TransformUsers(items []db.User) []ResponseUser {
	var listItems []ResponseUser
	for _, data := range items {
		listItems = append(listItems, TransformUser(data))
	}
	return listItems
}
