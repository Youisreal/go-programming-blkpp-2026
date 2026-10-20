package customers

import (
	"strings"
)

type Customer struct {
	CustomerID   string `json:"customer_id"`
	CompanyName  string `json:"company_name"`
	CustomerName string `json:"contact_name"`
	ContactTitle string `json:"contact_title"`
}

func (c Customer) Valid() bool {
	return strings.TrimSpace(c.CustomerID) != "" && strings.TrimSpace(c.CustomerName) != ""
}
