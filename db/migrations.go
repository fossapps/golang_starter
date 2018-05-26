package db

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// MigrationManager deals with any info related to migration saved to disk
type MigrationManager interface {
	MarkApplied(key string, description string)
	ShouldRun(key string) bool
}

// MigrationInfo information about migration
type MigrationInfo struct {
	Key         string    `json:"key"`
	Description string    `json:"description"`
	AppliedAt   time.Time `json:"applied_at"`
}

type migrationManagerLayer struct {
	db *mgo.Database
}

// seedingCollectionName name of collection which holds info about migration
const seedingCollectionName = "migration"

// GetMigrationManager returns implementation of MigrationManager
func GetMigrationManager(db *mgo.Database) MigrationManager {
	return migrationManagerLayer{
		db: db,
	}
}

// MarkApplied marks a migration key as applied
func (dbLayer migrationManagerLayer) MarkApplied(key string, description string) {
	info := MigrationInfo{
		Key:         key,
		Description: description,
		AppliedAt:   time.Now(),
	}
	dbLayer.db.C(seedingCollectionName).Insert(info)
}

// ShouldRun determine weather or not a migration should run, if it's already applied, then it returns false
func (dbLayer migrationManagerLayer) ShouldRun(key string) bool {
	collection := dbLayer.db.C(seedingCollectionName)
	result := new(MigrationInfo)
	collection.Find(bson.M{
		"key": key,
	}).One(&result)
	return result.Key != key
}
