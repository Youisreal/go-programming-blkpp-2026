package product

import (
	"net/http"
	"strings"

	"gin-auth/internal/auth"
	"gin-auth/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(group *gin.RouterGroup, db *gorm.DB) {
	group.POST("/products", middleware.MockEndpoint("private.products.create"), func(c *gin.Context) {
		var input struct {
			Name string `json:"name" binding:"required,max=100"`
		}
		if err := c.ShouldBindJSON(&input); err != nil || strings.TrimSpace(input.Name) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nama product wajib diisi, maksimal 100 karakter"})
			return
		}
		user := c.MustGet("user").(auth.User)
		item := Product{Name: strings.TrimSpace(input.Name), UserID: user.ID}
		if err := db.WithContext(c.Request.Context()).Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan product"})
			return
		}
		item.User = user
		c.JSON(http.StatusCreated, gin.H{"data": item})
	})
	group.GET("/products", middleware.MockEndpoint("private.products.list"), func(c *gin.Context) {
		user := c.MustGet("user").(auth.User)
		items := []Product{}
		// ponytail: daftar maksimal 100 product terbaru; tambah pagination saat data melebihi batas.
		if err := db.WithContext(c.Request.Context()).Where("user_id = ?", user.ID).
			Preload("User").Order("id DESC").Limit(100).Find(&items).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca product"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": items})
	})
}
