package supplier

// Supplier represents a supplier row from the database.
type Supplier struct {
	SupplierID   int    `json:"supplier_id"`
	SupplierName string `json:"supplier_name"`
	ContactName  string `json:"contact_name"`
	Address      string `json:"address"`
	City         string `json:"city"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
	Phone        string `json:"phone"`
}
