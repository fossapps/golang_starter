package helpers

type ICollection interface {
	Insert(...interface{}) error
}

type DbCollection struct {
	Collection ICollection
}

func (c DbCollection) Insert(docs ...interface{}) error {
	return c.Collection.Insert(docs...)
}
