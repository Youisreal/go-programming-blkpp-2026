package order

import (
	"net/http"

	"northwind-api/internal/httpjson"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("GET /api/orders", h.GetAll)
	mux.HandleFunc("GET /api/orders/{id}", h.GetByID)
	mux.HandleFunc("POST /api/orders", h.Create)
	httpjson.Fallbacks(mux, "/api/orders", "GET, HEAD, POST", "GET, HEAD")
}
