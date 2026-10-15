package category

import (
	"net/http"

	"northwind-api/internal/httpjson"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("GET /api/categories", h.GetAll)
	mux.HandleFunc("GET /api/categories/{id}", h.GetByID)
	mux.HandleFunc("POST /api/categories", h.Create)
	mux.HandleFunc("PUT /api/categories/{id}", h.Update)
	mux.HandleFunc("DELETE /api/categories/{id}", h.Delete)
	httpjson.Fallbacks(mux, "/api/categories", "GET, HEAD, POST", "GET, HEAD, PUT, DELETE")
}
