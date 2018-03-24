package migrations

import (
	"crazy_nl_backend/models"
	"fmt"
	"github.com/globalsign/mgo"
	"golang.org/x/crypto/bcrypt"
)

type UserSeed struct{}

func (UserSeed) GetKey() string {
	return "USER_SEED"
}

func (UserSeed) GetDescription() string {
	return "Create default users"
}

func (UserSeed) Apply(db *mgo.Database) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin1234"), bcrypt.DefaultCost)
	admin := models.User{
		Email:       "admin@example.com",
		Password:    string(hash),
		Permissions: []string{"sudo"},
	}
	err := db.C("users").Insert(admin)
	if err != nil {
		fmt.Println(err)
	}
}

func (UserSeed) Remove(db *mgo.Database) {

}
