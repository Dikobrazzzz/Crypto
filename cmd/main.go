package main

import (
	"context"
	"log/slog"
	"os"

	"crypto/config"
	migrations "crypto/database"
	"crypto/internal/app"
	"crypto/internal/cache"
	"crypto/internal/handler"
	"crypto/internal/repository"
	"crypto/internal/storage"
	"crypto/internal/usecase"
)

func loggerinit() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})))
}

func main() {
	config.Init()
	dbURL := config.AppConfig.DatabaseURL
	if err := migrations.Migrate(dbURL); err != nil {
		slog.Error("Migration failed", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := storage.GetConnection(ctx, "postgres://postgres:postgres@db:5432/postgres")
	if err != nil {
		slog.Error("Error getting connection", "error", err)
		return
	}
	defer pool.Close()

	walletRepo := repository.NewWalletProvider(pool)
	cacheDecorator := cache.CacheNewDecorator(walletRepo)
	walletUC := usecase.NewWalletProvider(cacheDecorator)
	handle := handler.New(walletUC)

	router := app.GetRouter(handle)

	if err := router.Run(); err != nil {
		slog.Error("Error", "error", err)
	}
}
