package product

import (
	"net/http"
	"strconv"
	"strings"

	"gin-auth/internal/auth"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required,max=100"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || strings.TrimSpace(input.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama product wajib diisi, maksimal 100 karakter"})
		return
	}
	user := c.MustGet("user").(auth.User)
	item, err := h.service.Create(c.Request.Context(), input.Name, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan product"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": item})
}

func (h *Handler) List(c *gin.Context) {
	user := c.MustGet("user").(auth.User)
	items, err := h.service.List(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func parseID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus berupa angka"})
		return 0, false
	}
	return id, true
}

func (h *Handler) Detail(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	user := c.MustGet("user").(auth.User)
	item, err := h.service.Detail(c.Request.Context(), user.ID, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Product tidak ditemukan"})
		return
	}
	// WAJIB: Kirim hasil produk yang ditemukan ke client
	c.JSON(http.StatusOK, item)

}

func (h *Handler) Update(c *gin.Context) {
	user := c.MustGet("user").(auth.User)
	id, ok := parseID(c)
	if !ok {
		return
	}
	_, err := h.service.Detail(c.Request.Context(), user.ID, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca product"})
		return
	}
	var input UpdateProductRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedItem, err := h.service.Update(c.Request.Context(), user.ID, uint(id), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Product berhasil diupdate",
		"data":    updatedItem,
	})
}

func (h *Handler) Deletes(c *gin.Context) {
	user := c.MustGet("user").(auth.User)

	id, ok := parseID(c)
	if !ok {
		return
	}

	err := h.service.Delete(c.Request.Context(), user.ID, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product berhasil dihapus",
	})
}
