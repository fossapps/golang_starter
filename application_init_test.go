package crazy_nl_backend_test

import (
	"testing"
	"github.com/globalsign/mgo"
	"crazy_nl_backend/config"
	"github.com/stretchr/testify/assert"
	"crazy_nl_backend/migrations"
)

func TestApplicationInit(t *testing.T) {
	session, err := mgo.Dial(config.GetTestingDbConnection())
	assert.Nil(t, err)
	assert.NotNil(t, session)
	db := session.DB(config.GetMongoConfig().DbName)
	assert.NotNil(t, db)
	defer session.Close()
	count, err := db.C(migrations.SeedingCollectionName).Count()
	assert.Nil(t, err)
	assert.NotNil(t, count)
	assert.True(t, count > 0)
}
