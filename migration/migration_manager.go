package migration

import (
	"fmt"

	"github.com/fossapps/starter/db"
)

// ApplyAll applies all migration to a specific database
func ApplyAll(dbLayer db.DB) {
	Apply(UserSeed{}, dbLayer)
	Apply(PermissionSeeds{}, dbLayer)
}

// Apply individual migration to a db
func Apply(migration Migration, dbLayer db.DB) {
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
