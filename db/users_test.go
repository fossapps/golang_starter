// +build integration

package db_test

import (
	"testing"

	"github.com/fossapps/starter/config"
	"github.com/fossapps/starter/db"

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
	result, err := userManager.FindByEmail("randomEmail@example.com")
	expect.Nil(result)
	expect.Nil(err)
}

func TestUserLayer_FindByIdReturnsNilIfNotExists(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()

	userManager := db.GetUserManager(database)
	// random id
	result, err := userManager.FindByID("5adbf94839b5b200068fa33e")
	expect.Nil(result)
	expect.Nil(err)
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
	user, err := userManager.FindByEmail(email)
	expect.Equal(email, user.Email)
	expect.Nil(err)
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
	user, err := userManager.FindByEmail(email)
	id := user.ID
	user, err = userManager.FindByID(id.Hex())
	expect.Equal(email, user.Email)
	expect.Nil(err)
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
	user, err := userManager.FindByEmail("mail_test@example.com")
	expect.Nil(err)
	expect.Equal("mail_test@example.com", user.Email)
	err = userManager.Update(user.ID.Hex(), db.User{Email: "mail_updated@example.com"})
	expect.Nil(err)
	updatedUser, err := userManager.FindByID(user.ID.Hex())
	expect.Nil(err)
	expect.Equal("mail_updated@example.com", updatedUser.Email)
	// expect.Equal("sudo", updatedUser.Permissions[0])
}

func TestUserLayer_CreateReturnsErrorIfEmailAlreadyExists(t *testing.T) {
	t.Skip()
}

func TestUserLayer_Exists(t *testing.T) {
	t.Skip()
}
