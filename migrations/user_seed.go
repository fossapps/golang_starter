package migrations

import (
	"fmt"
	"github.com/fossapps/starter/db"

	"github.com/globalsign/mgo"
)

// UserSeed seed initial list of users
type UserSeed struct{}

// GetKey returns key for user seeds
func (UserSeed) GetKey() string {
	return "USER_SEED"
}

// GetDescription returns description for user seeds
func (UserSeed) GetDescription() string {
	return "Create default users"
}

// Apply adds users to database
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

// Remove remove un does the migration
func (UserSeed) Remove(db *mgo.Database) {

}
