package migrations

import (
	"time"
	"github.com/globalsign/mgo"
)

type IMigration interface {
	GetKey() string
	GetDescription() string
	Apply(db *mgo.Database)
	Remove(db *mgo.Database)
}
type MigrationInfo struct {
	Key string `json:"key"`
	Description string `json:"description"`
	AppliedAt time.Time `json:"applied_at"`
}
