package helpers

import (
	"github.com/globalsign/mgo"
	"crazy_nl_backend/config"
)

type Mongo struct {
	session *mgo.Session
}

type IMongoClient interface {
	Close()
	Clone() IMongoClient
	Copy() IMongoClient
	DB(string) IDatabase
}

func GetMongo() (IMongoClient, error) {
	mongoConfig := config.GetMongoConfig()
	session, err := mgo.Dial(mongoConfig.Connection)
	if err != nil {
		return nil, err
	}
	return &Mongo{
		session: session,
	}, nil
}

func (m *Mongo) Close() {
	m.session.Close()
}

func (m *Mongo) Copy() IMongoClient {
	return &Mongo{
		session: m.session.Copy(),
	}
}

func (m *Mongo) Clone() IMongoClient {
	return &Mongo{
		session: m.session.Clone(),
	}
}

func (m *Mongo) DB(name string) IDatabase {
	return Db{
		DB: m.session.DB(name),
	}
}
