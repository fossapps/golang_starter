package helpers

import (
	"time"
	//"crazy_nl_backend"
)

func Remember(key string, handler func() string ,duration time.Duration) string {
	// check if the key is there.
	// check if it's expired
	// if everything is ok, return
	// if expired/not exists, call handler func
	//redis := crazy_nl_backend.GetRedis()
	//result, _ := redis.Get("cache" + key).Result()
	//redis.Set("cache|" + key, handler(), duration).Result()
	return ""
}
