package usecase

import (
	"context"
	"crypto/internal/models"
	"crypto/internal/repository"
)

type WalletUsecase struct {
	WalletCache repository.WalletProvider
}

func NewWalletProvider(cache repository.WalletProvider) *WalletUsecase {
	return &WalletUsecase{
		WalletCache: cache,
	}
}

func (w *WalletUsecase) CreateAddress(ctx context.Context, req *models.AddressRequest) (*models.Address, error) {

	return w.WalletCache.CreateAddress(ctx, req)
}

func (w *WalletUsecase) GetID(ctx context.Context, id uint64) (*models.Address, error) {

	return w.WalletCache.GetID(ctx, id)
}

func (w *WalletUsecase) GetAllWallets(ctx context.Context) ([]models.Address, error) {

	return w.WalletCache.GetAllWallets(ctx)
}

func (w *WalletUsecase) EditTag(ctx context.Context, req *models.TagUpdateRequest) error {

	return w.WalletCache.EditTag(ctx, req)
}
