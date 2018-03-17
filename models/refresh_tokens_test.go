package models_test

import (
	"testing"
	"github.com/globalsign/mgo"
	"crazy_nl_backend/config"
	"github.com/stretchr/testify/assert"
	"crazy_nl_backend/models"
)

func TestRefreshToken_FindOne(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetTestingDbConnection())
	expect.NotNil(session)
	expect.Nil(err)
	// first simply add one
	session.DB(config.GetMongoConfig().DbName).C("refresh_tokens").Insert(
		models.RefreshToken{
			User: "some_id",
			Token: "my_token",
		},
	)
	// then try to get
	token := models.RefreshToken{}.FindOne("my_token", session.DB(config.GetMongoConfig().DbName))
	expect.Equal("my_token", token.Token)
	expect.Equal("some_id", token.User)
}

func TestRefreshToken_FindOneReturnsNilIfNotFound(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetTestingDbConnection())
	expect.NotNil(session)
	expect.Nil(err)
	token := models.RefreshToken{}.FindOne("some_random_token_which_should_not_exist", session.DB(config.GetMongoConfig().DbName))
	expect.Nil(token)
}
