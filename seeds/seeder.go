package seeds

import (
	"fmt"
	"crazy_nl_backend/helpers"
	"crazy_nl_backend/config"
	"time"
)

const SeedingCollectionName = "seeds"

func SeedDb(dbName string) {
	session, err := helpers.GetMongo(config.GetMongoConfig())
	defer session.Close()
	if err != nil {
		panic(err)
	}
	Seed(UserSeed{}, session.DB(dbName))
}

func Seed(seeder ISeeder, db helpers.IDatabase) {
	if !shouldRun(seeder, db) {
		return
	}
	key := seeder.GetKey()
	description := seeder.GetDescription()
	fmt.Printf("seeding file: %s", key)
	fmt.Println(description)
	seeder.Seed(db)
	markSeeded(seeder, db)
}

func shouldRun(seeder ISeeder, db helpers.IDatabase) bool {
	key := seeder.GetKey()
	collection := db.C(SeedingCollectionName)
	result := new(SeedInfo)
	collection.Find(SeedInfo{
		Key: key,
	}).One(&result)
	return result.Key != key
}
func markSeeded(seeder ISeeder, db helpers.IDatabase) {
	info := SeedInfo{
		Key:seeder.GetKey(),
		Description:seeder.GetDescription(),
		AppliedAt:time.Now(),
	}
	db.C(SeedingCollectionName).Insert(info)
}
