// +build integration

package db_test

import (
	"testing"

	"github.com/fossapps/starter/config"
	"github.com/fossapps/starter/db"

	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"
)

func TestPermissionLayer_Create(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	database.DropDatabase()
	permissionManager := db.GetPermissionManager(database)
	key := "test_permission"
	description := "random description"
	err = permissionManager.Create(key, description)
	expect.Nil(err)
	permissions, err := permissionManager.List()
	expect.Nil(err)
	expect.Equal(key, permissions[0].Key)
	expect.Equal(description, permissions[0].Description)
}

func TestPermissionLayer_Create_ReturnsErrorIfAlreadyExists(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()
	permissionManager := db.GetPermissionManager(database)
	key := "permission1"
	description := "permission description"
	permissionManager.Create(key, description)
	err = permissionManager.Create(key, description)
	expect.NotNil(err)
}

func TestPermissionManager_Exists(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()
	permissionManager := db.GetPermissionManager(database)
	permissionManager.Create("test_permission", "description")
	res, err := permissionManager.Exists("test_permission")
	expect.True(res)
	expect.Nil(err)
}
func TestPermissionManager_List(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()
	permissionManager := db.GetPermissionManager(database)
	permissionManager.Create("test_permission", "description")
	list, err := permissionManager.List()
	expect.Nil(err)
	expect.True(len(list) > 0)
}
