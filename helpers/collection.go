package helpers

import "github.com/globalsign/mgo"

type ICollection interface {
	Insert(...interface{}) error
	Find(interface{}) *mgo.Query
}

type DbCollection struct {
	Collection ICollection
}

func (c DbCollection) Insert(docs ...interface{}) error {
	return c.Collection.Insert(docs...)
}

func (c DbCollection) Find(query interface{}) *mgo.Query {
	return c.Collection.Find(query)
}
