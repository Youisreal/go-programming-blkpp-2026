package category

import "strings"

type Category struct {
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Description  *string `json:"description"`
}

func (c Category) Valid() bool {
	return c.CategoryID > 0 && c.CategoryID <= 32767 && strings.TrimSpace(c.CategoryName) != "" && len([]rune(c.CategoryName)) <= 15
}
