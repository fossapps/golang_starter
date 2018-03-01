package workers

import (
	"fmt"
	"time"

	"crazy_nl_backend/config"

	"github.com/cyberhck/captain"

	"crazy_nl_backend/helpers"
	// "github.com/cyberhck/pushy"
)

type ImportDevices struct {
	Db    helpers.IDatabase
	Redis helpers.IRedisClient
}

func (job ImportDevices) LockProvider() captain.LockProvider {
	return nil
}

func (job ImportDevices) Job() captain.Worker {
	return func(channels captain.CommChan) {
		// start a while pop doesn't return nil, (maybe don't pop, but if popped, make sure it's imported?)
		// invoke import on a goroutine which will take care of everything (maybe a maximum of 100 go routines?)
		channels.Logs <- "working..."
		time.Sleep(5 * time.Second)
	}
}

func (job ImportDevices) ResultProcessor() captain.ResultProcessor {
	return nil
}

func (job ImportDevices) RuntimeProcessor() captain.RuntimeProcessor {
	return func(tick time.Time, message string, startTime time.Time) {
		fmt.Println(tick, message)
	}
}

func (job ImportDevices) ShouldRun(key string) bool {
	return key == "import_devices"
}

func ProcessQueue(channels captain.CommChan) {

	db, redis, err := getDbAndRedis()

	if err != nil {

		channels.Result <- "connection error" + err.Error()

		return

	}

	processAll(redis, db, channels)

}

func getDbAndRedis() (helpers.IMongoClient, helpers.IRedisClient, error) {

	db, err := helpers.GetMongo(config.GetMongoConfig())

	if err != nil {

		return nil, nil, err

	}

	redis, err := helpers.GetRedis()

	if err != nil {

		return nil, nil, err

	}

	return db, redis, nil

}

func processAll(redis helpers.IRedisClient, db helpers.IMongoClient, channels captain.CommChan) {

	//token, err := redis.SPop("registration")

	//if err != nil {

	//	channels.Result <- "registration pop error: " + err.Error()

	//	return

	//}

	//// check if token is correct

	//device := pushy.Device{

	//	Token: token,

	//}

	//if !device.IsValid() {

	//	channels.Logs <- "Invalid device detected: " + token

	//}

	//helpers.RegisterDevice(device.Token, db)

	//// simply save it.

	//channels.Result <- "Valid one: " + token

}
