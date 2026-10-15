// Package httpjson berisi format JSON dan validasi HTTP yang dipakai semua module.
package httpjson

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/lib/pq"
)

func JSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func Data(w http.ResponseWriter, status int, value any) {
	JSON(w, status, map[string]any{"data": value})
}

func Error(w http.ResponseWriter, status int, code, message string) {
	JSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}

func ID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 || id > 32767 {
		Error(w, http.StatusBadRequest, "INVALID_ID", "ID must be between 1 and 32767")
		return 0, false
	}
	return id, true
}

func Decode(w http.ResponseWriter, r *http.Request, value any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
		return false
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		Error(w, http.StatusBadRequest, "INVALID_JSON", "Expected one JSON value")
		return false
	}
	return true
}

func DatabaseError(w http.ResponseWriter, err error, entity string) {
	var pgErr *pq.Error
	switch {
	case errors.Is(err, sql.ErrNoRows):
		Error(w, http.StatusNotFound, entity+"_NOT_FOUND", "Record not found")
	case errors.As(err, &pgErr) && pgErr.Code == "23505":
		Error(w, http.StatusConflict, "DUPLICATE_KEY", "Record already exists")
	case errors.As(err, &pgErr) && pgErr.Code == "23503":
		Error(w, http.StatusConflict, "FOREIGN_KEY_CONFLICT", "Related record missing or still in use")
	case errors.As(err, &pgErr) && (pgErr.Code == "23514" || pgErr.Code == "23502" || pgErr.Code == "22003" || pgErr.Code == "22001"):
		Error(w, http.StatusBadRequest, "INVALID_INPUT", "Input violates database constraints")
	default:
		log.Printf("%s database error: %v", entity, err)
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}

// Fallbacks menjaga error routing tetap JSON; ServeMux tetap menangani HEAD sebagai GET.
func Fallbacks(mux *http.ServeMux, base string, methods, itemMethods string) {
	for _, route := range []struct{ path, allow string }{{base, methods}, {base + "/{id}", itemMethods}} {
		mux.HandleFunc(route.path, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Allow", route.allow)
			Error(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		})
	}
}
