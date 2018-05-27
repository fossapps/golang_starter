package db

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// MigrationManager deals with any info related to migration saved to disk
type MigrationManager interface {
	MarkApplied(key string, description string) error
	IsApplied(key string) (bool, error)
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
func (dbLayer migrationManagerLayer) MarkApplied(key string, description string) error {
	info := MigrationInfo{
		Key:         key,
		Description: description,
		AppliedAt:   time.Now(),
	}
	return dbLayer.db.C(seedingCollectionName).Insert(info)
}

// IsApplied determine weather or not a migration should run, if it's already applied, then it returns false
func (dbLayer migrationManagerLayer) IsApplied(key string) (bool, error) {
	collection := dbLayer.db.C(seedingCollectionName)
	result := new(MigrationInfo)
	err := collection.Find(bson.M{
		"key": key,
	}).One(&result)
	if err == mgo.ErrNotFound {
		return false, nil
	}
	return result.Key == key, err
}
