package migrations

import (
	"fmt"
	"crazy_nl_backend/helpers"
	"crazy_nl_backend/config"
	"time"
	"github.com/globalsign/mgo/bson"
)

const SeedingCollectionName = "migrations"

func ApplyAll(dbName string) {
	session, err := helpers.GetMongo(config.GetMongoConfig())
	defer session.Close()
	if err != nil {
		panic(err)
	}
	Apply(UserSeed{}, session.DB(dbName))
	Apply(PermissionSeeds{}, session.DB(dbName))
}

func Apply(seeder IMigration, db helpers.IDatabase) {
	if !shouldRun(seeder, db) {
		return
	}
	key := seeder.GetKey()
	description := seeder.GetDescription()
	fmt.Printf("applying migration file: %s\n", key)
	fmt.Println(description)
	seeder.Apply(db)
	markApplied(seeder, db)
}

func shouldRun(seeder IMigration, db helpers.IDatabase) bool {
	key := seeder.GetKey()
	collection := db.C(SeedingCollectionName)
	result := new(MigrationInfo)
	collection.Find(bson.M{
		"key": key,
	}).One(&result)
	return result.Key != key
}
func markApplied(seeder IMigration, db helpers.IDatabase) {
	info := MigrationInfo{
		Key:seeder.GetKey(),
		Description:seeder.GetDescription(),
		AppliedAt:time.Now(),
	}
	db.C(SeedingCollectionName).Insert(info)
}
