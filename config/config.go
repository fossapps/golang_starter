package config

import (
	"os"
	"strconv"
	"time"
)

// GetRedisConfig returns redis config
func GetRedisConfig() *RedisConfig {
	redisDb, _ := strconv.Atoi(os.Getenv("REDIS_DB_VALUE"))
	return &RedisConfig{
		Host:     os.Getenv("REDIS_DB"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Db:       redisDb,
	}
}

// GetMongoConfig returns mongodb config
func GetMongoConfig() *MongoConfig {
	return &MongoConfig{
		Connection: os.Getenv("MONGO_DB_CONNECTION"),
		DbName:     os.Getenv("MONGO_DB_NAME"),
	}
}

// GetLogLevel returns log level to display
func GetLogLevel() string {
	return os.Getenv("LOG_LEVEL")
}

// GetAPIPort returns port to run on
func GetAPIPort() int {
	port, _ := strconv.Atoi(os.Getenv("API_PORT"))
	return port
}

// GetPushyToken returns token to access pushy API
func GetPushyToken() string {
	return os.Getenv("PUSHY_TOKEN")
}

// GetTestingDbName db name to use for integration testing
func GetTestingDbName() string {
	return "crazy_nl_test_db"
}

// GetApplicationConfig returns application config
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
