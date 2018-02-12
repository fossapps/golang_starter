package helpers

import (
	"github.com/go-redis/redis"
)

func GetRedis() (*Redis, error){
	config, err := DotEnv{}.GetRedisConfig()
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr: config.Host,
		Password: config.Password,
		DB: config.Db,
	})
	return &Redis{
		Client: client,
	}, nil
}

type IRedisClient interface {
	SIsMember(string, interface{}) (bool, error)
	SAdd(string, ...interface{}) (int64, error)
	Close() error
}

type Redis struct {
	Client *redis.Client
}

func (channel *Redis) SIsMember(key string, member interface{}) (bool, error) {
	return channel.Client.SIsMember(key, member).Result()
}

func (channel *Redis) SAdd(key string, member ...interface{}) (int64, error) {
	return channel.Client.SAdd(key, member...).Result()
}

func (channel *Redis) Close() error {
	return channel.Client.Close()
}
