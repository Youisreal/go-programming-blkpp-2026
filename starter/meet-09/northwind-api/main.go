package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"northwind-api/internal/category"
	"northwind-api/internal/customers"
	"northwind-api/internal/database"
	"northwind-api/internal/httpjson"
	"northwind-api/internal/order"
	"northwind-api/internal/product"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load .env: %w", err)
	}
	db, err := database.ConnectDatabase(os.Getenv("DATABASE_URL"))
	if err != nil {
		return databaseStartupError(err)
	}
	defer db.Close()
	mux := newMux(db)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{Addr: ":" + port, Handler: mux, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second}
	log.Printf("HTTP server listening on :%s", port)
	return server.ListenAndServe()
}

func databaseStartupError(err error) error {
	// Error parsing DSN bisa berisi password; tampilkan hanya error jaringan atau SQLSTATE.
	var networkError *net.OpError
	if errors.As(err, &networkError) {
		return fmt.Errorf("koneksi PostgreSQL gagal: %s; periksa jaringan, host, dan dukungan IPv4/IPv6", networkError.Err)
	}
	var postgresError *pq.Error
	if errors.As(err, &postgresError) {
		return fmt.Errorf("koneksi PostgreSQL gagal: SQLSTATE %s (%s); periksa user, password, database, dan SSL", postgresError.Code, postgresError.Code.Name())
	}
	return errors.New("koneksi PostgreSQL gagal; periksa format DATABASE_URL, jaringan, SSL, dan kredensial")
}

func newMux(db *sql.DB) *http.ServeMux {
	productRepo := product.NewRepository(db)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)
	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService)
	orderRepo := order.NewRepository(db)
	orderService := order.NewService(orderRepo)
	orderHandler := order.NewHandler(orderService)
	customerRepo := customers.NewRepository(db)
	customerService := customers.NewService(customerRepo)
	customerHandler := customers.NewHandler(customerService)

	mux := http.NewServeMux()
	product.RegisterRoutes(mux, productHandler)
	category.RegisterRoutes(mux, categoryHandler)
	order.RegisterRoutes(mux, orderHandler)
	customers.RegisterRoutes(mux, customerHandler)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		httpjson.Data(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", "GET, HEAD")
		httpjson.Error(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpjson.Error(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found")
	})
	return mux
}
