package main

import (
	"crazy_nl_backend/migrations"
	"crazy_nl_backend/config"
)

func main() {
	migrations.ApplyAll(config.GetMongoConfig().DbName)
}
