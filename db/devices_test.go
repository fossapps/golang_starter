//go:generate mockgen -destination=../mocks/mock_devices.go -package=mocks crazy_nl_backend/db IDeviceManager
// +build integration

package db_test

import (
	"testing"
	"github.com/globalsign/mgo"
	"crazy_nl_backend/config"
	"crazy_nl_backend/db"
	"github.com/stretchr/testify/assert"
)

func TestDeviceManager_Register(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()
	deviceManager := db.GetDeviceManager(database)
	token := "random_device_token"
	deviceManager.Register(token)
	expect.Equal(token, deviceManager.FindByToken(token).Token)
}

func TestDeviceManager_RegisterReturnsErrorOnDuplicate(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()
	deviceManager := db.GetDeviceManager(database)
	token := "random_token"
	expect.Nil(deviceManager.Register(token))
	expect.NotNil(deviceManager.Register(token))
}

func TestDeviceManager_Exists(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()
	deviceManager := db.GetDeviceManager(database)
	token := "token"
	expect.False(deviceManager.Exists(token))
	deviceManager.Register(token)
	expect.True(deviceManager.Exists(token))
}

func TestDeviceManager_FindByTokenReturnsNilIfNotFound(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()
	deviceManager := db.GetDeviceManager(database)
	token := "find_token"
	expect.Nil(deviceManager.FindByToken(token))
	deviceManager.Register(token)
	expect.NotNil(deviceManager.FindByToken(token))
}
