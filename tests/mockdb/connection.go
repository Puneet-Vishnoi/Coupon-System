package mockdb

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"

	redisPkg "github.com/Puneet-Vishnoi/Coupon-System/cache/redis"
	redisProvider "github.com/Puneet-Vishnoi/Coupon-System/cache/redis/providers"
	"github.com/Puneet-Vishnoi/Coupon-System/db/postgres"
	providers "github.com/Puneet-Vishnoi/Coupon-System/db/postgres/providers"
	"github.com/Puneet-Vishnoi/Coupon-System/repository"
	"github.com/Puneet-Vishnoi/Coupon-System/service"

	_ "github.com/lib/pq"
)

type TestDeps struct {
	Service        *service.CouponService
	Repo           *repository.CouponRepository
	PostgresClient *postgres.Db
	RedisClient    *redisPkg.RedisDb
	Cleanup        func()
}

// ConnectTestDB connects to the test PostgreSQL DB using TEST_POSTGRES_* env vars
func ConnectTestDB() *postgres.Db {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("TEST_POSTGRES_HOST"),
		os.Getenv("TEST_POSTGRES_PORT"),
		os.Getenv("TEST_POSTGRES_USER"),
		os.Getenv("TEST_POSTGRES_PASSWORD"),
		os.Getenv("TEST_POSTGRES_DB"),
	)
	log.Println(connStr, "//;;")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open test database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to test PostgreSQL DB: %v", err)
	}

	log.Println("Connected to test PostgreSQL database successfully!")
	return &postgres.Db{PostgresClient: db}
}

// ConnectTestRedis connects to Redis using TEST_REDIS_* env vars
func ConnectTestRedis() *redisPkg.RedisDb {
	ctx := context.Background()

	addr := os.Getenv("TEST_REDIS_ADDR")
	password := os.Getenv("TEST_REDIS_PASSWORD")
	dbIndex, _ := strconv.Atoi(os.Getenv("TEST_REDIS_DB"))

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
			log.Printf("Redis test connection attempt %d failed: %v", i+1, err)
			continue
		}
		break
	}

	if err != nil {
		log.Fatalf("Failed to connect to test Redis: %v", err)
	}

	log.Println("Connected to test Redis successfully")
	return &redisPkg.RedisDb{RedisClient: redisClient}
}

// Returns an instance of initialized test services and clients
func GetTestInstance() *TestDeps {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Failed to load env file: ", err)
	}

	// 1. Connect to Postgres (you can point this to a dedicated test DB)
	pgClient := ConnectTestDB()
	err = pgClient.InitSchema()
	if err != nil {
		log.Fatalf("failed to init test schema: %v", err)
	}

	// 2. Setup Postgres Provider
	dbHelper, err := providers.NewDbProvider(pgClient.PostgresClient)
	if err != nil {
		log.Fatalf("failed to get dbHelper: %v", err)
	}
	repo := repository.NewCouponRepository(dbHelper)

	// 3. Connect Redis
	redisClient := ConnectTestRedis()
	redisHelper := redisProvider.NewRedisProvider(redisClient.RedisClient)

	// 4. Build service
	svc := service.NewCouponService(repo, redisHelper)

	return &TestDeps{
		Service:        svc,
		Repo:           repo,
		PostgresClient: pgClient,
		RedisClient:    redisClient,
		Cleanup: func() {
			pgClient.Stop()
			redisClient.Stop()
		},
	}
}
