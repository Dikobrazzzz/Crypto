package main

import (
	"context"
	"os"
	"log/slog"

	"crypto/internal/app"
	"crypto/internal/handler"
	"crypto/migrations"
	"crypto/internal/repository"
	"crypto/internal/storage"
	usecase "crypto/internal/usecase"
)

func loggerinit() {
    slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    })))
}

func main() {

	dbURL := "postgres://postgres:postgres@postgres:5432"
	if err := migrate.Migrate(dbURL); err != nil {
        slog.Error("Migration failed", "error", err)
        os.Exit(1)
    }
	
	ctx := context.Background()
	conn, err := storage.GetConnection(ctx)
    if err != nil {
        slog.Error("Error getting connection", "error", err)
        return
    }
    defer conn.Close(ctx)

	walletRepo := repository.NewWalletProvider(conn)
	walletUC := usecase.NewWalletProvider(walletRepo)
	handle := handler.New(walletUC)

	router:= app.GetRouter(handle)

	if err:= router.Run(); err != nil {
		slog.Error("Error","error", err)
	}
}
