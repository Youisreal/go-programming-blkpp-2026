package product

import (
	"gin-auth/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(group *gin.RouterGroup, h *Handler) {
	group.POST("/products", middleware.MockEndpoint("private.products.create"), h.Create)
	group.GET("/products", middleware.MockEndpoint("private.products.list"), h.List)
	group.GET("/products/:id", middleware.MockEndpoint("private.products.detail"), h.Detail)
	group.PUT("products/:id", middleware.MockEndpoint("private.products.update"), h.Update)
	group.DELETE("products/:id", middleware.MockEndpoint("private.products.delete"), h.Deletes)
}
