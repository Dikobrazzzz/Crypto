package storage

import (
	"log/slog"
	"os"
	"github.com/jackc/pgx/v5"
	"context"
)

func GetConnection(ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("DB connection error", "error", err)
		return nil, err
	}
	return conn, nil
}