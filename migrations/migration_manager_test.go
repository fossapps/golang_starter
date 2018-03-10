// +build integration

package migrations_test

import (
	"github.com/stretchr/testify/mock"
	"crazy_nl_backend/helpers"
	"testing"
	"crazy_nl_backend/migrations"
	"crazy_nl_backend/config"
	"github.com/stretchr/testify/assert"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// region MockSeed
type MockSeeder struct {
	mock.Mock
}

func(m *MockSeeder) GetKey() string {
	args := m.Called()
	return args.String(0)
}

func(m *MockSeeder) GetDescription() string {
	args := m.Called()
	return args.String(0)
}

func(m *MockSeeder) Apply(db helpers.IDatabase) {
	m.Called(db)
}

func (m *MockSeeder) Remove(db helpers.IDatabase) {}
// endregion

func TestSeedCallsFirstTime(t *testing.T) {
	mockedSeeder := new(MockSeeder)
	mockedSeeder.On("GetKey").Return("test_seed")
	mockedSeeder.On("GetDescription").Return("test_seed description")
	mockedSeeder.On("Apply", mock.Anything).Return(nil)

	session, err := helpers.GetMongo(config.GetMongoConfig())
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	assert.Nil(t, err)
	assert.NotNil(t, session)
	db := session.DB(config.GetTestingDbName())
	db.DropDatabase()
	migrations.Apply(mockedSeeder, db)
	c := db.C(migrations.SeedingCollectionName)
	info := new(migrations.MigrationInfo)
	query := c.Find(bson.M{
		"key": "test_seed",
	})
	query.One(&info)
	assert.NotNil(t, info.Key)
	assert.Equal(t, "test_seed", info.Key)
}

func TestSeedDoesNotExecuteDuplicates(t *testing.T) {
	mockedSeeder := new(MockSeeder)
	mockedSeeder.On("GetKey").Return("test_seed")

	session, err := helpers.GetMongo(config.GetMongoConfig())
	defer session.Close()

	assert.Nil(t, err)
	assert.NotNil(t, session)
	db := session.DB(config.GetTestingDbName())
	db.DropDatabase()
	c := db.C(migrations.SeedingCollectionName)
	c.Insert(migrations.MigrationInfo{
		Key: "test_seed",
	})
	migrations.Apply(mockedSeeder, db)
}
