package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Set up the database connection string
	connStr := "host=localhost port=5432 user=postgres password=your_password dbname=market_data_db sslmode=disable"

	// Attempt to connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")
}
