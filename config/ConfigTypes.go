package config

import "time"

type RedisConfig struct {
	Host     string
	Password string
	Db       int
}
type MongoConfig struct {
	Connection string
	DbName     string
}

type ApplicationConfig struct {
	SlackLoggingAppConfig string
	JWTExpiryTime         time.Duration
	JWTSecret             string
	RefreshTokenSize      int // should always be more than 128 for security reason
	SlackLogLevel         string
}
