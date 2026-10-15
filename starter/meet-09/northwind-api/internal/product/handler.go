package product

import (
	"net/http"

	"northwind-api/internal/httpjson"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	values, err := h.service.GetAll(r.Context())
	if err != nil {
		httpjson.DatabaseError(w, err, "PRODUCT")
		return
	}
	httpjson.Data(w, http.StatusOK, values)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := httpjson.ID(w, r)
	if !ok {
		return
	}
	value, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		httpjson.DatabaseError(w, err, "PRODUCT")
		return
	}
	httpjson.Data(w, http.StatusOK, value)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var value Product
	if !httpjson.Decode(w, r, &value) {
		return
	}
	if !value.Valid() {
		httpjson.Error(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid product fields")
		return
	}
	if err := h.service.Create(r.Context(), value); err != nil {
		httpjson.DatabaseError(w, err, "PRODUCT")
		return
	}
	httpjson.Data(w, http.StatusCreated, value)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := httpjson.ID(w, r)
	if !ok {
		return
	}
	var value Product
	if !httpjson.Decode(w, r, &value) {
		return
	}
	if value.ProductID != 0 && value.ProductID != id {
		httpjson.Error(w, http.StatusBadRequest, "INVALID_INPUT", "Body ID must match path ID")
		return
	}
	value.ProductID = id
	if !value.Valid() {
		httpjson.Error(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid product fields")
		return
	}
	if err := h.service.Update(r.Context(), value); err != nil {
		httpjson.DatabaseError(w, err, "PRODUCT")
		return
	}
	httpjson.Data(w, http.StatusOK, value)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := httpjson.ID(w, r)
	if !ok {
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		httpjson.DatabaseError(w, err, "PRODUCT")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
