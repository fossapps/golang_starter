// +build integration

package db_test

import (
	"testing"

	"crazy_nl_backend/config"
	"crazy_nl_backend/db"

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
	permissionManager.Create(key, description)
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
