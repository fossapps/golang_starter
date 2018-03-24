package main

import (
	"crazy_nl_backend/config"
	"crazy_nl_backend/helpers"
	"crazy_nl_backend/migrations"
)

func main() {
	session, err := helpers.GetMongo(config.GetMongoConfig())
	if err != nil {
		panic(err)
	}
	migrations.ApplyAll(config.GetMongoConfig().DbName, session)
}
