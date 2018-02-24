// +build integration

package seeds_test

import (
	"github.com/stretchr/testify/mock"
	"crazy_nl_backend/helpers"
	"testing"
	"crazy_nl_backend/seeds"
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

func(m *MockSeeder) Seed(db helpers.IDatabase) {
	m.Called(db)
}

func (m *MockSeeder) Remove(db helpers.IDatabase) {}
// endregion

func TestSeedCallsFirstTime(t *testing.T) {
	mockedSeeder := new(MockSeeder)
	mockedSeeder.On("GetKey").Return("test_seed")
	mockedSeeder.On("GetDescription").Return("test_seed description")
	mockedSeeder.On("Seed", mock.Anything).Return(nil)

	session, err := helpers.GetMongo(config.GetMongoConfig())
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	assert.Nil(t, err)
	assert.NotNil(t, session)
	db := session.DB(config.GetTestingDbName())
	defer db.DropDatabase()
	seeds.Seed(mockedSeeder, db)
	c := db.C(seeds.SeedingCollectionName)
	info := new(seeds.SeedInfo)
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
	defer db.DropDatabase()
	c := db.C(seeds.SeedingCollectionName)
	c.Insert(seeds.SeedInfo{
		Key: "test_seed",
	})
	seeds.Seed(mockedSeeder, db)
}
