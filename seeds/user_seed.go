package seeds

import (
	"crazy_nl_backend/helpers"
	"crazy_nl_backend/models"
	"fmt"
)

type UserSeed struct{}

func (UserSeed) GetKey() string {
	return "USER_SEED"
}

func (UserSeed) GetDescription() string {
	return "Create default users"
}

func (UserSeed) Seed(db helpers.IDatabase) {
	admin := models.User{
		Email:       "admin@example.com",
		Password:    "admin1234",
		Permissions: []string{"sudo"},
	}
	err := db.C("users").Insert(admin)
	if err != nil {
		fmt.Println(err)
	}
}

func (UserSeed) Remove(db helpers.IDatabase) {

}
