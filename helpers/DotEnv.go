package helpers

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"errors"
)

func InitDotEnv() error {
	return godotenv.Load(".env")
}

type DotEnv struct{}

type RedisConfig struct {
	Host     string
	Password string
	Db       int
}

type MongoConfig struct {
	Connection string
}

func (DotEnv) GetRedisConfig() (*RedisConfig, error) {
	host := os.Getenv("REDIS_HOST")
	password := os.Getenv("REDIS_PASSWORD")
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, errors.New("REDIS_DB is invalid")
	}
	return &RedisConfig{
		Host:     host,
		Password: password,
		Db:       db,
	}, nil
}

func (DotEnv) GetMongoConfig() *MongoConfig {
	return &MongoConfig{
		Connection: os.Getenv("DB_CONNECTION"),
	}
}

func (DotEnv) GetLogLevel() string {
	return os.Getenv("LOG_LEVEL")
}
func (DotEnv) GetServerPort() string {
	return os.Getenv("API_PORT")
}