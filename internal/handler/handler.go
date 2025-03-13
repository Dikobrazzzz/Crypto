package handler

import (
	"net/http"
	"strconv"
	"log/slog" 
	"context"
	"crypto/internal/models"
	usecase "crypto/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	walletUC usecase.WalletProvider
	ctx      context.Context
}


func New(walletUC usecase.WalletProvider) *Handler {
	return &Handler{
		walletUC: walletUC,
		ctx:      context.Background(),
	}
}

func (h *Handler) CreateAddressHandler(c *gin.Context) {
	
	var req models.AddressRequest

	if err := c.BindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.walletUC.CreateAddress(c.Request.Context(), &req) 
	if err != nil {
		slog.Error("Error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

    c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetIdHandler(c *gin.Context) {
	
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		slog.Error("Failed to parse 'id' param", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	addr, err := h.walletUC.GetId(c.Request.Context(), id)
	if err != nil {
		slog.Error("Error", "error", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	c.JSON(http.StatusOK, addr)
}


func (h *Handler) GetAllWalletsHandler(c *gin.Context) {

	list, err := h.walletUC.GetAllWallets(c.Request.Context())
	if err != nil {
		slog.Error("Error", "error", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
	c.JSON(http.StatusOK, list)
}


func (h *Handler) EditTagHandler(c *gin.Context) {

	var req models.TagUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON for tag update", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.walletUC.EditTag(c.Request.Context(), &req); err != nil {
		slog.Error("Failed to update tag", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": req.ID, "tag": req.Tag})
}
