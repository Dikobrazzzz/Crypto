package app

import (
	"github.com/gin-gonic/gin"
	"crypto/internal/handler"
)


func GetRouter(h *handler.Handler) *gin.Engine {
	router := gin.Default()

	router.POST("/address", h.CreateAddressHandler)
	router.GET("/address/:id", h.GetIdHandler)
	router.GET("/allwallets", h.GetAllWalletsHandler)
	router.PUT("/address/tag", h.EditTagHandler)

	return router
}