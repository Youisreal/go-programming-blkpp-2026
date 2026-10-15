package product

import (
	"math"
	"testing"
)

func TestProductValidation(t *testing.T) {
	valid := Product{ProductID: 30000, ProductName: "Workshop", UnitPrice: 10, UnitsInStock: 5}
	if !valid.Valid() {
		t.Fatal("valid product rejected")
	}
	for _, price := range []float64{-1, math.Inf(1), math.NaN(), math.MaxFloat64} {
		invalid := valid
		invalid.UnitPrice = price
		if invalid.Valid() {
			t.Fatalf("invalid price accepted: %v", price)
		}
	}
	invalid := valid
	invalid.UnitsInStock = 32768
	if invalid.Valid() {
		t.Fatal("smallint overflow accepted")
	}
}
