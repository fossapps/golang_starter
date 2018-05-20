package helpers

import (
	"github.com/fossapps/starter/config"
	"github.com/go-redis/redis"
	"time"
)

// GetRedis returns an implementation of IRedisClient
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

// IRedisClient interface of redis which this application depends on
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

// Redis implementation of IRedisClient
type Redis struct {
	Client *redis.Client
}

// ZCard returns cardinal number of a sorted set
func (channel *Redis) ZCard(key string) (int64, error) {
	return channel.Client.ZCard(key).Result()
}

// ZAdd adds a member to sorted set
func (channel *Redis) ZAdd(key string, members ...redis.Z) (int64, error) {
	return channel.Client.ZAdd(key, members...).Result()
}

// ZRemRangeByScore removes members of a sorted set between min and max scores
func (channel *Redis) ZRemRangeByScore(key, min, max string) (int64, error) {
	return channel.Client.ZRemRangeByScore(key, min, max).Result()
}

// Expire sets time to live on a key
func (channel *Redis) Expire(key string, expiration time.Duration) (bool, error) {
	return channel.Client.Expire(key, expiration).Result()
}

// SIsMember checks if a a member exists in a set
func (channel *Redis) SIsMember(key string, member interface{}) (bool, error) {
	return channel.Client.SIsMember(key, member).Result()
}

// SAdd adds a member to a set
func (channel *Redis) SAdd(key string, member ...interface{}) (int64, error) {
	return channel.Client.SAdd(key, member...).Result()
}

// Close closes connection to redis
func (channel *Redis) Close() error {
	return channel.Client.Close()
}

// SPop pops an element out of set
func (channel *Redis) SPop(key string) (string, error) {
	return channel.Client.SPop(key).Result()
}

// LRange returns element between start and stop
func (channel *Redis) LRange(key string, start int64, stop int64) ([]string, error) {
	return channel.Client.LRange(key, start, stop).Result()
}

// SRem remove members from a set
func (channel *Redis) SRem(key string, members ...interface{}) (int64, error) {
	return channel.Client.SRem(key, members...).Result()
}

// SMembers returns members of set
func (channel *Redis) SMembers(key string) ([]string, error) {
	return channel.Client.SMembers(key).Result()
}
