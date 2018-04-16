//go:generate mockgen -destination=../mocks/mock_users.go -package=mocks crazy_nl_backend/db IUserManager
// +build integration

package db_test

import (
	"testing"

	"crazy_nl_backend/config"
	"crazy_nl_backend/db"

	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"
)

func TestUserLayer_FindByEmailReturnsNilIfNotExists(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()

	userManager := db.GetUserManager(database)
	expect.Nil(userManager.FindByEmail("randomEmail@example.com"))
}

func TestUserLayer_FindByIdReturnsNilIfNotExists(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()

	userManager := db.GetUserManager(database)
	// random id
	expect.Nil(userManager.FindById("5adbf94839b5b200068fa33e"))
}

func TestUserLayer_CreateNewUser(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()

	userManager := db.GetUserManager(database)
	email := "testEmail@example.com"
	password := "test_password"
	userManager.Create(db.User{
		Permissions: []string{"test"},
		Password:    password,
		Email:       email,
	})
	user := userManager.FindByEmail(email)
	expect.Equal(email, user.Email)
	expect.NotEmpty(password, user.Password)
}

func TestUserLayer_FindById(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()

	userManager := db.GetUserManager(database)
	email := "testEmail1@example.com"
	password := "test_password"
	userManager.Create(db.User{
		Permissions: []string{"test"},
		Password:    password,
		Email:       email,
	})
	id := userManager.FindByEmail(email).ID
	user := userManager.FindById(id.Hex())
	expect.Equal(email, user.Email)
}

func TestUserLayer_CreateReturnsErrorIfEmailAlreadyExists(t *testing.T) {
	t.Skip()
}

func TestUserLayer_Exists(t *testing.T) {
	t.Skip()
}
