// +build integration

package db_test

import (
	"testing"

	"golang_starter/config"
	"golang_starter/db"

	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"
)

func TestMigrationManagerLayer_ShouldRun(t *testing.T) {
	expect := assert.New(t)
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	expect.Nil(err)
	database := session.DB(config.GetTestingDbName())
	defer database.DropDatabase()
	expect.NotNil(database)
	migrationManagerLayer := db.GetMigrationManager(database)
	expect.True(migrationManagerLayer.ShouldRun("migration_key"))
	migrationManagerLayer.MarkApplied("migration_key", "description")
	expect.False(migrationManagerLayer.ShouldRun("migration_key"))
}
