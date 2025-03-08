package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log/slog" 
	"context"
)

type AddressRequest struct {
	WalletAddress string `json:"walletAddress"`
	ChainName     string `json:"chainName"`
	CryptoName    string `json:"cryptoName"`
	Tag           string `json:"tag"`
}

type Address struct {
	ID            uint64 `json:"id"`
	WalletAddress string `json:"wallet_address"`
	ChainName     string `json:"chain_name"`
	CryptoName    string `json:"crypto_name"`
	Tag           string `json:"tag"`
	Balance       int64  `json:"balance"`
}

type TagUpdateRequest struct {
	ID  uint64 `json:"id"`
	Tag string `json:"tag"`
}

type Handle struct {
	conn *pgx.Conn
	cb   context.Context
}

func New(conn *pgx.Conn, cb context.Context) *Handle {
	return &Handle{conn: conn, cb: cb}
}

func (h *Handle) CreateAddressHandler(c *gin.Context) {
	var req AddressRequest

	if err := c.BindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	insertSQL := `
		INSERT INTO main (wallet_address, chain_name, crypto_name, tag, balance)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	var newID uint64
	err := h.conn.QueryRow(h.cb, insertSQL,
		req.WalletAddress, req.ChainName, req.CryptoName, req.Tag, 0,
	).Scan(&newID)
	if err != nil {
		slog.Error("Failed to insert into table", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Address{
		ID:            newID,
		WalletAddress: req.WalletAddress,
		ChainName:     req.ChainName,
		CryptoName:    req.CryptoName,
		Tag:           req.Tag,
		Balance:       0,
	})
}

func (h *Handle) GetIdHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		slog.Error("Failed to parse 'id' param", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var addr Address

	err = h.conn.QueryRow(h.cb, `
		SELECT id, wallet_address, chain_name, crypto_name, tag, balance
		FROM main 
		WHERE id = $1
	`, id).Scan(
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
			c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		}
		slog.Error("QueryRow failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, addr)
}


func (h *Handle) GetAllWalletsHandler(c *gin.Context) {

	rows, err := h.conn.Query(h.cb, `
		SELECT id, wallet_address, chain_name, crypto_name, tag, balance
		FROM main
	`)
	if err != nil {
		slog.Error("Query failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var list []Address
	for rows.Next() {
		var addr Address
		if err := rows.Scan(
			&addr.ID,
			&addr.WalletAddress,
			&addr.ChainName,
			&addr.CryptoName,
			&addr.Tag,
			&addr.Balance,
		); err != nil {
			slog.Error("Row scan failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		list = append(list, addr)
	}

	c.JSON(http.StatusOK, list)
}


func (h *Handle) EditTagHandler(c *gin.Context) {
	var req TagUpdateRequest

	if err := c.BindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON for tag update", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.conn.Exec(h.cb, `
		UPDATE main
		SET tag = $1
		WHERE id = $2
	`, req.Tag, req.ID)
	if err != nil {
		slog.Error("Update failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.RowsAffected() == 0 {
		slog.Error("No rows were affected by the update", "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": req.ID, "tag": req.Tag})
}
