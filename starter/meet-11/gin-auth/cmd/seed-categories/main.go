package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gin-auth/internal/category"
	"gin-auth/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatal(fmt.Errorf("load .env: %w", err))
	}

	db, err := database.Open(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if conn, cerr := db.DB(); cerr == nil {
			_ = conn.Close()
		}
	}()

	if err := db.AutoMigrate(&category.Category{}); err != nil {
		log.Fatal(fmt.Errorf("migrasi category gagal: %w", err))
	}

	for _, item := range category.DefaultCategories {
		result := db.Where(category.Category{Name: item.Name}).FirstOrCreate(&item)
		if result.Error != nil {
			log.Fatal(fmt.Errorf("seed kategori %q gagal: %w", item.Name, result.Error))
		}
	}

	log.Println("seed categories berhasil")
}
