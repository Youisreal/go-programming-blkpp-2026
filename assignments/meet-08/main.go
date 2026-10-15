package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"meet08b/internal/analysiscust"
	"meet08b/internal/analysispro"
	"meet08b/internal/category"
	"meet08b/internal/config"
	"meet08b/internal/database"
	"meet08b/internal/product"
	"meet08b/internal/summary"
	"meet08b/internal/supplier"
)

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	db, err := database.Open(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	repository := product.Repository{DB: db}
	if _, err := repository.Get(context.Background(), 1); err != nil {
		return fmt.Errorf("database belum siap; import northwind.sql dahulu: %w", err)
	}

	productHandler := product.Handler{Service: product.Service{Repository: repository}}
	categoryHandler := category.Handler{Service: category.Service{Repository: category.Repository{DB: db}}}
	supplierHandler := supplier.Handler{Service: supplier.Service{Repository: supplier.Repository{DB: db}}}
	summaryHandler := summary.Handler{Service: summary.Service{Repository: summary.Repository{DB: db}}}
	analysiscustHandler := analysiscust.Handler{Service: analysiscust.Service{Repository: analysiscust.Repository{DB: db}}}
	analysisproHandler := analysispro.Handler{Service: analysispro.Service{Repository: analysispro.Repository{DB: db}}}

	mux := http.NewServeMux()
	mux.Handle("/products", productHandler)
	mux.Handle("/products/", productHandler)
	mux.Handle("/categories", categoryHandler)
	mux.Handle("/suppliers", supplierHandler)
	mux.Handle("/summary", summaryHandler)
	mux.Handle("/analysiscust", analysiscustHandler)
	mux.Handle("/analysispro", analysisproHandler)

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second,
	}
	log.Printf("Northwind API: http://%s", cfg.Addr)
	return server.ListenAndServe()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
