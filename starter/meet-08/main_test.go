package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"meet08/internal/product"
)

// PostgreSQL asli, schema sementara. Gunakan koneksi lokal/session, bukan transaction pooler.
func TestCRUD(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL untuk integration check PostgreSQL")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatal("koneksi database test gagal")
	}
	schema := fmt.Sprintf("northwind_test_%d", time.Now().UnixNano())
	if _, err := db.ExecContext(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatal(err)
	}
	defer func() {
		cleanup, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()
		if _, err := db.ExecContext(cleanup, "DROP SCHEMA "+schema+" CASCADE"); err != nil {
			t.Errorf("cleanup schema test: %v", err)
		}
	}()
	if _, err := db.ExecContext(ctx, "SET search_path TO "+schema); err != nil {
		t.Fatal(err)
	}
	sample, err := os.ReadFile("northwind.sql")
	if err != nil {
		t.Fatal(err)
	}
	a := product.Handler{Service: product.Service{Repository: product.Repository{DB: db}}}
	if _, err := db.ExecContext(ctx, string(sample)); err != nil {
		t.Fatalf("import: %v", err)
	}
	request := func(method, path, body string, want int) []byte {
		t.Helper()
		r := httptest.NewRequest(method, path, bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		a.ServeHTTP(w, r)
		if w.Code != want || !json.Valid(w.Body.Bytes()) {
			t.Fatalf("%s %s: status=%d body=%s, want=%d", method, path, w.Code, w.Body, want)
		}
		return w.Body.Bytes()
	}
	request("GET", "/products", "", 200)
	request("GET", "/categories", "", 200)
	request("GET", "/suppliers", "", 200)
	request("GET", "/products/1", "", 200)
	request("GET", "/products/99999", "", 404)
	request("GET", "/products/1%20OR%201=1", "", 400)
	request("POST", "/products/1", "{}", 405)
	request("PATCH", "/products", "{}", 405)
	request("POST", "/products", "{", 400)
	request("POST", "/products", `{"product_name":"Tea","supplier_id":1,"category_id":1,"typo":1}`, 400)
	request("POST", "/products", `{"product_name":"Tea","supplier_id":1,"category_id":1} {}`, 400)
	request("POST", "/products", `{"product_name":"Tea","supplier_id":1,"category_id":1,"unit_price":-1}`, 400)
	request("POST", "/products", `{"product_name":"Tea","supplier_id":1,"category_id":1,"unit_price":10000000000}`, 400)
	request("POST", "/products", `{"product_name":"Tea","supplier_id":99999,"category_id":1}`, 409)

	p := product.Input{ProductName: "Tea'; DROP TABLE products;--\n.test", SupplierID: 1, CategoryID: 1, UnitPrice: 12.5, UnitsInStock: 3}
	body, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	created := request("POST", "/products", string(body), 201)
	var row struct {
		ProductID int    `json:"product_id"`
		Name      string `json:"product_name"`
	}
	if err := json.Unmarshal(created, &row); err != nil || row.ProductID <= 8 || row.Name != p.ProductName {
		t.Fatalf("create: %s, err=%v", created, err)
	}
	path := fmt.Sprintf("/products/%d", row.ProductID)
	request("GET", path, "", 200)
	p.ProductName = "Teh diperbarui"
	body, err = json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	updated := request("PUT", path, string(body), 200)
	if err := json.Unmarshal(updated, &row); err != nil || row.Name != p.ProductName {
		t.Fatalf("update: %s, err=%v", updated, err)
	}
	request("PUT", "/products/99999", string(body), 404)
	request("DELETE", path, "", 200)
	request("DELETE", path, "", 404)
	request("GET", path, "", 404)
	request("GET", "/products/1", "", 200)

	// Import ulang tidak menimpa data yang sudah diedit.
	if _, err := db.ExecContext(ctx, "UPDATE products SET unit_price = 99 WHERE product_id = 1"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, string(sample)); err != nil {
		t.Fatalf("reimport: %v", err)
	}
	var saved struct {
		Price float64 `json:"unit_price"`
	}
	if err := json.Unmarshal(request("GET", "/products/1", "", 200), &saved); err != nil || saved.Price != 99 {
		t.Fatalf("reimport menimpa data: price=%v err=%v", saved.Price, err)
	}
}
