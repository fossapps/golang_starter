package helpers

import (
	"errors"
	"github.com/globalsign/mgo"
)

func QueueDeviceRegistration(registrationToken string, redisClient IRedisClient) error {
	// todo error handling
	result, _ := redisClient.SIsMember("registration", registrationToken)
	if result {
		return errors.New("token already exists")
	}
	// todo also check if it's already in MongoDb
	redisClient.SAdd("registration", registrationToken)
	// SPop will pop it out.
	return nil
}

func RegisterDevice(registrationToken string, db *mgo.Database) {
	// todo check if registration token already exists
	db.C("devices").Insert(struct {
		Token string `json:"token"`
	}{Token: registrationToken})
}
