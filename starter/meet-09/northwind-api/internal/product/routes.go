package product

import (
	"net/http"

	"northwind-api/internal/httpjson"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("GET /api/products", h.GetAll)
	mux.HandleFunc("GET /api/products/{id}", h.GetByID)
	mux.HandleFunc("POST /api/products", h.Create)
	mux.HandleFunc("PUT /api/products/{id}", h.Update)
	mux.HandleFunc("DELETE /api/products/{id}", h.Delete)
	httpjson.Fallbacks(mux, "/api/products", "GET, HEAD, POST", "GET, HEAD, PUT, DELETE")
}
