package customers

import (
	"net/http"

	"northwind-api/internal/httpjson"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("GET /api/customers", h.GetAll)
	mux.HandleFunc("GET /api/customers/{id}", h.GetByID)
	mux.HandleFunc("POST /api/customers", h.Create)
	mux.HandleFunc("PUT /api/customers/{id}", h.Update)
	mux.HandleFunc("DELETE /api/customers/{id}", h.Delete)
	httpjson.Fallbacks(mux, "/api/customers", "GET, HEAD, POST", "GET, HEAD, PUT, DELETE")
}
