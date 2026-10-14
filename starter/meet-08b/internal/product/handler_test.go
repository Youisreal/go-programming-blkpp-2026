package product

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutingDeterministic(t *testing.T) {
	h := Handler{}
	routes := h.Routes()

	// 404 on unknown route
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	routes.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("want 404 for unknown route, got %d", rec.Code)
	}

	// 400 on bad ID for item route
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/products/abc", nil)
	routes.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400 for bad ID, got %d", rec.Code)
	}

	// 405 on POST to /products/1
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/products/1", bytes.NewBufferString("{}"))
	routes.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("want 405 for POST item, got %d", rec.Code)
	}
	if rec.Header().Get("Allow") != "GET, PUT, DELETE" {
		t.Fatalf("want Allow: GET, PUT, DELETE, got %s", rec.Header().Get("Allow"))
	}

	// 405 on DELETE to /products
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/products", nil)
	routes.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("want 405 for DELETE collection, got %d", rec.Code)
	}
	if rec.Header().Get("Allow") != "GET, POST" {
		t.Fatalf("want Allow: GET, POST, got %s", rec.Header().Get("Allow"))
	}

	// 405 on POST to /categories
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/categories", nil)
	routes.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("want 405 for POST /categories, got %d", rec.Code)
	}
	if rec.Header().Get("Allow") != "GET" {
		t.Fatalf("want Allow: GET, got %s", rec.Header().Get("Allow"))
	}
}
