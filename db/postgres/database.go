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

	// Open the database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	// Test the database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL database: %v", err)
	}

	// Configure the connection pool
	db.SetMaxOpenConns(25)   // Maximum number of open connections
	db.SetMaxIdleConns(5)    // Maximum number of idle connections
	db.SetConnMaxLifetime(0) // Unlimited connection lifetime

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
		expiry_date TIMESTAMPTZ,
		usage_type TEXT,
		applicable_medicine_ids JSONB,
		applicable_categories JSONB,
		min_order_value DOUBLE PRECISION,
		valid_start TIMESTAMPTZ,
		valid_end TIMESTAMPTZ,
		terms_and_conditions TEXT,
		discount_type TEXT,
		discount_value DOUBLE PRECISION,
		max_usage_per_user INTEGER,
		discount_target TEXT,
		max_discount_amount DOUBLE PRECISION
	);

	CREATE TABLE coupon_usages (
		id SERIAL PRIMARY KEY,
		user_id TEXT,
		coupon_code TEXT REFERENCES coupons(coupon_code) ON DELETE CASCADE,
		used_at TIMESTAMPTZ DEFAULT NOW()
	);`

	_, err := db.PostgresClient.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	fmt.Println("Database schema initialized successfully.")
	return nil
}