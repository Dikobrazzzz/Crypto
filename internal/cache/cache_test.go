package cache_test

import (
	"context"
	"crypto/internal/cache"
	"crypto/internal/models"
	"crypto/internal/repository"
	"log/slog"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCacheDecorator_CreateAddress(t *testing.T) {
	tests := []struct {
		name               string
		req                *models.AddressRequest
		createAddressResp  *models.Address
		createAddressError error
		expectedWallet     *models.Address
		expectError        error
	}{
		{
			name: "success call",
			req: &models.AddressRequest{
				WalletAddress: "0x123",
				ChainName:     "Ethereum",
				CryptoName:    "ETH",
				Tag:           "test",
			},
			createAddressResp: &models.Address{
				WalletAddress: "0x123",
				ChainName:     "Ethereum",
				CryptoName:    "ETH",
				Tag:           "test",
			},
			createAddressError: nil,
			expectedWallet: &models.Address{
				WalletAddress: "0x123",
				ChainName:     "Ethereum",
				CryptoName:    "ETH",
				Tag:           "test",
			},
			expectError: nil,
		},
		{
			name: "repo error",
			req: &models.AddressRequest{
				WalletAddress: "0xabc",
				ChainName:     "Ethereum",
				CryptoName:    "ETH",
				Tag:           "error",
			},
			createAddressResp:  nil,
			createAddressError: errors.New("db error"),
			expectedWallet:     nil,
			expectError:        errors.New("The operation cannot be performed"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockWalletRepo := repository.NewMockWalletProvider(ctrl)
			mockWalletRepo.
				EXPECT().
				CreateAddress(gomock.Any(), tc.req).
				Return(tc.createAddressResp, tc.createAddressError)

			cacheDecorator := &cache.CacheDecorator{
				WalletRepo: mockWalletRepo,
				Wallets:    make(map[uint64]cache.WrapWallet),
			}

			wallet, err := cacheDecorator.CreateAddress(context.Background(), tc.req)
			if tc.expectError != nil {
				slog.Error("Error with create address")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedWallet, wallet)
			}
		})
	}
}
