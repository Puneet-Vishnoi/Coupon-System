package providers

import (
	redis "github.com/go-redis/redis/v8"
)

type RedisHelper struct {
	RedisClient *redis.Client
}

func NewRedisProvider(redisClient *redis.Client) *RedisHelper {
	return &RedisHelper{
		RedisClient: redisClient,
	}
}
