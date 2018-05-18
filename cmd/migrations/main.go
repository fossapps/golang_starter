package main

import (
	"golang_starter/config"
	"golang_starter/db"
	"golang_starter/migrations"

	"github.com/globalsign/mgo"
)

func main() {
	mongo, err := mgo.Dial(config.GetMongoConfig().Connection)
	if err != nil {
		panic(err)
	}
	migrationManagerLayer := db.GetDbImplementation(mongo)
	migrations.ApplyAll(migrationManagerLayer)
}
