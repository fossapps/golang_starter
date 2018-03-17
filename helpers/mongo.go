package helpers

import (
	"github.com/globalsign/mgo"
	"crazy_nl_backend/config"
)

func GetMongo(config *config.MongoConfig) (*mgo.Session, error) {
	return mgo.Dial(config.Connection)
}
