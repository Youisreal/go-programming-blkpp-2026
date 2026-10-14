package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"meet08/internal/config"
	"meet08/internal/database"
	"meet08/internal/product"
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
	service := product.Service{Repository: repository}
	handler := product.Handler{Service: service}
	server := &http.Server{
		Addr: cfg.Addr, Handler: handler,
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
