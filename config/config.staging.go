// +build staging

package config

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
