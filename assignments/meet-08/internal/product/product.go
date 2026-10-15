package product

import "errors"

type Input struct {
	ProductName     string  `json:"product_name"`
	SupplierID      int     `json:"supplier_id"`
	CategoryID      int     `json:"category_id"`
	QuantityPerUnit string  `json:"quantity_per_unit"`
	UnitPrice       float64 `json:"unit_price"`
	UnitsInStock    int     `json:"units_in_stock"`
	UnitsOnOrder    int     `json:"units_on_order"`
	ReorderLevel    int     `json:"reorder_level"`
	Discontinued    int     `json:"discontinued"`
}

var (
	ErrInvalidInput = errors.New("nama wajib 1–40 karakter; kemasan maksimal 20; supplier/category positif; harga 0–9999999999.99; stok/order/reorder 0–32767; discontinued 0 atau 1")
	ErrNotFound     = errors.New("produk tidak ditemukan")
	ErrForeignKey   = errors.New("supplier/category tidak ditemukan atau data masih dipakai")
	ErrConflict     = errors.New("data dengan ID tersebut sudah ada")
	ErrConstraint   = errors.New("data melanggar batasan database")
)
