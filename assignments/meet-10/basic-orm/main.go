package main

import (
	"log"
	"net/http"
	"os"

	"meet10/internal/product"

	"github.com/gin-gonic/gin"
)

func setupEngine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"status": "ok"}})
	})

	repo := product.NewRepository()
	service := product.NewService(repo)
	handler := product.NewHandler(service)

	product.RegisterRoutes(r, handler)
	return r
}

func main() {
	r := setupEngine()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
