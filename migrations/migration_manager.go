package migrations

import (
	"fmt"

	"github.com/fossapps/starter/db"
)

// ApplyAll applies all migration to a specific database
func ApplyAll(dbLayer db.Db) {
	Apply(UserSeed{}, dbLayer)
	Apply(PermissionSeeds{}, dbLayer)
}

// Apply individual migration to a db
func Apply(migration IMigration, dbLayer db.Db) {
	if !dbLayer.Migrations().ShouldRun(migration.GetKey()) {
		return
	}
	key := migration.GetKey()
	description := migration.GetDescription()
	fmt.Printf("applying migration file: %s\n", key)
	fmt.Println(description)
	migration.Apply(dbLayer)
	dbLayer.Migrations().MarkApplied(migration.GetKey(), migration.GetDescription())
}
