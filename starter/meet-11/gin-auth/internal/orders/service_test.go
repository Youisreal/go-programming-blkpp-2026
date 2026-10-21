package orders

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	"gin-auth/internal/auth"
	"gin-auth/internal/database"
	"gin-auth/internal/product"
)

func TestOrderValidation(t *testing.T) {
	s := NewService(NewRepository(nil))
	for _, input := range []struct {
		id   uint
		qty  int
		user uint
	}{{0, 1, 1}, {1, 0, 1}, {1, -1, 1}, {1, 1, 0}} {
		if _, err := s.Create(context.Background(), input.id, input.qty, auth.User{ID: input.user}); !errors.Is(err, ErrInvalidOrder) {
			t.Fatalf("invalid order: %v", err)
		}
	}
	for _, input := range []struct{ price, qty int }{{-1, 1}, {1, 0}, {math.MaxInt, 2}} {
		if _, err := orderTotal(input.price, input.qty); !errors.Is(err, ErrInvalidOrder) {
			t.Fatalf("invalid total: %v", err)
		}
	}
	if total, err := orderTotal(15000, 2); err != nil || total != 30000 {
		t.Fatalf("total=%d err=%v", total, err)
	}
}

func TestOrderTransaction(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("Isi TEST_DATABASE_URL dengan database PostgreSQL khusus test")
	}
	db, err := database.Open(dsn)
	if err != nil {
		t.Fatal(err)
	}
	conn, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if err := db.AutoMigrate(&auth.User{}, &product.Product{}, &Order{}, &OrderDetail{}); err != nil {
		t.Fatal(err)
	}
	user := auth.User{Name: "Order test", Email: fmt.Sprintf("order-%d@example.com", time.Now().UnixNano()), PasswordHash: "unused-test-hash"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	item := product.Product{Name: "Buku", UserID: user.ID, Stock: 3, Price: 15000}
	t.Cleanup(func() {
		for _, result := range []error{
			db.Where("order_id IN (?)", db.Model(&Order{}).Select("id").Where("user_id = ?", user.ID)).Delete(&OrderDetail{}).Error,
			db.Where("user_id = ?", user.ID).Delete(&Order{}).Error,
			db.Where("user_id = ?", user.ID).Delete(&product.Product{}).Error,
			db.Delete(&user).Error,
		} {
			if result != nil {
				t.Error(result)
			}
		}
	})
	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}
	s := NewService(NewRepository(db))
	ctx := context.Background()
	order, err := s.Create(ctx, item.ID, 2, user)
	if err != nil {
		t.Fatal(err)
	}
	if order.Total != 30000 || len(order.Details) != 1 || order.Details[0].Qty != 2 || order.Details[0].Price != 15000 {
		t.Fatalf("order=%+v", order)
	}
	if _, err := s.Create(ctx, item.ID, 2, user); !errors.Is(err, ErrInsufficientStock) {
		t.Fatalf("stock: %v", err)
	}
	if err := db.First(&item, item.ID).Error; err != nil || item.Stock != 1 {
		t.Fatalf("stock=%d err=%v", item.Stock, err)
	}
	list, err := s.List(ctx, user.ID)
	if err != nil || len(list) != 1 {
		t.Fatalf("list=%v err=%v", list, err)
	}
	hidden, err := s.GetByID(ctx, order.ID, 0)
	if err != nil || hidden != nil {
		t.Fatalf("order pengguna lain terbaca: %v %v", hidden, err)
	}
	list, err = s.List(ctx, 0)
	if err != nil || len(list) != 0 {
		t.Fatalf("order pengguna lain terbaca: %v %v", list, err)
	}
}
