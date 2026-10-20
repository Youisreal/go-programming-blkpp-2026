package customer

import (
	"net/http"

	"northwind-api/internal/httpjson"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	values, err := h.service.GetAll(r.Context())
	if err != nil {
		httpjson.DatabaseError(w, err, "CUSTOMERS")
		return
	}
	httpjson.Data(w, http.StatusOK, values)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := httpjson.IDString(w, r)
	if !ok {
		return
	}
	value, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		httpjson.DatabaseError(w, err, "CUSTOMERS")
		return
	}
	httpjson.Data(w, http.StatusOK, value)
}
