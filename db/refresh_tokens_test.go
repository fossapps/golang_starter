//go:generate mockgen -destination=../mocks/mock_refresh_tokens.go -package=mocks crazy_nl_backend/db IRefreshTokenManager
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
