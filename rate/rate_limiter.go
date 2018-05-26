package rate

import (
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

// RedisClient interface need to satisfy for Limiter to work
type RedisClient interface {
	ZAdd(key string, members ...redis.Z) (int64, error)
	ZRemRangeByScore(key, min, max string) (int64, error)
	Expire(key string, expiration time.Duration) (bool, error)
	ZCard(key string) (int64, error)
}

// Limiter helper for making sure we don't call something too much
type Limiter struct {
	RedisClient RedisClient
	Decay       time.Duration
	Limit       int
}

// Count returns number of attempts done till now
func (limiter Limiter) Count(key string) (int64, error) {
	duration := strconv.FormatInt(time.Now().Add(-limiter.Decay).Unix(), 10)
	_, err := limiter.RedisClient.ZRemRangeByScore(key, "0", duration)
	if err != nil {
		return -1, err
	}
	return limiter.RedisClient.ZCard(key)
}

// Hit adds a hit to redis and finally returns a new Count
func (limiter Limiter) Hit(key string) (int64, error) {
	_, err := limiter.RedisClient.ZAdd(key, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: time.Now().Unix(),
	})
	if err != nil {
		return -1, err
	}

	_, err = limiter.RedisClient.Expire(key, limiter.Decay)
	if err != nil {
		return -1, err
	}

	return limiter.Count(key)
}
