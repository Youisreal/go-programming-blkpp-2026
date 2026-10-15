package order

import (
	"context"
	"errors"
	"math"
	"testing"
)

func TestSubtotal(t *testing.T) {
	for _, tc := range []struct {
		price          float64
		quantity       int
		discount, want float64
	}{
		{18, 2, 0, 36}, {10, 3, 0.1, 27}, {25, 2, 1, 0}, {0, 2, 0.5, 0},
	} {
		if got := Subtotal(tc.price, tc.quantity, tc.discount); math.Abs(got-tc.want) > 1e-9 {
			t.Fatalf("got %v want %v", got, tc.want)
		}
	}
}

func TestValidate(t *testing.T) {
	valid := func() CreateRequest {
		return CreateRequest{OrderID: 30000, CustomerID: "ALFKI", Items: []ItemRequest{{ProductID: 1, Quantity: 2, Discount: 0.1}}}
	}
	if err := Validate(valid()); err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(*CreateRequest){
		func(r *CreateRequest) { r.OrderID = 32768 },
		func(r *CreateRequest) { r.CustomerID = "" },
		func(r *CreateRequest) { r.CustomerID = "TOOLONG" },
		func(r *CreateRequest) { r.Items = nil },
		func(r *CreateRequest) { r.Items[0].ProductID = 0 },
		func(r *CreateRequest) { r.Items[0].Quantity = 0 },
		func(r *CreateRequest) { r.Items[0].Quantity = 32768 },
		func(r *CreateRequest) { r.Items[0].Discount = -0.1 },
		func(r *CreateRequest) { r.Items[0].Discount = 1.1 },
		func(r *CreateRequest) { r.Items[0].Discount = math.NaN() },
		func(r *CreateRequest) { r.Items[0].Discount = math.Inf(1) },
		func(r *CreateRequest) { r.Items = append(r.Items, r.Items[0]) },
	} {
		request := valid()
		mutate(&request)
		// Invalid input harus ditolak sebelum akses repository/transaction.
		_, err := NewService(nil).Create(context.Background(), request)
		if !errors.Is(err, ErrInvalidInput) {
			t.Fatalf("expected invalid input, got %v", err)
		}
	}
}
