// +build integration

package crazy_nl_backend_test

import (
	"testing"
	"crazy_nl_backend/seeds"
	"crazy_nl_backend/config"
	"crazy_nl_backend/helpers"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestApplicationInit(t *testing.T) {
	seeds.SeedDb(config.GetTestingDbName())
	session, err := helpers.GetMongo(config.GetMongoConfig())
	assert.Nil(t, err)
	assert.NotNil(t, session)
	db := session.DB(config.GetTestingDbName())
	assert.NotNil(t, db)
	defer session.Close()
	defer db.DropDatabase()
	count, err := db.C(seeds.SeedingCollectionName).Count()
	assert.Nil(t, err)
	assert.NotNil(t, count)
	fmt.Println(count)
	assert.True(t, count > 0)
}
