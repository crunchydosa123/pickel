package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func ConnectSupabase() {
	_ = godotenv.Load()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set in environment")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to Supabase: %v", err)
	}

	// Optional: ping to ensure connectivity
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Ping to Supabase failed: %v", err)
	}

	DB = pool
	fmt.Println("Connected to Supabase Postgres")
}

func GetDB() *pgxpool.Pool {
	if DB == nil {
		log.Fatal("Supabase DB not initialized. Call ConnectSupabase() first.")
	}

	return DB
}
