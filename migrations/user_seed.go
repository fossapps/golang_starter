package migrations

import (
	"golang_starter/db"
	"fmt"

	"github.com/globalsign/mgo"
)

type UserSeed struct{}

func (UserSeed) GetKey() string {
	return "USER_SEED"
}

func (UserSeed) GetDescription() string {
	return "Create default users"
}

func (UserSeed) Apply(dbLayer db.Db) {
	admin := db.User{
		Email:       "admin@example.com",
		Password:    "admin1234",
		Permissions: []string{"sudo"},
	}
	err := dbLayer.Users().Create(admin)
	if err != nil {
		fmt.Println(err)
	}
}

func (UserSeed) Remove(db *mgo.Database) {

}
