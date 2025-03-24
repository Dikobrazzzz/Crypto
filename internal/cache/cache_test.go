package cache_test

import (
	"context"
	"crypto/internal/cache"
	"crypto/internal/models"
	"crypto/internal/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCacheDecorator_CreateAddress(t *testing.T) {
	tests := []struct {
		name           string
		req            *models.AddressRequest
		repoWallet     *models.Address
		repoErr        error
		expectedWallet *models.Address
		expectError    bool
	}{
		{
			name: "success call",
			req: &models.AddressRequest{
				WalletAddress: "0x123",
				ChainName:     "Ethereum",
				CryptoName:    "ETH",
				Tag:           "test",
			},
			repoWallet: &models.Address{
				WalletAddress: "0x123",
				ChainName:     "Ethereum",
				CryptoName:    "ETH",
				Tag:           "test",
			},
			repoErr: nil,
			expectedWallet: &models.Address{
				WalletAddress: "0x123",
				ChainName:     "Ethereum",
				CryptoName:    "ETH",
				Tag:           "test",
			},
			expectError: false,
		},
		{
			name: "repo error",
			req: &models.AddressRequest{
				WalletAddress: "0xabc",
				ChainName:     "Ethereum",
				CryptoName:    "ETH",
				Tag:           "error",
			},
			repoWallet:     nil,
			repoErr:        errors.New("db error"),
			expectedWallet: nil,
			expectError:    true,
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
				Return(tc.repoWallet, tc.repoErr)

			c := &cache.CacheDecorator{
				WalletRepo: mockWalletRepo,
				Wallets:    make(map[uint64]cache.WrapWallet),
			}

			wallet, err := c.CreateAddress(context.Background(), tc.req)
			if tc.expectError {
				require.Error(t, err)
				require.Nil(t, wallet)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedWallet, wallet)
			}
		})
	}
}
