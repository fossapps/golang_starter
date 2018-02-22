package seeds

import (
	"crazy_nl_backend/helpers"
	"fmt"
)

type UserSeed struct{}

func (UserSeed) GetKey() string {
	return "USER_SEED"
}

func (UserSeed) GetDescription() string {
	return "Create default users"
}

func (UserSeed) Seed(db helpers.IMongoClient) {
	fmt.Println("seeding")
}

func (UserSeed) Remove() {

}
