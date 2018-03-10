// +build integration

package crazy_nl_backend_test

import (
	"testing"
	"crazy_nl_backend/migrations"
	"crazy_nl_backend/config"
	"crazy_nl_backend/helpers"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestApplicationInit(t *testing.T) {
	session, err := helpers.GetMongo(config.GetMongoConfig())
	assert.Nil(t, err)
	assert.NotNil(t, session)
	db := session.DB(config.GetTestingDbName())
	assert.NotNil(t, db)
	defer session.Close()
	db.DropDatabase()
	migrations.ApplyAll(config.GetTestingDbName())
	count, err := db.C(migrations.SeedingCollectionName).Count()
	assert.Nil(t, err)
	assert.NotNil(t, count)
	fmt.Println(count)
	assert.True(t, count > 0)
}
