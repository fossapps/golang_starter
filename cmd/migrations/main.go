package main

import (
	"starter/config"
	"starter/db"
	"starter/migrations"

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
