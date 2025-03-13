package repository

import (
    "context"
	"crypto/internal/models"
	"database/sql"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

type WalletRepo struct {
    conn *pgx.Conn
}


func NewWalletProvider(conn *pgx.Conn) *WalletRepo {
    return &WalletRepo{conn: conn}
}

func (w *WalletRepo) CreateAddress(ctx context.Context, req *models.AddressRequest) (*models.Address, error) {
    
    insertSQL := `
    INSERT INTO main (wallet_address, chain_name, crypto_name, tag, balance)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id;
    `

    var newID uint64
    err := w.conn.QueryRow(ctx, insertSQL,
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

func (w *WalletRepo) GetId(ctx context.Context, id uint64) (*models.Address, error){
    
    var addr models.Address

	query :=  `
		SELECT id, wallet_address, chain_name, crypto_name, tag, balance
		FROM main 
		WHERE id = $1
	`
	err := w.conn.QueryRow(ctx, query, id).Scan(
		&addr.ID,
		&addr.WalletAddress,
		&addr.ChainName,
		&addr.CryptoName,
		&addr.Tag,
		&addr.Balance,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("No rows found", "error", err)
			return nil, err
		}
		slog.Error("QueryRow failed", "error", err)
		return nil, err
	}
    return &addr, nil

}

func (w *WalletRepo) GetAllWallets(ctx context.Context) ([]models.Address, error){

    query := `
		SELECT id, wallet_address, chain_name, crypto_name, tag, balance
		FROM main
	`

	rows, err := w.conn.Query(ctx, query) 
	if err != nil {
		slog.Error("Query failed", "error", err)
		return nil, err
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
        UPDATE main
        SET tag = $1
        WHERE id = $2
    `

    result, err := w.conn.Exec(ctx, query, req.Tag, req.ID)
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