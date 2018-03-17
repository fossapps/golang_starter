package workers

import (
	"fmt"
	"time"

	"github.com/cyberhck/captain"
	"crazy_nl_backend/helpers"
	"github.com/cyberhck/pushy"
	"github.com/globalsign/mgo"
)

type ImportDevices struct {
	Db    *mgo.Database
	Redis helpers.IRedisClient
	Pushy *pushy.Pushy
}

func (job ImportDevices) LockProvider() captain.LockProvider {
	return nil
}

func (job ImportDevices) Job() captain.Worker {
	return func(channels captain.CommChan) {
		members, err := job.Redis.SMembers("registration")
		if err != nil {
			channels.Result <- "error: " + err.Error()
			return
		}
		for _, member := range members {
			job.importDevice(member)
		}
	}
}

func (job ImportDevices) ResultProcessor() captain.ResultProcessor {
	return nil
}

func (job ImportDevices) RuntimeProcessor() captain.RuntimeProcessor {
	return func(tick time.Time, message string, startTime time.Time) {
	}
}

func (job ImportDevices) ShouldRun(key string) bool {
	return key == "import_devices"
}

func (job ImportDevices) importDevice(token string) {
	if !job.isDeviceValid(token) {
		fmt.Println("invalid token: ", token)
		return
	}
	helpers.RegisterDevice(token, job.Db)
	// delete from redis
}

func (job ImportDevices) isDeviceValid(token string) bool {
	_, pushyErr, err := job.Pushy.DeviceInfo(token)
	if pushyErr != nil || err != nil {
		return false
	}
	return true
}
