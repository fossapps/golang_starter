// +build staging

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
	return 8080
}

func GetPushyToken() string {
	return "9a64fc25eb0dee5c0fcb88c6dbc033041a919024279814489fd12c5906184eae"
}

func GetTestingDbName() string {
	return "integration_tests_staging"
}

func GetApplicationConfig() ApplicationConfig {
	return ApplicationConfig{
		JWTExpiryTime:10 * time.Minute,
		JWTSecret:"MY_STAGING_SECRET",
		RefreshTokenSize:256,
		SlackLoggingAppConfig:"https://hooks.slack.com/services/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",
	}
}
