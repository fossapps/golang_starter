package migrations

import (
	"fmt"
	"time"
	"github.com/globalsign/mgo/bson"
	"github.com/globalsign/mgo"
)

const SeedingCollectionName = "migrations"

func ApplyAll(dbName string, session *mgo.Session) {
	Apply(UserSeed{}, session.DB(dbName))
	Apply(PermissionSeeds{}, session.DB(dbName))
}

func Apply(seeder IMigration, db *mgo.Database) {
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

func shouldRun(seeder IMigration, db *mgo.Database) bool {
	key := seeder.GetKey()
	collection := db.C(SeedingCollectionName)
	result := new(MigrationInfo)
	collection.Find(bson.M{
		"key": key,
	}).One(&result)
	return result.Key != key
}

func markApplied(seeder IMigration, db *mgo.Database) {
	info := MigrationInfo{
		Key:seeder.GetKey(),
		Description:seeder.GetDescription(),
		AppliedAt:time.Now(),
	}
	db.C(SeedingCollectionName).Insert(info)
}
