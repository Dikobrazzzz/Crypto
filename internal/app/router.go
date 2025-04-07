package app

import (
	"crypto/internal/handler"
	"crypto/middleware"

	"github.com/gin-gonic/gin"
)

func GetRouter(h *handler.Handler) *gin.Engine {
	router := gin.Default()
	address := router.Group("/address")
	router.Use(middleware.TracingMiddleware())
	address.POST("/", h.CreateAddressHandler)
	address.GET("/:id", h.GetIDHandler)
	address.GET("/allwallets", h.GetAllWalletsHandler)
	address.PUT("/tag", h.EditTagHandler)

	return router
}
