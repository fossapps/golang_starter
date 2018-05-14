package helpers

import (
	"crazy_nl_backend/config"
	"github.com/go-redis/redis"
	"time"
)

func GetRedis() (IRedisClient, error) {
	redisConfig := config.GetRedisConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Host,
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
	})
	return &Redis{
		Client: client,
	}, nil
}

type IRedisClient interface {
	SIsMember(string, interface{}) (bool, error)
	SAdd(string, ...interface{}) (int64, error)
	Close() error
	SPop(string) (string, error)
	SRem(string, ...interface{}) (int64, error)
	SMembers(string) ([]string, error)
	ZAdd(key string, members ...redis.Z) (int64, error)
	ZRemRangeByScore(key, min, max string) (int64, error)
	ZCard(key string) (int64, error)
	Expire(key string, expiration time.Duration) (bool, error)
}

type Redis struct {
	Client *redis.Client
}

func (channel *Redis) ZCard(key string) (int64, error) {
	return channel.Client.ZCard(key).Result()
}

func (channel *Redis) ZAdd(key string, members ...redis.Z) (int64, error) {
	return channel.Client.ZAdd(key, members...).Result()
}

func (channel *Redis) ZRemRangeByScore(key, min, max string) (int64, error) {
	return channel.Client.ZRemRangeByScore(key, min, max).Result()
}

func (channel *Redis) Expire(key string, expiration time.Duration) (bool, error) {
	return channel.Client.Expire(key, expiration).Result()
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

func (channel *Redis) LRange(key string, start int64, stop int64) ([]string, error) {
	return channel.Client.LRange(key, start, stop).Result()
}

func (channel *Redis) SRem(key string, members ...interface{}) (int64, error) {
	return channel.Client.SRem(key, members...).Result()
}

func (channel *Redis) SMembers(key string) ([]string, error) {
	return channel.Client.SMembers(key).Result()
}
