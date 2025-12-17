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

type DatabaseConf struct {
	conn *pgxpool.Pool
}

func (dc *DatabaseConf) NewConn() {
	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		slog.Error(fmt.Sprintf("Unable to create connection pool: %v", err))
		panic("Unable to connect to database")
	}

	slog.Info("Database connection established successfully")

	dc.conn = conn
}

func (dc *DatabaseConf) GetConn() *pgxpool.Pool {
	return dc.conn
}
