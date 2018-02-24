// +build integration

package helpers_test

import (
	"testing"
	"crazy_nl_backend/helpers"
	"crazy_nl_backend/config"
	"github.com/stretchr/testify/assert"
)
type Name struct {
	Name string `json:"name"`
}

func TestGetMongo(t *testing.T) {
	mongo, _ := helpers.GetMongo(config.GetMongoConfig())
	defer mongo.Close()
	db := mongo.DB(config.GetTestingDbName())
	defer DropDb(db)
	collection := db.C("test")
	collection.Insert(Name{Name:"Test_Name"})
	d := collection.Find(Name{Name: "Test_Name"})
	result := new(Name)
	d.One(&result)
	assert.NotNil(t, result.Name)
	assert.Equal(t, result.Name, "Test_Name")
}

func TestInsertAfterCopy(t *testing.T) {
	mongo, _ := helpers.GetMongo(config.GetMongoConfig())
	defer mongo.Close()
	db := mongo.DB(config.GetTestingDbName())
	defer DropDb(db)
	collection := db.C("test")
	collection.Insert(Name{Name:"Test_Name"})
	d := collection.Find(Name{Name: "Test_Name"})
	result := new(Name)
	d.One(&result)
	assert.NotNil(t, result.Name)
	assert.Equal(t, result.Name, "Test_Name")
}

