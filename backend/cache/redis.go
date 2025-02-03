package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func SetCache(client *redis.Client, key string, value interface{}, expiration time.Duration) error {
	return client.Set(Ctx, key, value, expiration).Err()
}

func GetCache(client *redis.Client, key string) (string, error) {
	return client.Get(Ctx, key).Result()
}

func DeleteCache(client *redis.Client, key string) error {
	return client.Del(Ctx, key).Err()
}