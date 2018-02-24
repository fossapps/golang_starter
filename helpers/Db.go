package helpers

import "github.com/globalsign/mgo"

type IDatabase interface {
	C(string) ICollection
	DropDatabase() error
}

type IDatabaseInstance interface {
	C(string) *mgo.Collection
	DropDatabase() error
}

type Db struct {
	DB IDatabaseInstance
}

func (d Db) C(key string) ICollection {
	return DbCollection{
		Collection: d.DB.C(key),
	}
}

func (d Db) DropDatabase() error {
	return d.DB.DropDatabase()
}
