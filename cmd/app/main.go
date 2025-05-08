package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/Puneet-Vishnoi/Coupon-System/cache/redis"
	redisProvider "github.com/Puneet-Vishnoi/Coupon-System/cache/redis/providers"
	"github.com/Puneet-Vishnoi/Coupon-System/db/postgres"
	providers "github.com/Puneet-Vishnoi/Coupon-System/db/postgres/providers"
	"github.com/Puneet-Vishnoi/Coupon-System/repository"
	"github.com/Puneet-Vishnoi/Coupon-System/routes"
	couponService "github.com/Puneet-Vishnoi/Coupon-System/service"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to load env file: ", err)
	}

	// 1. Connect Redis
	redisClient := redis.ConnectRedis()
	redisHelper := redisProvider.NewRedisProvider(redisClient.RedisClient)

	// 2. Connect PostgreSQL
	postgresClient := postgres.ConnectDB()
	defer postgresClient.Stop()

	// 2.1 Init Schema (optional — only for development)
	if err := postgresClient.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	// 3. DB Helper
	dbHelper, err := providers.NewDbProvider(postgresClient.PostgresClient)
	if err != nil {
		log.Fatalf("Failed to initialize DB helper: %v", err)
	}

	// 4. Repo & Service
	couponRepo := repository.NewCouponRepository(dbHelper)
	couponSrv := couponService.NewCouponService(couponRepo, redisHelper)

	// 5. Gin Router & Handlers
	router := gin.Default()
	routes.RegisterRoutes(router, couponSrv)

	// 6. Run REST API
	port := ":8080"
	fmt.Printf("Coupon REST API running on %s\n", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start Gin server: %v", err)
	}
}
