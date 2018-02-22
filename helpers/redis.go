package helpers

import (
	"github.com/go-redis/redis"
	"crazy_nl_backend/config"
)

func GetRedis() (IRedisClient, error){
	redisConfig := config.GetRedisConfig()
	client := redis.NewClient(&redis.Options{
		Addr: redisConfig.Host,
		Password: redisConfig.Password,
		DB: redisConfig.Db,
	})
	return &Redis{
		Client: client,
	}, nil
}

type IRedisClient interface {
	SIsMember(string, interface{}) (bool, error)
	SAdd(string, ...interface{}) (int64, error)
	Close() error
	SPop(key string) (string, error)
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

func (channel *Redis) SPop(key string) (string, error) {
	return channel.Client.SPop(key).Result()
}
