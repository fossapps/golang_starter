package models_test

import (
	"crazy_nl_backend/config"
	"crazy_nl_backend/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_FindUserByEmail(t *testing.T) {
	// first insert
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetTestingDbConnection())
	expect.NotNil(session)
	expect.Nil(err)
	session.DB(config.GetMongoConfig().DbName).C("users").Insert(models.User{
		Email:       "test@example.com",
		Permissions: []string{"sudo"},
	})
	// try to get
	user := models.User{}.FindUserByEmail("test@example.com", session.DB(config.GetMongoConfig().DbName))
	expect.Equal("test@example.com", user.Email)
}
func TestUser_FindUserById(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetTestingDbConnection())
	expect.NotNil(session)
	expect.Nil(err)
	id := bson.NewObjectId()
	session.DB(config.GetMongoConfig().DbName).C("users").Insert(models.User{
		Email:       "test1@example.com",
		ID:          id,
		Permissions: []string{"sudo"},
	})
	// try to get
	user := models.User{}.FindUserById(id.Hex(), session.DB(config.GetMongoConfig().DbName))
	expect.Equal("test1@example.com", user.Email)
	expect.Equal(id.Hex(), user.ID.Hex())
}

func TestUser_FindUserByIdReturnsNilIfIdNotFound(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetTestingDbConnection())
	expect.NotNil(session)
	expect.Nil(err)
	user := models.User{}.FindUserById("aaaaaaaaaaaaaaaaaaaaaaaa", session.DB(config.GetMongoConfig().DbName))
	expect.Nil(user)
}
