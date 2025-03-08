package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"context"
	"os"
	"log/slog"
	"crypto/internal/handler" 
    "crypto/internal/storage" 
)

var (
	conn 	*pgx.Conn
	cb = context.Background()
)

type Address struct {
	ID            uint64 `json:"id"`
	WalletAddress string `json:"wallet_address"`
	ChainName     string `json:"chain_name"`
	CryptoName    string `json:"crypto_name"`
	Tag           string `json:"tag"`
	Balance       int    `json:"balance"`
}

type AddressRequest struct {
	WalletAddress 	string `json:"wallet_address" binding:"required"`
	ChainName 		string `json:"chain_name" binding:"required"`
	CryptoName 		string `json:"crypto_name" binding:"required"`
	Tag 			string `json:"tag" binding:"required"`
}

type TagUpdateRequest struct {
	ID 	uint64 `json:"id"`
	Tag string `json:"tag"`
}

func init() {
    slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    })))
}

func main() {
	
	conn,err := storage.GetConnection(cb)
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		return
	}
	
	defer conn.Close(cb)

	h := handler.New(conn, cb)

	r := gin.Default()

	r.POST("/address", h.CreateAddressHandler)
	r.GET("/address/:id", h.GetIdHandler)
	r.GET("/allwallets", h.GetAllWalletsHandler)
	r.PUT("/address/tag", h.EditTagHandler)

	r.Run()
}