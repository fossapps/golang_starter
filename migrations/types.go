package migrations

import (
	"crazy_nl_backend/db"

	"github.com/globalsign/mgo"
)

type IMigration interface {
	GetKey() string
	GetDescription() string
	Apply(dbLayer db.Db)
	Remove(db *mgo.Database)
}
