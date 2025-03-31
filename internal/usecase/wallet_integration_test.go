package usecase

import (
	"context"
	"crypto/config"
	"crypto/internal/models"
	"crypto/internal/repository"
	"crypto/internal/storage"
	"log/slog"
	"testing"

	_ "github.com/lib/pq"

	"github.com/stretchr/testify/require"
)

func TestWalletUsecase_CreateAddress(t *testing.T) {

	ctx := context.Background()
	pool, err := storage.GetConnection(ctx, config.AppConfig.DatabaseURL)
	if err != nil {
		slog.Error("Error getting connection", "error", err)
		return
	}
	defer pool.Close()

	repo := repository.NewPostgresWalletRepo(pool)
	walletUsecase := NewWalletProvider(repo)

	testCases := []struct {
		name          string
		WalletAddress string
		ChainName     string
		CryptoName    string
		Tag           string
		expectError   bool
	}{
		{
			name:          "success call",
			WalletAddress: "9999999999999999999999999FFFFFFF",
			ChainName:     "Ethereum",
			CryptoName:    "ETH",
			Tag:           "myTag",
			expectError:   false,
		},
		{
			name:          "success call",
			WalletAddress: "",
			ChainName:     "Ethereum",
			CryptoName:    "ETH",
			Tag:           "myTag",
			expectError:   false,
		},
		{
			name:          "success call",
			WalletAddress: "0xABCASDFIO",
			ChainName:     "X",
			CryptoName:    "ETH",
			Tag:           "5",
			expectError:   false,
		},
		{
			name:          "success call",
			WalletAddress: "0xABC",
			ChainName:     "Ethereum",
			CryptoName:    "ETH",
			Tag:           "",
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		req := &models.AddressRequest{
			WalletAddress: tc.WalletAddress,
			ChainName:     tc.ChainName,
			CryptoName:    tc.CryptoName,
			Tag:           tc.Tag,
		}
		address, err := walletUsecase.CreateAddress(ctx, req)
		if tc.expectError {
			require.Error(t, err)
			return
		}
		require.NoError(t, err)
		require.NotNil(t, address)
		// require.NotZero(t, address.ID)
		require.Equal(t, tc.WalletAddress, address.WalletAddress)
		require.Equal(t, tc.ChainName, address.ChainName)
		require.Equal(t, tc.CryptoName, address.CryptoName)
		require.Equal(t, tc.Tag, address.Tag)
	}
}
