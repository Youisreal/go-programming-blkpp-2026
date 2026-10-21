package category

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"gin-auth/internal/auth"
	"gin-auth/internal/database"

	"github.com/gin-gonic/gin"
)

func TestRejectInvalidName(t *testing.T) {
	r := gin.New()
	RegisterRoutes(r.Group("/api/private"), NewHandler(NewService(NewRepository(nil))))
	for _, body := range []string{`{`, `{}`, `{"name":"   "}`, `{"name":"` + strings.Repeat("a", 101) + `"}`} {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/private/products", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		if w.Code != 400 {
			t.Fatalf("status=%d, body=%s", w.Code, w.Body.String())
		}
	}
}

func TestProductRelation(t *testing.T) {
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
	t.Cleanup(func() { conn.Close() })
	if err := db.AutoMigrate(&auth.User{}, &Category{}); err != nil {
		t.Fatal(err)
	}
	user := auth.User{Name: "Product test", Email: fmt.Sprintf("product-%d@example.com", time.Now().UnixNano()), PasswordHash: "unused-test-hash"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Where("user_id = ?", user.ID).Delete(&Category{}).Error; err != nil {
			t.Error(err)
		}
		if err := db.Delete(&user).Error; err != nil {
			t.Error(err)
		}
	})
	r := gin.New()
	currentUser := user
	group := r.Group("/api/private", func(c *gin.Context) { c.Set("user", currentUser) })
	RegisterRoutes(group, NewHandler(NewService(NewRepository(db))))
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/private/products", strings.NewReader(`{"name":" Buku ","user_id":0}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	var created struct {
		Data Category `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if w.Code != 201 || created.Data.UserID != user.ID || created.Data.Name != "Buku" || strings.Contains(w.Body.String(), "password") {
		t.Fatalf("create: %d %s", w.Code, w.Body.String())
	}
	for _, owner := range []auth.User{user, {ID: 0}} {
		currentUser = owner
		w = httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("GET", "/api/private/products", nil))
		var listed struct {
			Data []Category `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &listed); err != nil {
			t.Fatal(err)
		}
		if w.Code != 200 {
			t.Fatalf("list: %d %s", w.Code, w.Body.String())
		}
		if owner.ID == 0 {
			if len(listed.Data) != 0 {
				t.Fatal("Product pengguna lain bocor")
			}
		} else if len(listed.Data) != 1 || listed.Data[0].User.ID != owner.ID || listed.Data[0].ID != created.Data.ID {
			t.Fatal("Relasi product dan user tidak dimuat")
		}
	}
	if err := db.Delete(&user).Error; err == nil {
		t.Fatal("Foreign key harus menolak penghapusan pembuat product")
	}
}
