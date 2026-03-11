package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect() {

	connStr := os.Getenv("DATABASE_URL")

	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/coupons?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	DB = pool

	log.Println("Connected to PostgreSQL")
}
