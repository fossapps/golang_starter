package helpers

import (
	"errors"
)

func QueueDeviceRegistration(registrationToken string, redisClient IRedisClient) error {
	// todo error handling
	result, _ := redisClient.SIsMember("registration", registrationToken)
	if result {
		return errors.New("token already exists")
	}
	redisClient.SAdd("registration", registrationToken)
	// SPop will pop it out.
	return nil
}

func RegisterDevice(registrationToken string, db IMongoClient) {
	session := db.Clone()
	defer session.Close()
}
