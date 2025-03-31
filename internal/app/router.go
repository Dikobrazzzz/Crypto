package app

import (
	"crypto/internal/handler"
	"crypto/middleware"

	"github.com/gin-gonic/gin"
)

func GetRouter(h *handler.Handler) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.TracingMiddleware())
	router.POST("/address", h.CreateAddressHandler)
	router.GET("/address/:id", h.GetIDHandler)
	router.GET("/allwallets", h.GetAllWalletsHandler)
	router.PUT("/address/tag", h.EditTagHandler)

	return router
}
