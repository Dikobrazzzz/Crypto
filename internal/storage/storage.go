package storage

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetConnection(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		slog.Error("DB connection error", "error", err)
		return nil, err
	}
	return pool, nil
}
