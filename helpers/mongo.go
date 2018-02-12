package helpers

import (
	"github.com/globalsign/mgo"
)

type Mongo struct {
	session *mgo.Session
}

type IMongoClient interface {
	Close()
	Clone() IMongoClient
	Copy() IMongoClient
}

func GetMongo() (*Mongo, error) {
	mongoConfig := DotEnv{}.GetMongoConfig()
	db, err := mgo.Dial(mongoConfig.Connection)
	if err != nil {
		return nil, err
	}
	return &Mongo{
		session: db,
	}, nil
}

func (m *Mongo) Close() {
	m.session.Close()
}

func (m *Mongo) Copy() *mgo.Session {
	return m.session.Copy()
}

func (m *Mongo) Clone() *mgo.Session {
	return m.session.Clone()
}
