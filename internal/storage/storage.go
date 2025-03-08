package storage

import (
	"log/slog"
	"os"
	"github.com/jackc/pgx/v5"
	"context"
)
	
var (
	conn *pgx.Conn
)

func GetConnection(cb context.Context) (*pgx.Conn, error) {
	c, err := pgx.Connect(cb, os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("DB connection error", "error", err)
		return nil, err
	}
	return c, nil
}