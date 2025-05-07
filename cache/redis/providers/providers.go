package providers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisHelper struct {
	RedisClient *redis.Client
}

func NewRedisProvider(redisClient *redis.Client) *RedisHelper {
	return &RedisHelper{
		RedisClient: redisClient,
	}
}

func (r *RedisHelper) GetJSON(ctx context.Context, key string, dest interface{}) (bool, error) {
	val, err := r.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil 
	}
	if err != nil {
		return false, err
	}
	err = json.Unmarshal([]byte(val), dest)
	return err == nil, err
}

func (r *RedisHelper) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.RedisClient.Set(ctx, key, data, ttl).Err()
}


func (r *RedisHelper) Delete(ctx context.Context, key string) error {
	return r.RedisClient.Del(ctx, key).Err()
}
