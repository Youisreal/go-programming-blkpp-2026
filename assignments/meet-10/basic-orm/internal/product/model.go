package product

import "strings"

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (p Product) Valid() bool {
	return strings.TrimSpace(p.Name) != "" && p.Price >= 0
}
