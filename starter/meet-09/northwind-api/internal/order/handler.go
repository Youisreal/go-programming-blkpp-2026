package order

import (
	"errors"
	"net/http"

	"northwind-api/internal/httpjson"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	values, err := h.service.GetAll(r.Context())
	if err != nil {
		httpjson.DatabaseError(w, err, "ORDER")
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
		httpjson.DatabaseError(w, err, "ORDER")
		return
	}
	httpjson.Data(w, http.StatusOK, value)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateRequest
	if !httpjson.Decode(w, r, &input) {
		return
	}
	value, err := h.service.Create(r.Context(), input)
	switch {
	case errors.Is(err, ErrInvalidInput):
		httpjson.Error(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid order fields or duplicate products")
	case errors.Is(err, ErrCustomerNotFound):
		httpjson.Error(w, http.StatusNotFound, "CUSTOMER_NOT_FOUND", "Customer not found")
	case errors.Is(err, ErrProductNotFound):
		httpjson.Error(w, http.StatusNotFound, "PRODUCT_NOT_FOUND", "Product not found")
	case errors.Is(err, ErrStock):
		httpjson.Error(w, http.StatusConflict, "INSUFFICIENT_STOCK", "Insufficient stock or discontinued product")
	case errors.Is(err, ErrPrice):
		httpjson.Error(w, http.StatusConflict, "INVALID_PRODUCT_PRICE", "Product has no valid price")
	case err != nil:
		httpjson.DatabaseError(w, err, "ORDER")
	default:
		httpjson.Data(w, http.StatusCreated, value)
	}
}
