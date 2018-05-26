package transformer

import (
	"github.com/fossapps/starter/db"
	"github.com/globalsign/mgo/bson"
)

// ResponseUser response friendly version of db.User
type ResponseUser struct {
	ID          bson.ObjectId `json:"id"`
	Email       string        `json:"email"`
	Permissions []string      `json:"permissions"`
}

// TransformUser takes db data and returns response friendly data
func TransformUser(user db.User) ResponseUser {
	return ResponseUser{
		ID:          user.ID,
		Permissions: user.Permissions,
		Email:       user.Email,
	}
}

// TransformUsers takes slice of db data and returns response friendly data
func TransformUsers(items []db.User) []ResponseUser {
	var listItems []ResponseUser
	for _, data := range items {
		listItems = append(listItems, TransformUser(data))
	}
	return listItems
}
