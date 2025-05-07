package redis

import (
	"context"
	"log"
	"time"

	constant "github.com/Puneet-Vishnoi/Coupon-System/models/constants"
	"github.com/go-redis/redis/v8"
)

type RedisDb struct {
	RedisClient *redis.Client
}

func ConnectRedis() *RedisDb {
	ctx := context.Background()

	var redisClient *redis.Client
	var err error

	// Retry logic for Redis connection
	for i := 0; i < constant.MAX_DB_ATTEMPTS; i++ {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     constant.RedisAddr,
			Password: constant.RedisPassword,
			DB:       constant.RedisDB,
		})

		err = redisClient.Ping(ctx).Err()
		if err != nil {
			log.Printf("Redis connection attempt %d failed: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	if err == nil {
		log.Print("Redis connected successfully")
		return &RedisDb{RedisClient: redisClient}
	} else {
		log.Print("Failed to connect to Redis after maximum attempts")
		return &RedisDb{}
	}
}

func (db *RedisDb) Stop() {
	if db.RedisClient == nil {
		log.Println("Redis client is nil, skipping stop.")
		return
	}

	ctx := context.Background()

	if err := db.RedisClient.FlushAll(ctx).Err(); err != nil && err != redis.ErrClosed {
		log.Printf("Failed to flush Redis: %v", err)
	}

	if err := db.RedisClient.Close(); err != nil && err != redis.ErrClosed {
		log.Printf("Error closing Redis connection: %v", err)
	} else {
		log.Print("Redis connection closed successfully")
	}
}