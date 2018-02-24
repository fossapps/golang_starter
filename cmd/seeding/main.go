package main

import (
	"crazy_nl_backend/seeds"
	"crazy_nl_backend/config"
)

func main() {
	seeds.SeedDb(config.GetMongoConfig().DbName)
}
