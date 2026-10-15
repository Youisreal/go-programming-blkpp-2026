package product

import (
	"math"
	"strings"
)

type Product struct {
	ProductID    int     `json:"product_id"`
	ProductName  string  `json:"product_name"`
	SupplierID   *int    `json:"supplier_id"`
	CategoryID   *int    `json:"category_id"`
	UnitPrice    float64 `json:"unit_price"`
	UnitsInStock int     `json:"units_in_stock"`
	Discontinued int     `json:"discontinued"`
}

func (p Product) Valid() bool {
	return p.ProductID > 0 && p.ProductID <= 32767 && strings.TrimSpace(p.ProductName) != "" && len([]rune(p.ProductName)) <= 40 &&
		!math.IsNaN(p.UnitPrice) && !math.IsInf(p.UnitPrice, 0) && p.UnitPrice >= 0 && p.UnitPrice <= math.MaxFloat32 &&
		p.UnitsInStock >= 0 && p.UnitsInStock <= 32767 && (p.Discontinued == 0 || p.Discontinued == 1) && validOptionalID(p.SupplierID) && validOptionalID(p.CategoryID)
}

func validOptionalID(id *int) bool { return id == nil || (*id > 0 && *id <= 32767) }
