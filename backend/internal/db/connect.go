package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Pool *pgxpool.Pool

func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env: %v", err)
	}

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("База данных недоступна: %v", err)
	}

	fmt.Println("✅ Подключение к PostgreSQL успешно!")
	Pool = pool
}
