package main

type Customer struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type TransactionItem struct {
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	UnitPrice int64  `json:"unit_price"`
	Quantity  int64  `json:"quantity"`
}

type Product struct {
	SKU    string `json:"sku"`
	Name   string `json:"name"`
	Price  int64  `json:"price"`
	Active bool   `json:"active"`
}

type StockItem struct {
	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Quantity int64  `json:"quantity"`
}

type StockAdjustment struct {
	SKU      string `json:"sku"`
	Quantity int64  `json:"quantity"`
}

type CheckoutItem struct {
	SKU      string `json:"sku"`
	Quantity int64  `json:"quantity"`
}

type CheckoutRequest struct {
	InvoiceNumber   string         `json:"invoice_number"`
	Customer        Customer       `json:"customer"`
	Items           []CheckoutItem `json:"items"`
	DiscountPercent int64          `json:"discount_percent"`
	TaxPercent      int64          `json:"tax_percent"`
	ShippingFee     int64          `json:"shipping_fee"`
	AmountPaid      int64          `json:"amount_paid"`
}

type InvoiceRequest struct {
	InvoiceNumber   string            `json:"invoice_number"`
	Customer        Customer          `json:"customer"`
	Items           []TransactionItem `json:"items"`
	DiscountPercent int64             `json:"discount_percent"`
	TaxPercent      int64             `json:"tax_percent"`
	ShippingFee     int64             `json:"shipping_fee"`
	AmountPaid      int64             `json:"amount_paid"`
}

type InvoiceItemResult struct {
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	Quantity  int64  `json:"quantity"`
	UnitPrice int64  `json:"unit_price"`
	LineTotal int64  `json:"line_total"`
}

type InvoiceSummary struct {
	Subtotal       int64  `json:"subtotal"`
	DiscountAmount int64  `json:"discount_amount"`
	AfterDiscount  int64  `json:"after_discount"`
	TaxAmount      int64  `json:"tax_amount"`
	ShippingFee    int64  `json:"shipping_fee"`
	GrandTotal     int64  `json:"grand_total"`
	AmountPaid     int64  `json:"amount_paid"`
	Change         int64  `json:"change"`
	Status         string `json:"status"`
}

type InvoiceResult struct {
	InvoiceNumber string              `json:"invoice_number"`
	CustomerName  string              `json:"customer_name"`
	Items         []InvoiceItemResult `json:"items"`
	Summary       InvoiceSummary      `json:"summary"`
}
