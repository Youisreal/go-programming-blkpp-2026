package customer

import "strings"

type Customer struct {
	CustomerID   string  `json:"customer_id"`
	CustomerName string  `json:"customer_name"`
	ContactName  *string `json:"contact_name"`
	ContactTitle *string `json:"contact_title"`
}

func (c Customer) Valid() bool {
	return strings.TrimSpace(c.CustomerID) != "" && strings.TrimSpace(c.CustomerName) != "" && len([]rune(c.CustomerName)) <= 15
}
