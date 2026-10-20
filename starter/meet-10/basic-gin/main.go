package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type User struct {
	Name   string `form:"name" json:"name" binding:"required"`
	Age    int    `form:"age" json:"age" binding:"required"`
	Gender string `form:"gender" json:"gender" binding:"required"`
}

func formPostHandler(c *gin.Context) {
	message := c.PostForm("message")
	nick := c.DefaultPostForm("nick", "anonymous")
	gender := c.DefaultPostForm("gender", "")

	c.JSON(200, gin.H{
		"status":  "posted",
		"message": message,
		"nick":    nick,
		"gender":  gender,
	})
}

func main() {
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	router.POST("/users", func(c *gin.Context) {
		name := c.PostForm("name")
		c.JSON(http.StatusCreated, gin.H{"user": name})
	})

	router.POST("/form_post", formPostHandler)

	router.POST("/post", func(c *gin.Context) {
		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		message := c.PostForm("message")

		c.JSON(http.StatusOK, gin.H{
			"id":      id,
			"page":    page,
			"name":    name,
			"message": message,
		})
	})

	router.POST("/upload", func(c *gin.Context) {
		// single file
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Println(file.Filename)

		// Upload the file to specific dst.
		dst := filepath.Join("./files/", filepath.Base(file.Filename))
		c.SaveUploadedFile(file, dst)

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	router.Run(":8080")
}
