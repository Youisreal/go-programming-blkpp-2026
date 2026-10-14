package product

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct{ Service Service }

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("response: %v", err)
	}
}

func fail(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func respond(w http.ResponseWriter, status int, rows json.RawMessage, err error) {
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidInput), errors.Is(err, ErrConstraint):
			fail(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrNotFound):
			fail(w, http.StatusNotFound, err.Error())
		case errors.Is(err, ErrForeignKey), errors.Is(err, ErrConflict):
			fail(w, http.StatusConflict, err.Error())
		default:
			log.Printf("database: %v", err)
			fail(w, http.StatusInternalServerError, "query database gagal")
		}
		return
	}
	writeJSON(w, status, rows)
}

func decodeBody(w http.ResponseWriter, r *http.Request) (Input, bool) {
	var p Input
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&p); err != nil {
		fail(w, http.StatusBadRequest, "body harus JSON produk valid, maksimal 1 MiB")
		return p, false
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		fail(w, http.StatusBadRequest, "body harus berisi satu objek JSON")
		return p, false
	}
	return p, true
}

func parseID(w http.ResponseWriter, path, prefix string) (int, bool) {
	idStr := strings.TrimPrefix(path, prefix)
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		fail(w, http.StatusBadRequest, "ID produk harus bilangan bulat positif")
		return 0, false
	}
	return id, true
}

func (h Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/categories", h.Categories)
	mux.HandleFunc("/suppliers", h.Suppliers)
	mux.HandleFunc("/products", h.Products)
	mux.HandleFunc("/products/", h.ProductItem)
	return mux
}

func (h Handler) Categories(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/categories" {
		fail(w, http.StatusNotFound, "endpoint tidak ditemukan")
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		fail(w, http.StatusMethodNotAllowed, "method tidak didukung")
		return
	}
	if h.Service.Repository.DB == nil {
		fail(w, http.StatusInternalServerError, "database belum siap")
		return
	}
	rows, err := h.Service.Categories(r.Context())
	respond(w, http.StatusOK, rows, err)
}

func (h Handler) Suppliers(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/suppliers" {
		fail(w, http.StatusNotFound, "endpoint tidak ditemukan")
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		fail(w, http.StatusMethodNotAllowed, "method tidak didukung")
		return
	}
	if h.Service.Repository.DB == nil {
		fail(w, http.StatusInternalServerError, "database belum siap")
		return
	}
	rows, err := h.Service.Suppliers(r.Context())
	respond(w, http.StatusOK, rows, err)
}

func (h Handler) Products(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/products" {
		fail(w, http.StatusNotFound, "endpoint tidak ditemukan")
		return
	}
	switch r.Method {
	case http.MethodGet:
		if h.Service.Repository.DB == nil {
			fail(w, http.StatusInternalServerError, "database belum siap")
			return
		}
		rows, err := h.Service.List(r.Context())
		respond(w, http.StatusOK, rows, err)
	case http.MethodPost:
		if h.Service.Repository.DB == nil {
			fail(w, http.StatusInternalServerError, "database belum siap")
			return
		}
		p, ok := decodeBody(w, r)
		if !ok {
			return
		}
		rows, err := h.Service.Create(r.Context(), p)
		respond(w, http.StatusCreated, rows, err)
	default:
		w.Header().Set("Allow", "GET, POST")
		fail(w, http.StatusMethodNotAllowed, "method tidak didukung")
	}
}

func (h Handler) ProductItem(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r.URL.Path, "/products/")
	fmt.Println("ProductItem ID:", id)
	if !ok {
		return
	}
	switch r.Method {
	case http.MethodGet:
		if h.Service.Repository.DB == nil {
			fail(w, http.StatusInternalServerError, "database belum siap")
			return
		}
		rows, err := h.Service.Get(r.Context(), id)
		respond(w, http.StatusOK, rows, err)
	case http.MethodPut:
		if h.Service.Repository.DB == nil {
			fail(w, http.StatusInternalServerError, "database belum siap")
			return
		}
		p, ok := decodeBody(w, r)
		if !ok {
			return
		}
		rows, err := h.Service.Update(r.Context(), id, p)
		respond(w, http.StatusOK, rows, err)
	case http.MethodDelete:
		if h.Service.Repository.DB == nil {
			fail(w, http.StatusInternalServerError, "database belum siap")
			return
		}
		rows, err := h.Service.Delete(r.Context(), id)
		respond(w, http.StatusOK, rows, err)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		fail(w, http.StatusMethodNotAllowed, "method tidak didukung")
	}
}

// ServeHTTP satisfies http.Handler delegation to Routes mux
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Routes().ServeHTTP(w, r)
}
