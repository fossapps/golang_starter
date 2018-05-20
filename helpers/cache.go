package helpers

import (
	"time"
)

// Remember is for caching something
// todo: implement this
func Remember(key string, handler func() string, duration time.Duration) string {
	// check if the key is there.
	// check if it's expired
	// if everything is ok, return
	// if expired/not exists, call handler func
	// redis := starter.GetRedis()
	// result, _ := redis.Get("cache" + key).Result()
	// redis.Set("cache|" + key, handler(), duration).Result()
	return ""
}
