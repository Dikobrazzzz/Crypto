package repository

import (
	"context"
	"crypto/internal/apperr"
	"crypto/internal/models"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type WalletRepo struct {
	pool *pgxpool.Pool
}

func NewWalletProvider(pool *pgxpool.Pool) *WalletRepo {
	return &WalletRepo{pool: pool}
}

func (w *WalletRepo) CreateAddress(ctx context.Context, req *models.AddressRequest) (*models.Address, error) {

	insertSQL := `
    INSERT INTO wallet (wallet_address, chain_name, crypto_name, tag, balance)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id;
    `

	var newID uint64
	err := w.pool.QueryRow(ctx, insertSQL,
		req.WalletAddress, req.ChainName, req.CryptoName, req.Tag, 0,
	).Scan(&newID)
	if err != nil {
		slog.Error("Failed to insert into table", "error", err)
		return nil, err
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

func (w *WalletRepo) GetID(ctx context.Context, id uint64) (*models.Address, error) {

	var addr models.Address

	query := `
		SELECT id, wallet_address, chain_name, crypto_name, tag, balance
		FROM wallet 
		WHERE id = $1
	`
	err := w.pool.QueryRow(ctx, query, id).Scan(
		&addr.ID,
		&addr.WalletAddress,
		&addr.ChainName,
		&addr.CryptoName,
		&addr.Tag,
		&addr.Balance,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			errors.Is(err, pgx.ErrNoRows)
			return nil, apperr.ErrNotFound
		}
		slog.Error("QueryRow failed", "error", err)
		return nil, err
	}
	return &addr, nil

}

func (w *WalletRepo) GetAllWallets(ctx context.Context) ([]models.Address, error) {

	query := `
		SELECT id, wallet_address, chain_name, crypto_name, tag, balance
		FROM wallet
	`

	rows, err := w.pool.Query(ctx, query)
	if err != nil {
		slog.Error("Query failed", "error", err)
		return nil, errors.Wrap(err, "GetAllWallets Query")
	}
	defer rows.Close()

	var list []models.Address
	for rows.Next() {
		var addr models.Address
		if err := rows.Scan(
			&addr.ID,
			&addr.WalletAddress,
			&addr.ChainName,
			&addr.CryptoName,
			&addr.Tag,
			&addr.Balance,
		); err != nil {
			slog.Error("Row scan failed", "error", err)
			return nil, err
		}
		list = append(list, addr)

	}
	if err := rows.Err(); err != nil {
		slog.Error("Rows error", "error", err)
		return nil, err
	}

	return list, nil
}

func (w *WalletRepo) EditTag(ctx context.Context, req *models.TagUpdateRequest) error {

	query := `
        UPDATE wallet
        SET tag = $1
        WHERE id = $2
    `

	result, err := w.pool.Exec(ctx, query, req.Tag, req.ID)
	if err != nil {
		slog.Error("Update failed", "error", err)
		return err
	}

	if result.RowsAffected() == 0 {
		slog.Error("No rows were affected by the update")
		return err
	}
	return nil
}
