package interfaces
import (
    "context"
	"crypto/internal/models"
	"crypto/internal/repository"
)

type WalletUsecase struct {
    WalletRepo repository.WalletProvider
}


func NewWalletProvider(repo repository.WalletProvider) *WalletUsecase {
    return &WalletUsecase{
        WalletRepo: repo,
    }
}


func (w *WalletUsecase) CreateAddress(ctx context.Context, req *models.AddressRequest) (*models.Address, error) { 

    return w.WalletRepo.CreateAddress(ctx, req)
}

func (w *WalletUsecase) GetId(ctx context.Context, id uint64) (*models.Address, error) {

    return w.WalletRepo.GetId(ctx,id)
}

func (w *WalletUsecase) GetAllWallets(ctx context.Context) ([]models.Address, error) {

    return w.WalletRepo.GetAllWallets(ctx)
}

func (w *WalletUsecase) EditTag(ctx context.Context, req *models.TagUpdateRequest) error {

    return w.WalletRepo.EditTag(ctx, req)
}
