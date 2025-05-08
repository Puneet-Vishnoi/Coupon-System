package redis

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisDb struct {
	RedisClient *redis.Client
}

func ConnectRedis() *RedisDb {
	ctx := context.Background()

	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	dbIndex, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	var redisClient *redis.Client
	var err error

	for i := 0; i < 5; i++ {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       dbIndex,
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
	}

	log.Print("Failed to connect to Redis after multiple attempts")
	return &RedisDb{}
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
