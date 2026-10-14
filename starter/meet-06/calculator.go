package main

import (
	"fmt"
	"math"
)

func CalculateLineTotal(price int64, qty int64) (int64, error) {
	if price < 0 {
		return 0, fmt.Errorf("unit_price tidak boleh negatif")
	}
	if qty <= 0 {
		return 0, fmt.Errorf("quantity harus lebih dari 0")
	}
	if price > math.MaxInt64/qty {
		return 0, fmt.Errorf("line_total overflow")
	}
	return price * qty, nil
}

func CalculateDiscount(subtotal int64, percent int64) (int64, error) {
	if subtotal < 0 {
		return 0, fmt.Errorf("subtotal tidak boleh negatif")
	}
	if err := ValidatePercent("discount_percent", percent); err != nil {
		return 0, err
	}
	return calculatePercent(subtotal, percent)
}

func CalculateTax(amount int64, percent int64) (int64, error) {
	if amount < 0 {
		return 0, fmt.Errorf("amount tidak boleh negatif")
	}
	if err := ValidatePercent("tax_percent", percent); err != nil {
		return 0, err
	}
	return calculatePercent(amount, percent)
}

func CalculateChange(grandTotal int64, amountPaid int64) (int64, error) {
	if grandTotal < 0 {
		return 0, fmt.Errorf("grand_total tidak boleh negatif")
	}
	if amountPaid < 0 {
		return 0, fmt.Errorf("amount_paid tidak boleh negatif")
	}
	if amountPaid < grandTotal {
		return 0, fmt.Errorf("amount_paid kurang dari grand_total")
	}
	return amountPaid - grandTotal, nil
}

func CalculateInvoice(request InvoiceRequest) (InvoiceResult, error) {
	if err := ValidateRequest(request); err != nil {
		return InvoiceResult{}, err
	}

	items := make([]InvoiceItemResult, 0, len(request.Items))
	var subtotal int64
	for _, item := range request.Items {
		lineTotal, err := CalculateLineTotal(item.UnitPrice, item.Quantity)
		if err != nil {
			return InvoiceResult{}, err
		}
		if subtotal > math.MaxInt64-lineTotal {
			return InvoiceResult{}, fmt.Errorf("subtotal overflow")
		}
		subtotal += lineTotal
		items = append(items, InvoiceItemResult{
			SKU:       item.SKU,
			Name:      item.Name,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			LineTotal: lineTotal,
		})
	}

	discountAmount, err := CalculateDiscount(subtotal, request.DiscountPercent)
	if err != nil {
		return InvoiceResult{}, err
	}
	afterDiscount := subtotal - discountAmount

	taxAmount, err := CalculateTax(afterDiscount, request.TaxPercent)
	if err != nil {
		return InvoiceResult{}, err
	}
	grandTotal, err := addMoney(afterDiscount, taxAmount, request.ShippingFee)
	if err != nil {
		return InvoiceResult{}, err
	}
	change, err := CalculateChange(grandTotal, request.AmountPaid)
	if err != nil {
		return InvoiceResult{}, err
	}

	return InvoiceResult{
		InvoiceNumber: request.InvoiceNumber,
		CustomerName:  request.Customer.Name,
		Items:         items,
		Summary: InvoiceSummary{
			Subtotal:       subtotal,
			DiscountAmount: discountAmount,
			AfterDiscount:  afterDiscount,
			TaxAmount:      taxAmount,
			ShippingFee:    request.ShippingFee,
			GrandTotal:     grandTotal,
			AmountPaid:     request.AmountPaid,
			Change:         change,
			Status:         "paid",
		},
	}, nil
}

func calculatePercent(amount int64, percent int64) (int64, error) {
	base := (amount / 100) * percent
	remainder := (amount % 100) * percent
	rounded := remainder / 100
	if remainder%100 >= 50 {
		rounded++
	}
	return addMoney(base, rounded)
}

func addMoney(values ...int64) (int64, error) {
	var total int64
	for _, value := range values {
		if value < 0 {
			return 0, fmt.Errorf("nilai uang tidak boleh negatif")
		}
		if total > math.MaxInt64-value {
			return 0, fmt.Errorf("perhitungan uang overflow")
		}
		total += value
	}
	return total, nil
}
