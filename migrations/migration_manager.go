package migrations

import (
	"fmt"

	"golang_starter/db"
)

func ApplyAll(dbLayer db.Db) {
	Apply(UserSeed{}, dbLayer)
	Apply(PermissionSeeds{}, dbLayer)
}

func Apply(seeder IMigration, dbLayer db.Db) {
	if !dbLayer.Migrations().ShouldRun(seeder.GetKey()) {
		return
	}
	key := seeder.GetKey()
	description := seeder.GetDescription()
	fmt.Printf("applying migration file: %s\n", key)
	fmt.Println(description)
	seeder.Apply(dbLayer)
	dbLayer.Migrations().MarkApplied(seeder.GetKey(), seeder.GetDescription())
}
