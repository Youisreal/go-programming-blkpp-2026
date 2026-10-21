package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProductCRUD(t *testing.T) {
	r := setupEngine()

	// CREATE
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/products", bytes.NewBufferString(`{"name":"Kopi","price":15000}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("POST gagal: %d", w.Code)
	}

	// GET ALL
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/products", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET all gagal: %d", w.Code)
	}

	// GET BY ID
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/products/1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET by id gagal: %d", w.Code)
	}

	// UPDATE
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/products/1", bytes.NewBufferString(`{"name":"Teh","price":10000}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("PUT gagal: %d", w.Code)
	}

	// DELETE
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/products/1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("DELETE gagal: %d", w.Code)
	}
}
