// +build integration

package db_test

import (
	"testing"

	"github.com/fossapps/starter/config"
	"github.com/fossapps/starter/db"
	"github.com/globalsign/mgo"
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
	err = deviceManager.Register(token)
	expect.Nil(err)
	result, err := deviceManager.FindByToken(token)
	expect.Equal(token, result.Token)
	expect.Nil(err)
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
	result, err := deviceManager.Exists(token)
	expect.False(result)
	expect.Nil(err)
	deviceManager.Register(token)
	result, err = deviceManager.Exists(token)
	expect.True(result)
	expect.Nil(err)
}

func TestDeviceManager_FindByTokenReturnsNilIfNotFound(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()
	deviceManager := db.GetDeviceManager(database)
	token := "find_token"
	res, err := deviceManager.FindByToken(token)
	expect.Nil(res)
	expect.Nil(err)
	deviceManager.Register(token)
	res, err = deviceManager.FindByToken(token)
	expect.NotNil(res)
	expect.Nil(err)
}
