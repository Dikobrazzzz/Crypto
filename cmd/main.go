package main

import (
	"context"
	"log/slog"
	"os"

	"crypto/config"
	"crypto/database"
	"crypto/internal/app"
	"crypto/internal/cache"
	"crypto/internal/handler"
	"crypto/internal/metrics"
	"crypto/internal/repository"
	"crypto/internal/storage"
	"crypto/internal/usecase"
)

func loggerinit() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level:     config.AppConfig.Level,
		AddSource: true,
	})))
}

func main() {

	config.Init()
	if err := database.Migrate(config.AppConfig.DatabaseURL); err != nil {
		slog.Error("Migration failed", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := storage.GetConnection(ctx, config.AppConfig.DatabaseURL)
	if err != nil {
		slog.Error("Error getting connection", "error", err)
		return
	}
	defer pool.Close()

	walletRepo := repository.NewWalletProvider(pool)
	cacheDecorator := cache.CacheNewDecorator(walletRepo, config.AppConfig.TTL)
	metrics.InitMetrics("9090", cacheDecorator)
	walletUC := usecase.NewWalletProvider(cacheDecorator)
	handle := handler.New(walletUC)

	router := app.GetRouter(handle)

	if err := router.Run(); err != nil {
		slog.Error("Error", "error", err)
	}
}
