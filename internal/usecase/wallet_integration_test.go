package usecase_test

import (
	"context"
	"crypto/internal/models"
	"crypto/internal/repository"
	"crypto/internal/usecase"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var db *sql.DB

func TestMain(m *testing.M) {
	dsn := "postgres://postgres:postgres@localhost:5432/test_db"

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("failed to ping DB: %v", err)
	}

	code := m.Run()

	_ = db.Close()
	os.Exit(code)
}

func TestWalletUsecase_CreateAddress(t *testing.T) {
	repo := repository.NewPostgresWalletRepo(db)
	walletUsecase := usecase.NewWalletProvider(repo)

	ctx := context.Background()
	req := &models.AddressRequest{
		WalletAddress: "0xABC",
		ChainName:     "Ethereum",
		CryptoName:    "ETH",
		Tag:           "myTag",
	}

	address, err := walletUsecase.CreateAddress(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, address)
	require.NotZero(t, address.ID)
	require.Equal(t, req.WalletAddress, address.WalletAddress)
	require.Equal(t, req.ChainName, address.ChainName)
	require.Equal(t, req.CryptoName, address.CryptoName)
	require.Equal(t, req.Tag, address.Tag)
}
