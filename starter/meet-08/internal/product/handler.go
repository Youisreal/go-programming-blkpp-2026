package product

import (
	"encoding/json"
	"errors"
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

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var id int
	collection := r.URL.Path == "/products"
	switch {
	case collection:
	case strings.HasPrefix(r.URL.Path, "/products/"):
		var err error
		id, err = strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/products/"))
		if err != nil || id <= 0 {
			fail(w, http.StatusBadRequest, "ID produk harus bilangan bulat positif")
			return
		}
	case r.URL.Path == "/categories" || r.URL.Path == "/suppliers":
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			fail(w, http.StatusMethodNotAllowed, "method tidak didukung")
			return
		}
		var rows json.RawMessage
		var err error
		if r.URL.Path == "/categories" {
			rows, err = h.Service.Categories(r.Context())
		} else {
			rows, err = h.Service.Suppliers(r.Context())
		}
		respond(w, http.StatusOK, rows, err)
		return
	default:
		fail(w, http.StatusNotFound, "endpoint tidak ditemukan")
		return
	}

	var rows json.RawMessage
	var err error
	status := http.StatusOK
	switch {
	case r.Method == http.MethodGet && collection:
		rows, err = h.Service.List(r.Context())
	case r.Method == http.MethodGet:
		rows, err = h.Service.Get(r.Context(), id)
	case (r.Method == http.MethodPost && collection) || (r.Method == http.MethodPut && !collection):
		var p Input
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&p); err != nil {
			fail(w, http.StatusBadRequest, "body harus JSON produk valid, maksimal 1 MiB")
			return
		}
		var extra any
		if err := decoder.Decode(&extra); err != io.EOF {
			fail(w, http.StatusBadRequest, "body harus berisi satu objek JSON")
			return
		}
		if collection {
			rows, err = h.Service.Create(r.Context(), p)
			status = http.StatusCreated
		} else {
			rows, err = h.Service.Update(r.Context(), id, p)
		}
	case r.Method == http.MethodDelete && !collection:
		rows, err = h.Service.Delete(r.Context(), id)
	default:
		allow := "GET, PUT, DELETE"
		if collection {
			allow = "GET, POST"
		}
		w.Header().Set("Allow", allow)
		fail(w, http.StatusMethodNotAllowed, "method tidak didukung")
		return
	}
	respond(w, status, rows, err)
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
