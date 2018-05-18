// +build integration

package db_test

import (
	"testing"

	"golang_starter/config"
	"golang_starter/db"

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

func TestUserLayer_List(t *testing.T) {
	mockUsers := []db.User{
		{Email: "mail@example.com", Permissions: []string{"sudo"}, Password: "test_password"},
		{Email: "mail2@example.com", Permissions: []string{"sudo"}, Password: "test_password2"},
	}
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()
	userManager := db.GetUserManager(database)
	userManager.Create(mockUsers[0])
	userManager.Create(mockUsers[1])
	users, err := userManager.List()
	expect.Nil(err)
	expect.NotNil(users)
	expect.Equal(mockUsers[0].Email, users[0].Email)
	expect.Equal(mockUsers[0].Permissions, users[0].Permissions)
	expect.Equal(mockUsers[1].Email, users[1].Email)
	expect.Equal(mockUsers[1].Permissions, users[1].Permissions)
}

func TestUserLayer_Edit(t *testing.T) {
	mockUsers := []db.User{
		{Email: "mail_test@example.com", Permissions: []string{"sudo"}, Password: "test_password"},
		{Email: "mail2@example.com", Permissions: []string{"sudo"}, Password: "test_password2"},
	}
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()
	userManager := db.GetUserManager(database)
	expect.Nil(userManager.Create(mockUsers[0]))
	expect.Nil(userManager.Create(mockUsers[1]))
	user := userManager.FindByEmail("mail_test@example.com")
	expect.Equal("mail_test@example.com", user.Email)
	err = userManager.Edit(user.ID.Hex(), db.User{Email: "mail_updated@example.com"})
	expect.Nil(err)
	updatedUser := userManager.FindById(user.ID.Hex())
	expect.Equal("mail_updated@example.com", updatedUser.Email)
	// expect.Equal("sudo", updatedUser.Permissions[0])
}

func TestUserLayer_CreateReturnsErrorIfEmailAlreadyExists(t *testing.T) {
	t.Skip()
}

func TestUserLayer_Exists(t *testing.T) {
	t.Skip()
}
