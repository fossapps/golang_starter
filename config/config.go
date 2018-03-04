// +build dev

package config

import "time"

func GetRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host: "localhost:6379",
		Password: "",
		Db: 0,
	}
}

func GetMongoConfig() *MongoConfig {
	return &MongoConfig{
		Connection: "localhost:27017",
		DbName: "crazy_nl",
	}
}

func GetLogLevel() string {
	return "info"
}

func GetApiPort() int {
	return 8080
}

func GetPushyToken() string {
	return "9a64fc25eb0dee5c0fcb88c6dbc033041a919024279814489fd12c5906184eae"
}

func GetTestingDbName() string {
	return "integration_tests"
}

func GetApplicationConfig() ApplicationConfig {
	return ApplicationConfig{
		JWTExpiryTime:10 * time.Minute,
		JWTSecret:"SECRET",
		RefreshTokenSize:256,
		SlackLoggingAppConfig:"https://hooks.slack.com/services/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",
	}
}
