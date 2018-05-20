package config

import "time"

// RedisConfig information about redis connection
type RedisConfig struct {
	Host     string
	Password string
	Db       int
}

// MongoConfig information about mongodb connection
type MongoConfig struct {
	Connection string
	DbName     string
}

// ApplicationConfig information about various configuration of application
type ApplicationConfig struct {
	SlackLoggingAppConfig string
	JWTExpiryTime         time.Duration
	JWTSecret             string
	RefreshTokenSize      int // should always be more than 128 for security reason
	SlackLogLevel         string
}
