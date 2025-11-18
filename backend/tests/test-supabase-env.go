package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main2() {
	// Load .env file (optional path: .env by default)
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: No .env file found or could not load it")
	}

	// Now this will work
	dbURL := os.Getenv("DATABASE_URL")
	fmt.Println("DATABASE_URL:", dbURL)

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("Connected to:", version)
}
