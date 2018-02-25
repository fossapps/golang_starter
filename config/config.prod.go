// +build prod

package config

import "time"

func GetRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host: "",
		Password: "",
		Db: 0,
	}
}

func GetMongoConfig() *MongoConfig {
	return &MongoConfig{
		Connection: "",
		DbName: "crazy_nl",
	}
}

func GetLogLevel() string {
	return "debug"
}

func GetApiPort() int {
	return 80
}

func GetPushyToken() string {
	return "9a64fc25eb0dee5c0fcb88c6dbc033041a919024279814489fd12c5906184eae"
}

func GetTestingDbName() string {
	panic("shouldn't be running in production mode...")
}

func GetApplicationConfig() ApplicationConfig {
	return ApplicationConfig{
		JWTExpiryTime:10 * time.Minute,
		JWTSecret:"My_Prod_SECRET",
		RefreshTokenSize:256,
	}
}
