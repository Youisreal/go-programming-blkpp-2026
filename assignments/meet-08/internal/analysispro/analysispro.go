package analysispro

import "errors"

var (
	ErrInvalidInput = errors.New("nama wajib 1–40 karakter; kemasan maksimal 20; supplier/category positif; harga 0–9999999999.99; stok/order/reorder 0–32767; discontinued 0 atau 1")
	ErrNotFound     = errors.New("produk tidak ditemukan")
	ErrForeignKey   = errors.New("supplier/category tidak ditemukan atau data masih dipakai")
	ErrConflict     = errors.New("data dengan ID tersebut sudah ada")
	ErrConstraint   = errors.New("data melanggar batasan database")
)
