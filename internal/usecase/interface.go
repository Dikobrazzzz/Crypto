package interfaces

import (
	"context"
	"crypto/internal/models"
)

type WalletProvider interface {
	CreateAddress(ctx context.Context, req *models.AddressRequest) (*models.Address, error)
	GetId(ctx context.Context, id uint64) (*models.Address, error)
	GetAllWallets(ctx context.Context) ([]models.Address, error)
	EditTag(ctx context.Context, req *models.TagUpdateRequest) error
}
