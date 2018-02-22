package seeds

import (
	"fmt"
	"crazy_nl_backend/helpers"
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
	db, err := helpers.GetMongo()
	if err != nil {
		panic(err)
	}
	seeder.Seed(db)
}
