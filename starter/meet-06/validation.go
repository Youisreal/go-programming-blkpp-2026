package main

import (
	"fmt"
	"strings"
)

func ValidatePercent(name string, percent int64) error {
	if percent < 0 || percent > 100 {
		return fmt.Errorf("%s harus berada di rentang 0 sampai 100", name)
	}
	return nil
}

func ValidateProduct(product Product) error {
	if strings.TrimSpace(product.SKU) == "" {
		return fmt.Errorf("sku tidak boleh kosong")
	}
	if strings.TrimSpace(product.Name) == "" {
		return fmt.Errorf("name tidak boleh kosong")
	}
	if product.Price < 0 {
		return fmt.Errorf("price tidak boleh negatif")
	}
	return nil
}

func ValidateStockAdjustment(adjustment StockAdjustment) error {
	if strings.TrimSpace(adjustment.SKU) == "" {
		return fmt.Errorf("sku tidak boleh kosong")
	}
	if adjustment.Quantity == 0 {
		return fmt.Errorf("quantity tidak boleh 0")
	}
	return nil
}

func ValidateCheckoutRequest(request CheckoutRequest) error {
	if strings.TrimSpace(request.InvoiceNumber) == "" {
		return fmt.Errorf("invoice_number tidak boleh kosong")
	}
	if strings.TrimSpace(request.Customer.Name) == "" {
		return fmt.Errorf("customer.name tidak boleh kosong")
	}
	if len(request.Items) == 0 {
		return fmt.Errorf("items tidak boleh kosong")
	}
	if err := ValidatePercent("discount_percent", request.DiscountPercent); err != nil {
		return err
	}
	if err := ValidatePercent("tax_percent", request.TaxPercent); err != nil {
		return err
	}
	if request.ShippingFee < 0 {
		return fmt.Errorf("shipping_fee tidak boleh negatif")
	}
	if request.AmountPaid < 0 {
		return fmt.Errorf("amount_paid tidak boleh negatif")
	}
	for i, item := range request.Items {
		if strings.TrimSpace(item.SKU) == "" {
			return fmt.Errorf("items[%d].sku tidak boleh kosong", i)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("items[%d].quantity harus lebih dari 0", i)
		}
	}
	return nil
}

func ValidateRequest(request InvoiceRequest) error {
	if strings.TrimSpace(request.InvoiceNumber) == "" {
		return fmt.Errorf("invoice_number tidak boleh kosong")
	}
	if strings.TrimSpace(request.Customer.Name) == "" {
		return fmt.Errorf("customer.name tidak boleh kosong")
	}
	if len(request.Items) == 0 {
		return fmt.Errorf("items tidak boleh kosong")
	}
	if err := ValidatePercent("discount_percent", request.DiscountPercent); err != nil {
		return err
	}
	if err := ValidatePercent("tax_percent", request.TaxPercent); err != nil {
		return err
	}
	if request.ShippingFee < 0 {
		return fmt.Errorf("shipping_fee tidak boleh negatif")
	}
	if request.AmountPaid < 0 {
		return fmt.Errorf("amount_paid tidak boleh negatif")
	}

	for i, item := range request.Items {
		if strings.TrimSpace(item.SKU) == "" {
			return fmt.Errorf("items[%d].sku tidak boleh kosong", i)
		}
		if strings.TrimSpace(item.Name) == "" {
			return fmt.Errorf("items[%d].name tidak boleh kosong", i)
		}
		if _, err := CalculateLineTotal(item.UnitPrice, item.Quantity); err != nil {
			return fmt.Errorf("items[%d] tidak valid: %w", i, err)
		}
	}
	return nil
}
