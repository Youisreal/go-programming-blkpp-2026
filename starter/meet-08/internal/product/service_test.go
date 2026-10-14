package product

import (
	"context"
	"errors"
	"math"
	"strings"
	"testing"
)

func TestServiceValidation(t *testing.T) {
	// Tanpa DB: input salah harus ditolak sebelum repository dipanggil.
	s := Service{}
	for _, input := range []Input{
		{},
		{ProductName: "   ", SupplierID: 1, CategoryID: 1},
		{ProductName: strings.Repeat("界", 41), SupplierID: 1, CategoryID: 1},
		{ProductName: "Tea\x00", SupplierID: 1, CategoryID: 1},
		{ProductName: "Tea", SupplierID: 1, CategoryID: 1, UnitPrice: math.NaN()},
		{ProductName: "Tea", SupplierID: 1, CategoryID: 1, UnitPrice: math.Inf(1)},
		{ProductName: "Tea", SupplierID: 1, CategoryID: 1, UnitPrice: 10000000000},
		{ProductName: "Tea", SupplierID: 1, CategoryID: 1, UnitsInStock: -1},
		{ProductName: "Tea", SupplierID: 1, CategoryID: 1, Discontinued: 2},
	} {
		if _, err := s.Create(context.Background(), input); !errors.Is(err, ErrInvalidInput) {
			t.Fatalf("Create menerima input tidak valid: %+v err=%v", input, err)
		}
		if _, err := s.Update(context.Background(), 1, input); !errors.Is(err, ErrInvalidInput) {
			t.Fatalf("Update menerima input tidak valid: %+v err=%v", input, err)
		}
	}
	p := Input{ProductName: "  " + strings.Repeat("界", 40) + "  ", SupplierID: 1, CategoryID: 1}
	if err := normalize(&p); err != nil || p.ProductName != strings.Repeat("界", 40) {
		t.Fatalf("normalisasi/validasi karakter Unicode: %+v err=%v", p, err)
	}
}
