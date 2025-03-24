package repository

import (
	"context"
	"crypto/internal/models"
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresWalletRepo struct {
	db *sql.DB
}

func NewPostgresWalletRepo(db *sql.DB) *PostgresWalletRepo {
	return &PostgresWalletRepo{db: db}
}

func (p *PostgresWalletRepo) CreateAddress(ctx context.Context, req *models.AddressRequest) (*models.Address, error) {
	query := `INSERT INTO addresses (wallet_address, chain_name, crypto_name, tag) VALUES ($1, $2, $3, $4) RETURNING id`
	var id uint64
	err := p.db.QueryRowContext(ctx, query,
		req.WalletAddress,
		req.ChainName,
		req.CryptoName,
		req.Tag,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &models.Address{
		ID:            id,
		WalletAddress: req.WalletAddress,
		ChainName:     req.ChainName,
		CryptoName:    req.CryptoName,
		Tag:           req.Tag,
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
