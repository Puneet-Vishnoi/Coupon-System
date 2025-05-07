package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type Db struct {
	PostgresClient *sql.DB
}

// ConnectDB establishes a connection to the PostgreSQL database
func ConnectDB() *Db {
	connStr := "host=localhost port=5432 user=postgres password=Puneet dbname=coupon-system sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)

	fmt.Println("Connected to PostgreSQL database successfully!")

	return &Db{
		PostgresClient: db,
	}
}

// Stop gracefully closes the PostgreSQL connection
func (db *Db) Stop() {
	if db.PostgresClient != nil {
		err := db.PostgresClient.Close()
		if err != nil {
			log.Printf("Error closing PostgreSQL connection: %v", err)
		} else {
			fmt.Println("PostgreSQL connection closed successfully!")
		}
	}
}

// InitSchema creates the necessary tables in the PostgreSQL database
func (db *Db) InitSchema() error {
	schema := `
	DROP TABLE IF EXISTS coupon_usages;
	DROP TABLE IF EXISTS coupons;

	CREATE TABLE coupons (
		coupon_code TEXT PRIMARY KEY,
		expiry_date TIMESTAMPTZ NOT NULL,
		usage_type TEXT NOT NULL DEFAULT 'single_use' CHECK (usage_type IN ('single_use', 'multi_use')),
		applicable_medicine_ids JSONB NOT NULL DEFAULT '[]',
		applicable_categories JSONB NOT NULL DEFAULT '[]',
		min_order_value DOUBLE PRECISION NOT NULL DEFAULT 0,
		valid_start TIMESTAMPTZ NOT NULL,
		valid_end TIMESTAMPTZ NOT NULL,
		terms_and_conditions TEXT NOT NULL DEFAULT '',
		discount_type TEXT NOT NULL DEFAULT 'flat' CHECK (discount_type IN ('flat', 'percentage')),
		discount_value DOUBLE PRECISION NOT NULL DEFAULT 0,
		max_usage_per_user INTEGER NOT NULL DEFAULT 1,
		discount_target TEXT NOT NULL DEFAULT 'total_order_value' CHECK (discount_target IN ('medicine', 'delivery', 'total_order_value')),
		max_discount_amount DOUBLE PRECISION NOT NULL DEFAULT 0
	);

	CREATE TABLE coupon_usages (
		id SERIAL PRIMARY KEY,
		user_id TEXT NOT NULL,
		coupon_code TEXT NOT NULL REFERENCES coupons(coupon_code) ON DELETE CASCADE,
		used_at TIMESTAMPTZ DEFAULT NOW()
	);`

	_, err := db.PostgresClient.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	fmt.Println("Database schema initialized successfully.")
	return nil
}