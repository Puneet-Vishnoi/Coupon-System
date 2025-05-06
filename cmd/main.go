package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Puneet-Vishnoi/Coupon-System/coupon-system/routes"

	config "github.com/Puneet-Vishnoi/Coupon-System/db"
)

func initPostgres(ctx context.Context) *pgxpool.Pool {
	dsn := os.Getenv("POSTGRES_DSN")
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to PostgreSQL: %v", err)
	}
	return pool
}

func main() {
	ctx := context.Background()

	// Initialize PostgreSQL
	db := initPostgres(ctx)
	defer db.Close()

	// Initialize Redis
	config.InitRedis()

	// Initialize Gin
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r, db)

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Start server
	r.Run() // default :8080
}
