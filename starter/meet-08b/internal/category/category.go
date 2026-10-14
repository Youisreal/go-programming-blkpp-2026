package category

// Category represents a product category row from the database.
type Category struct {
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
	Description  string `json:"description"`
}
