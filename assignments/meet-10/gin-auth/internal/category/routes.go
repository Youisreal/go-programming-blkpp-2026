package category

import (
	"gin-auth/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(group *gin.RouterGroup, h *Handler) {
	group.POST("/category", middleware.MockEndpoint("private.category.create"), h.Create)
	group.GET("/category", middleware.MockEndpoint("private.category.list"), h.List)
	group.GET("/category/:id", middleware.MockEndpoint("private.category.detail"), h.Detail)
	group.PUT("category/:id", middleware.MockEndpoint("private.category.update"), h.Update)
	group.DELETE("category/:id", middleware.MockEndpoint("private.category.delete"), h.Deletes)
}
