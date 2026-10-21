package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		unauthorized := func() {
			c.Header("WWW-Authenticate", "Bearer")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau kedaluwarsa"})
		}
		parts := strings.Fields(c.GetHeader("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			unauthorized()
			return
		}
		user, err := h.service.Authenticate(c.Request.Context(), parts[1])
		if errors.Is(err, ErrInvalidToken) {
			unauthorized()
			return
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca pengguna"})
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
