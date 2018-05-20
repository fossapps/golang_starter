package migrations

import (
	"starter/db"

	"github.com/globalsign/mgo"
)

// IMigration interface to implement for a migration
type IMigration interface {
	GetKey() string
	GetDescription() string
	Apply(dbLayer db.Db)
	Remove(db *mgo.Database)
}
