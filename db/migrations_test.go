// +build integration

package db_test

import (
	"testing"

	"github.com/fossapps/starter/config"
	"github.com/fossapps/starter/db"
	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"
)

func TestMigrationManagerLayer_IsApplied(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()
	expect.NotNil(database)
	migrationManagerLayer := db.GetMigrationManager(database)
	result, err := migrationManagerLayer.IsApplied("migration_key")
	expect.False(result)
	expect.Nil(err)
	err = migrationManagerLayer.MarkApplied("migration_key", "description")
	expect.Nil(err)
	result, err = migrationManagerLayer.IsApplied("migration_key")
	expect.True(result)
	expect.Nil(err)
}
