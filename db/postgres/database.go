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
	connStr := "host=localhost port=5432 user=postgres password=Puneet dbname=food-delivery-backend sslmode=disable"

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
