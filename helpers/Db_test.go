// +build integration

package helpers_test

import (
	"testing"
	"crazy_nl_backend/helpers"
	"crazy_nl_backend/config"
	"github.com/stretchr/testify/assert"
	"github.com/globalsign/mgo/bson"
)
type Name struct {
	Name string `json:"name"`
}

func TestGetMongo(t *testing.T) {
	mongo, _ := helpers.GetMongo(config.GetMongoConfig())
	defer mongo.Close()
	db := mongo.DB(config.GetTestingDbName())
	db.DropDatabase()
	collection := db.C("test")
	collection.Insert(Name{Name:"Test_Name"})
	d := collection.Find(bson.M{
		"name": "Test_Name",
	})
	result := new(Name)
	d.One(&result)
	assert.NotNil(t, result.Name)
	assert.Equal(t, "Test_Name", result.Name)
}

func TestInsertAfterCopy(t *testing.T) {
	mongo, _ := helpers.GetMongo(config.GetMongoConfig())
	defer mongo.Close()
	db := mongo.DB(config.GetTestingDbName())
	db.DropDatabase()
	collection := db.C("test")
	collection.Insert(Name{Name:"Test_Name"})
	d := collection.Find(bson.M{
		"name": "Test_Name",
	})
	result := new(Name)
	d.One(&result)
	assert.NotNil(t, result.Name)
	assert.Equal(t, "Test_Name", result.Name)
}

