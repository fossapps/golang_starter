// +build integration

package db_test

import (
	"testing"

	"crazy_nl_backend/config"
	"crazy_nl_backend/db"

	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"
)

func TestRefreshTokenLayer_FindOneReturnsNilIfRefreshTokenDoesNotExist(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()

	refreshTokenManager := db.GetRefreshTokenManager(database)
	expect.Nil(refreshTokenManager.FindOne("tokenWhichNeverExist"))
}

func TestRefreshTokenLayer_AddWorksAsExpected(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()

	refreshTokenManager := db.GetRefreshTokenManager(database)
	token := "myTestToken"
	user := "testUser1"
	refreshTokenManager.Add(token, user)
	refToken := refreshTokenManager.FindOne(token)
	expect.Equal(token, refToken.Token)
	expect.Equal(user, refToken.User)
}

func TestRefreshTokenLayer_List(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()

	refreshTokenManager := db.GetRefreshTokenManager(database)
	user1 := "testUser1"
	user2 := "testUser2"
	refreshTokenManager.Add("test_token1", user1)
	refreshTokenManager.Add("test_token2", user1)
	refreshTokenManager.Add("test_token3", user2)
	refreshTokenManager.Add("test_token4", user2)
	tokens, err := refreshTokenManager.List(user1)
	expect.Nil(err)
	expect.Equal(tokens[0].Token, "test_token1")
	expect.Equal(tokens[1].Token, "test_token2")
	tokens, err = refreshTokenManager.List(user2)
	expect.Nil(err)
	expect.Equal(tokens[0].Token, "test_token3")
	expect.Equal(tokens[1].Token, "test_token4")
}

func TestRefreshTokenLayer_Delete(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()
	refreshTokenManager := db.GetRefreshTokenManager(database)
	// prepare data
	user1 := "testUser1"
	user2 := "testUser2"
	refreshTokenManager.Add("test_token1", user1)
	refreshTokenManager.Add("test_token2", user1)
	refreshTokenManager.Add("test_token3", user2)
	refreshTokenManager.Add("test_token4", user2)
	// perform delete
	expect.NotNil(refreshTokenManager.FindOne("test_token1"))
	expect.NotNil(refreshTokenManager.FindOne("test_token2"))
	expect.NotNil(refreshTokenManager.FindOne("test_token3"))
	expect.NotNil(refreshTokenManager.FindOne("test_token4"))
	// everything exists till now
	expect.Nil(refreshTokenManager.Delete("test_token1"))
	expect.Nil(refreshTokenManager.FindOne("test_token1"))
	expect.Nil(refreshTokenManager.Delete("test_token2"))
	expect.Nil(refreshTokenManager.FindOne("test_token2"))
	expect.Nil(refreshTokenManager.Delete("test_token3"))
	expect.Nil(refreshTokenManager.FindOne("test_token3"))
	expect.Nil(refreshTokenManager.Delete("test_token4"))
	expect.Nil(refreshTokenManager.FindOne("test_token4"))
	// maybe it's a good idea to add table driven tests here.
}
