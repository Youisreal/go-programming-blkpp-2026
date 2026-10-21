package product

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *Handler) {
	group := r.Group("/api/products")
	{
		group.GET("", h.GetAll)
		group.GET("/:id", h.GetByID)
		group.POST("", h.Create)
		group.PUT("/:id", h.Update)
		group.DELETE("/:id", h.Delete)
	}
}
