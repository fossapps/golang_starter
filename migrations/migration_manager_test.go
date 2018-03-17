//go:generate mockgen -destination=../mocks/mock_migration.go -package=mocks crazy_nl_backend/migrations IMigration

package migrations_test

import (
	"testing"
	"crazy_nl_backend/mocks"
	"github.com/golang/mock/gomock"
	"github.com/globalsign/mgo"
	"crazy_nl_backend/config"
	"github.com/stretchr/testify/assert"
	"crazy_nl_backend/migrations"
	"github.com/globalsign/mgo/bson"
)

func TestSeedCallsFirstTime(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockMigration := mocks.NewMockIMigration(controller)
	mockMigration.EXPECT().GetKey().MinTimes(1).Return("key")
	mockMigration.EXPECT().GetDescription().MinTimes(1).Return("description")
	mockMigration.EXPECT().Apply(gomock.Any()).Times(1)

	session, err := mgo.Dial(config.GetTestingDbConnection())
	assert.Nil(t, err)
	assert.NotNil(t, session)
	db := session.DB(config.GetMongoConfig().DbName)
	db.C(migrations.SeedingCollectionName).DropCollection()
	migrations.Apply(mockMigration, db)
	c := db.C(migrations.SeedingCollectionName)
	info := new(migrations.MigrationInfo)
	query := c.Find(bson.M{
		"key": "key",
	})
	query.One(&info)
	assert.NotNil(t, info.Key)
	assert.Equal(t, "key", info.Key)
}

func TestSeedDoesNotExecuteDuplicates(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockMigration := mocks.NewMockIMigration(controller)
	mockMigration.EXPECT().GetKey().Times(1).Return("key")
	mockMigration.EXPECT().GetDescription().Times(0)
	mockMigration.EXPECT().Apply(gomock.Any()).Times(0)
	session, err := mgo.Dial(config.GetTestingDbConnection())
	assert.Nil(t, err)
	assert.NotNil(t, session)
	db := session.DB(config.GetMongoConfig().DbName)
	db.C(migrations.SeedingCollectionName).DropCollection()
	c := db.C(migrations.SeedingCollectionName)
	c.Insert(migrations.MigrationInfo{
		Key: "key",
	})
	migrations.Apply(mockMigration, db)
}
