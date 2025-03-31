package repository

import (
	"context"
	"crypto/internal/models"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type PostgresWalletRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresWalletRepo(pool *pgxpool.Pool) *PostgresWalletRepo {
	return &PostgresWalletRepo{pool: pool}
}

func (p *PostgresWalletRepo) CreateAddress(ctx context.Context, req *models.AddressRequest) (*models.Address, error) {
	insertSQL := `
	INSERT INTO wallet (wallet_address, chain_name, crypto_name, tag, balance)
	VALUES ($1, $2, $3, $4, $5)
	RETURN id
	`
	var newID uint64
	err := p.pool.QueryRow(ctx, insertSQL, req.WalletAddress, req.ChainName, req.CryptoName, req.Tag, 0).Scan(&newID)
	if err != nil {
		slog.Error("Failed to insert into table")
	}

	return &models.Address{
		ID:            newID,
		WalletAddress: req.WalletAddress,
		ChainName:     req.ChainName,
		CryptoName:    req.CryptoName,
		Tag:           req.Tag,
		Balance:       0,
	}, nil
}

func (p *PostgresWalletRepo) GetID(ctx context.Context, id uint64) (*models.Address, error) {
	return nil, nil
}

func (p *PostgresWalletRepo) GetAllWallets(ctx context.Context) ([]models.Address, error) {
	return nil, nil
}

func (p *PostgresWalletRepo) EditTag(ctx context.Context, req *models.TagUpdateRequest) error {
	return nil
}
