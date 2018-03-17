package config

import (
	"time"
	"os"
	"strconv"
)

func GetRedisConfig() *RedisConfig {
	redisDb, _ := strconv.Atoi(os.Getenv("REDIS_DB_VALUE"))
	return &RedisConfig{
		Host:     os.Getenv("REDIS_DB"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Db:       redisDb,
	}
}

func GetMongoConfig() *MongoConfig {
	return &MongoConfig{
		Connection: os.Getenv("MONGO_DB_CONNECTION"),
		DbName:     os.Getenv("MONGO_DB_NAME"),
	}
}

func GetLogLevel() string {
	return os.Getenv("LOG_LEVEL")
}

func GetApiPort() int {
	port, _ := strconv.Atoi(os.Getenv("API_PORT"))
	return port
}

func GetPushyToken() string {
	return os.Getenv("PUSHY_TOKEN")
}

func GetTestingDbConnection() string {
	return os.Getenv("MGO_TEST_CONNECTION")
}

func GetApplicationConfig() ApplicationConfig {
	validity, _ := strconv.Atoi(os.Getenv("JWT_VALIDITY"))
	refreshTokenSize, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_SIZE"))
	return ApplicationConfig{
		JWTExpiryTime:         time.Duration(validity) * time.Minute,
		JWTSecret:             os.Getenv("JWT_SECRET"),
		RefreshTokenSize:      refreshTokenSize,
		SlackLoggingAppConfig: os.Getenv("SLACK_WEBHOOK"),
		SlackLogLevel:         os.Getenv("SLACK_LOG_LEVEL"),
	}
}
