package orders

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http/httptest"
	"os"
	"strings"
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
		if _, err := s.Create(context.Background(), []ItemInput{{input.id, input.qty}}, auth.User{ID: input.user}); !errors.Is(err, ErrInvalidOrder) {
			t.Fatalf("invalid order: %v", err)
		}
	}
	for _, items := range [][]ItemInput{nil, {}, {{1, 1}, {1, 2}}, {{1, 1}, {2, 0}}} {
		if _, err := s.Create(context.Background(), items, auth.User{ID: 1}); !errors.Is(err, ErrInvalidOrder) {
			t.Fatalf("invalid items: %v", err)
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
	second := product.Product{Name: "Pensil", UserID: user.ID, Stock: 4, Price: 5000}
	if err := db.Create(&second).Error; err != nil {
		t.Fatal(err)
	}

	s := NewService(NewRepository(db))
	ctx := context.Background()
	order, err := s.Create(ctx, []ItemInput{{item.ID, 2}, {second.ID, 1}}, user)
	if err != nil {
		t.Fatal(err)
	}
	if order.GrandTotal != 35000 || len(order.Details) != 2 || order.Details[0].Qty != 2 || order.Details[0].Price != 15000 || order.Details[0].Subtotal != 30000 || order.Details[1].Subtotal != 5000 {
		t.Fatalf("order=%+v", order)
	}
	if _, err := s.Create(ctx, []ItemInput{{item.ID, 1}, {second.ID, 4}}, user); !errors.Is(err, ErrInsufficientStock) {
		t.Fatalf("stock: %v", err)
	}
	if err := db.First(&item, item.ID).Error; err != nil || item.Stock != 1 {
		t.Fatalf("stock=%d err=%v", item.Stock, err)
	}
	if err := db.First(&second, second.ID).Error; err != nil || second.Stock != 3 {
		t.Fatalf("second stock=%d err=%v", second.Stock, err)
	}
	found, err := s.GetByID(ctx, order.ID, user.ID)
	if err != nil || found == nil || found.GrandTotal != 35000 || len(found.Details) != 2 || found.Details[0].Subtotal != 30000 {
		t.Fatalf("detail=%+v err=%v", found, err)
	}
	// Overflow penjumlahan antar item juga harus rollback pengurangan stok.
	if err := db.Model(&second).Update("price", math.MaxInt).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := s.Create(ctx, []ItemInput{{item.ID, 1}, {second.ID, 1}}, user); !errors.Is(err, ErrInvalidOrder) {
		t.Fatalf("grand total overflow: %v", err)
	}
	if err := db.First(&item, item.ID).Error; err != nil || item.Stock != 1 {
		t.Fatalf("rollback stock=%d err=%v", item.Stock, err)
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

func TestOrderRequestValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("user", auth.User{ID: 1}) })
	router.POST("/orders", NewHandler(NewService(NewRepository(nil))).Create)
	for _, body := range []string{
		`{}`, `{"items":[]}`, `{"product_id":1,"quantity":2}`,
		`{"items":[{"product_id":1,"quantity":0}]}`,
		`{"items":[{"product_id":0,"quantity":1}]}`,
		`{"items":[{"product_id":1,"quantity":1},{"product_id":1,"quantity":2}]}`,
	} {
		req := httptest.NewRequest("POST", "/orders", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)
		if response.Code != 400 {
			t.Fatalf("body=%s status=%d", body, response.Code)
		}
	}
}

func TestOrderDetailSubtotal(t *testing.T) {
	detail := OrderDetail{Price: 15000, Qty: 2}
	if err := detail.AfterFind(nil); err != nil || detail.Subtotal != 30000 {
		t.Fatalf("subtotal=%d err=%v", detail.Subtotal, err)
	}
}
