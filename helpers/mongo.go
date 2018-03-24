package helpers

import (
	"crazy_nl_backend/config"
	"github.com/globalsign/mgo"
)

func GetMongo(config *config.MongoConfig) (*mgo.Session, error) {
	return mgo.Dial(config.Connection)
}
