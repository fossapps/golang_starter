package main

import (
	"github.com/fossapps/starter/config"
	"github.com/fossapps/starter/db"
	"github.com/fossapps/starter/migration"
	"github.com/globalsign/mgo"
)

func main() {
	mongo, err := mgo.Dial(config.GetMongoConfig().Connection)
	if err != nil {
		panic(err)
	}
	migrationManagerLayer := db.GetDbImplementation(mongo)
	migration.ApplyAll(migrationManagerLayer)
}
