package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseGetter interface {
	NewConn()
	GetConn() *pgxpool.Pool
}

type Database struct {
	Pool *pgxpool.Pool
}

func New() *Database {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		slog.Error(fmt.Sprintf("Unable to create connection pool: %v", err))
		panic("Unable to connect to database")
	}

	slog.Info("Database connection established successfully")

	return &Database{Pool: pool}
}

func (db *Database) GetPool() *pgxpool.Pool {
	return db.Pool
}
