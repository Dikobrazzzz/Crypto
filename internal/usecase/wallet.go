package usecase

import (
	"context"
	"crypto/internal/models"
	"crypto/internal/repository"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type WalletUsecase struct {
	WalletRepo repository.WalletProvider
	Instrumentation
}

type Instrumentation struct {
	tracer trace.Tracer
}

func NewWalletProvider(repo repository.WalletProvider) *WalletUsecase {
	return &WalletUsecase{
		WalletRepo: repo,
		Instrumentation: Instrumentation{
			tracer: otel.Tracer("API service"),
		},
	}
}

func (w *WalletUsecase) CreateAddress(ctx context.Context, req *models.AddressRequest) (*models.Address, error) {
	ctx, span := w.tracer.Start(ctx, "CreateAddress")
	defer span.End()
	return w.WalletRepo.CreateAddress(ctx, req)
}

func (w *WalletUsecase) GetID(ctx context.Context, id uint64) (*models.Address, error) {

	ctx, span := w.tracer.Start(ctx, "GetID")
	defer span.End()
	return w.WalletRepo.GetID(ctx, id)
}

func (w *WalletUsecase) GetAllWallets(ctx context.Context) ([]models.Address, error) {

	ctx, span := w.tracer.Start(ctx, "GetAllWallets")
	defer span.End()
	return w.WalletRepo.GetAllWallets(ctx)
}

func (w *WalletUsecase) EditTag(ctx context.Context, req *models.TagUpdateRequest) error {

	ctx, span := w.tracer.Start(ctx, "EditTag")
	defer span.End()
	return w.WalletRepo.EditTag(ctx, req)
}
