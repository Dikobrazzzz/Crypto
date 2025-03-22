package app

import (
	"crypto/internal/handler"

	"github.com/gin-gonic/gin"
)

func GetRouter(h *handler.Handler) *gin.Engine {
	router := gin.Default()

	router.POST("/address", h.CreateAddressHandler)
	router.GET("/address/:id", h.GetIDHandler)
	router.GET("/allwallets", h.GetAllWalletsHandler)
	router.PUT("/address/tag", h.EditTagHandler)

	return router
}
