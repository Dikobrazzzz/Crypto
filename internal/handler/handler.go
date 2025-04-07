package handler

import (
	"crypto/internal/apperr"
	"crypto/internal/models"
	usecase "crypto/internal/usecase"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	walletUC usecase.WalletProvider
	Instrumentation
}

type Instrumentation struct {
	tracer trace.Tracer
}

func New(walletUC usecase.WalletProvider) *Handler {
	return &Handler{
		walletUC: walletUC,
		Instrumentation: Instrumentation{
			tracer: otel.Tracer("API service"),
		},
	}
}

func (h *Handler) CreateAddressHandler(c *gin.Context) {

	ctx := c.Request.Context()
	ctx, span := h.tracer.Start(ctx, "CreateAddress")
	defer span.End()

	c.Request = c.Request.WithContext(ctx)

	var req models.AddressRequest

	if err := c.BindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.walletUC.CreateAddress(ctx, &req)
	if err != nil {
		slog.Error("Error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetIDHandler(c *gin.Context) {

	ctx := c.Request.Context()
	ctx, span := h.tracer.Start(ctx, "GetID")
	defer span.End()

	c.Request = c.Request.WithContext(ctx)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		slog.Error("Failed to parse 'id' param", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if id == 0 {
		slog.Error("Invalid 'id' parame: must be greated than 0", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	addr, err := h.walletUC.GetID(ctx, id)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		slog.Error("Error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, addr)
}

func (h *Handler) GetAllWalletsHandler(c *gin.Context) {

	ctx := c.Request.Context()
	ctx, span := h.tracer.Start(ctx, "GetAllWallets")
	defer span.End()

	c.Request = c.Request.WithContext(ctx)

	list, err := h.walletUC.GetAllWallets(ctx)
	if err != nil {
		slog.Error("Error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) EditTagHandler(c *gin.Context) {

	ctx := c.Request.Context()
	ctx, span := h.tracer.Start(ctx, "EditTag")
	defer span.End()

	c.Request = c.Request.WithContext(ctx)

	var req models.TagUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON for tag update", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.walletUC.EditTag(ctx, &req); err != nil {
		slog.Error("Failed to update tag", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": req.ID, "tag": req.Tag})
}
