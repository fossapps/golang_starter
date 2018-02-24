package seeds

import (
	"fmt"
	"crazy_nl_backend/helpers"
	"crazy_nl_backend/config"
)

const SEEDING_COLLECTION_NAME = "seeds"

func SeedDb() {
	Seed(UserSeed{})
}

func Seed(seeder ISeeder) {
	key := seeder.GetKey()
	description := seeder.GetDescription()
	fmt.Printf("seeding file: %s", key)
	fmt.Println(description)
	db, err := helpers.GetMongo(config.GetMongoConfig())
	if err != nil {
		panic(err)
	}
	seeder.Seed(db)
}
