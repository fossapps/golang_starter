package main

import (
	"flag"
	"fmt"

	"crazy_nl_backend/config"
	"crazy_nl_backend/helpers"
	"crazy_nl_backend/workers"
	"github.com/cyberhck/pushy"
	"time"
)

func main() {
	job := flag.String("job", "", "job to run")
	flag.Parse()
	fmt.Println(*job)
	if *job == "" {
		flag.Usage()
	}
	importDevices(*job)
}

func importDevices(key string) {
	redis, err := helpers.GetRedis()
	if err != nil {
		fmt.Println(err)
		return
	}
	db, err := helpers.GetMongo(config.GetMongoConfig())
	if err != nil {
		fmt.Println(err)
		return
	}
	sdk := pushy.Create(config.GetPushyToken(), pushy.GetDefaultAPIEndpoint())
	sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(10 * time.Second))
	job := workers.ImportDevices{
		Redis: redis,
		Db:    db.DB(config.GetMongoConfig().DbName),
		Pushy: sdk,
	}
	workers.Run(key, job)
}
