package migration

import (
	"github.com/fossapps/starter/db"

	"github.com/globalsign/mgo"
)

// Migration interface to implement for a migration
type Migration interface {
	GetKey() string
	GetDescription() string
	Apply(dbLayer db.DB)
	Remove(db *mgo.Database)
}
