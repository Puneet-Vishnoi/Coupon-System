package mockdb

import (
	"log"

	"github.com/Puneet-Vishnoi/Coupon-System/cache/redis"
	redisProvider "github.com/Puneet-Vishnoi/Coupon-System/cache/redis/providers"
	"github.com/Puneet-Vishnoi/Coupon-System/db/postgres"
	providers "github.com/Puneet-Vishnoi/Coupon-System/db/postgres/providers"
	"github.com/Puneet-Vishnoi/Coupon-System/repository"
	"github.com/Puneet-Vishnoi/Coupon-System/service"
)

type TestDeps struct {
	Service        *service.CouponService
	Repo           *repository.CouponRepository
	PostgresClient *postgres.Db
	RedisClient    *redis.RedisDb
	Cleanup        func()
}

// Returns an instance of initialized test services and clients
func GetTestInstance() *TestDeps {
	// 1. Connect to Postgres (you can point this to a dedicated test DB)
	pgClient := postgres.ConnectDB()
	err := pgClient.InitSchema() 
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
	redisClient := redis.ConnectRedis()
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
