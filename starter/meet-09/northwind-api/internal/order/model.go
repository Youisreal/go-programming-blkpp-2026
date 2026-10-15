package order

import (
	"errors"
	"math"
	"strings"
	"time"
)

var (
	ErrInvalidInput     = errors.New("invalid order input")
	ErrCustomerNotFound = errors.New("customer not found")
	ErrProductNotFound  = errors.New("product not found")
	ErrStock            = errors.New("insufficient stock or discontinued product")
	ErrPrice            = errors.New("product has no valid price")
)

type CreateRequest struct {
	OrderID    int           `json:"order_id"`
	CustomerID string        `json:"customer_id"`
	Items      []ItemRequest `json:"items"`
}

type ItemRequest struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Discount  float64 `json:"discount"`
}

type Item struct {
	ProductID int     `json:"product_id"`
	UnitPrice float64 `json:"unit_price"`
	Quantity  int     `json:"quantity"`
	Discount  float64 `json:"discount"`
	Subtotal  float64 `json:"subtotal"`
}

type Order struct {
	OrderID    int        `json:"order_id"`
	CustomerID *string    `json:"customer_id"`
	OrderDate  *time.Time `json:"order_date"`
	Total      float64    `json:"total"`
	Items      []Item     `json:"items,omitempty"`
}

func Subtotal(price float64, quantity int, discount float64) float64 {
	return price * float64(quantity) * (1 - discount)
}

func Validate(input CreateRequest) error {
	if input.OrderID < 1 || input.OrderID > 32767 || strings.TrimSpace(input.CustomerID) == "" || len([]rune(input.CustomerID)) > 5 || len(input.Items) == 0 {
		return ErrInvalidInput
	}
	seen := make(map[int]bool, len(input.Items))
	for _, item := range input.Items {
		if item.ProductID < 1 || item.ProductID > 32767 || item.Quantity < 1 || item.Quantity > 32767 || math.IsNaN(item.Discount) || math.IsInf(item.Discount, 0) || item.Discount < 0 || item.Discount > 1 || seen[item.ProductID] {
			return ErrInvalidInput
		}
		seen[item.ProductID] = true
	}
	return nil
}
