package helpers

type IDatabase interface {
	C(string) ICollection
}

type IDatabaseInstance interface {}

type Db struct {
	DB IDatabaseInstance
}

func (d Db) C(key string) ICollection {
	return DbCollection{
		Collection: nil,
	}
}
