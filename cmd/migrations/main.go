package main

import (
	"crazy_nl_backend/config"
	"crazy_nl_backend/db"
	"crazy_nl_backend/migrations"

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
