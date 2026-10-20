package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"gin-auth/internal/auth"
	"gin-auth/internal/database"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func TestPublicRoutes(t *testing.T) {
	router := setupEngine(nil, []byte("test-secret-with-at-least-32-bytes"))
	for _, tc := range []struct {
		method, path, body, hook string
		status                   int
	}{
		{"GET", "/health", "", "health", http.StatusOK},
		{"POST", "/api/auth/register", `{}`, "auth.register", http.StatusBadRequest},
		{"POST", "/api/auth/login", `{}`, "auth.login", http.StatusBadRequest},
		{"GET", "/api/me", "", "auth.me", http.StatusUnauthorized},
	} {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		if w.Code != tc.status || w.Header().Get("X-Mock-Middleware") != tc.hook {
			t.Fatalf("%s %s: status=%d, hook=%q", tc.method, tc.path, w.Code, w.Header().Get("X-Mock-Middleware"))
		}
	}
}

func TestAuthFlow(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("Isi TEST_DATABASE_URL dengan URL PostgreSQL khusus test integrasi")
	}
	secret := []byte("test-secret-with-at-least-32-bytes")
	db, err := database.Open(dsn)
	if err != nil {
		t.Fatal(err)
	}
	conn, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	if err := db.AutoMigrate(&auth.User{}); err != nil {
		t.Fatal(err)
	}
	email := fmt.Sprintf("test-%d@example.com", time.Now().UnixNano())
	t.Cleanup(func() {
		if err := db.Where("email = ?", email).Delete(&auth.User{}).Error; err != nil {
			t.Errorf("Gagal membersihkan pengguna test: %v", err)
		}
	})
	router := setupEngine(db, secret)
	request := func(t *testing.T, method, path, body, authorization string, status int, hook string) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		if authorization != "" {
			req.Header.Set("Authorization", authorization)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != status {
			t.Fatalf("%s %s: status %d, ingin %d: %s", method, path, w.Code, status, w.Body.String())
		}
		if got := w.Header().Get("X-Mock-Middleware"); got != hook {
			t.Fatalf("%s: hook %q, ingin %q", path, got, hook)
		}
		if strings.Contains(w.Body.String(), "password") && status < 400 {
			t.Fatalf("Response membocorkan password: %s", w.Body.String())
		}
		return w
	}

	request(t, "GET", "/health", "", "", http.StatusOK, "health")
	for _, body := range []string{
		`{`,
		`{"name":"   ","email":"user@example.com","password":"password123"}`,
		`{"name":"User","email":"invalid","password":"password123"}`,
		`{"name":"User","email":"user@example.com","password":"short"}`,
		`{"name":"User","email":"user@example.com","password":"` + strings.Repeat("é", 40) + `"}`,
	} {
		request(t, "POST", "/api/auth/register", body, "", http.StatusBadRequest, "auth.register")
	}
	registration := fmt.Sprintf(`{"name":" User ","email":%q,"password":"password123"}`, strings.ToUpper(email))
	request(t, "POST", "/api/auth/register", registration, "", http.StatusCreated, "auth.register")
	request(t, "POST", "/api/auth/register", strings.ReplaceAll(registration, strings.ToUpper(email), email), "", http.StatusConflict, "auth.register")
	var user auth.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		t.Fatal(err)
	}
	if user.Name != "User" || user.Email != email || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("password123")) != nil {
		t.Fatal("Data pengguna atau hash password salah")
	}
	request(t, "POST", "/api/auth/login", fmt.Sprintf(`{"email":%q,"password":"wrongpassword"}`, email), "", http.StatusUnauthorized, "auth.login")
	request(t, "POST", "/api/auth/login", fmt.Sprintf(`{"email":%q,"password":"password123"}`, "missing-"+email), "", http.StatusUnauthorized, "auth.login")
	w := request(t, "POST", "/api/auth/login", fmt.Sprintf(`{"email":%q,"password":"password123"}`, strings.ToUpper(email)), "", http.StatusOK, "auth.login")
	var login struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &login); err != nil || login.Data.AccessToken == "" {
		t.Fatalf("Response login salah: %s", w.Body.String())
	}
	w = request(t, "GET", "/api/me", "", "Bearer "+login.Data.AccessToken, http.StatusOK, "auth.me")
	var profile struct {
		Data auth.User `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &profile); err != nil || profile.Data.ID != user.ID || profile.Data.Email != user.Email {
		t.Fatalf("Profil salah: %s", w.Body.String())
	}
	for _, header := range []string{"", "Bearer", "Basic abc", "Bearer invalid", "Bearer abc def"} {
		request(t, "GET", "/api/me", "", header, http.StatusUnauthorized, "auth.me")
	}

	valid := jwt.RegisteredClaims{
		Issuer: "gin-auth", Subject: strconv.FormatUint(uint64(user.ID), 10), ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	for _, name := range []string{"expired", "no expiry", "wrong issuer", "invalid subject", "missing user", "wrong signature", "wrong algorithm"} {
		t.Run(name, func(t *testing.T) {
			claims := valid
			key := secret
			var method jwt.SigningMethod = jwt.SigningMethodHS256
			switch name {
			case "expired":
				claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Hour))
			case "no expiry":
				claims.ExpiresAt = nil
			case "wrong issuer":
				claims.Issuer = "other-app"
			case "invalid subject":
				claims.Subject = "-1"
			case "missing user":
				claims.Subject = "9223372036854775807"
			case "wrong signature":
				key = []byte("other-secret-with-at-least-32-bytes")
			case "wrong algorithm":
				method = jwt.SigningMethodHS384
			}
			signed, err := jwt.NewWithClaims(method, claims).SignedString(key)
			if err != nil {
				t.Fatal(err)
			}
			request(t, "GET", "/api/me", "", "Bearer "+signed, http.StatusUnauthorized, "auth.me")
		})
	}

	// Koneksi terpisah memastikan data tersimpan di PostgreSQL.
	persisted, err := database.Open(dsn)
	if err != nil {
		t.Fatal(err)
	}
	persistedConn, err := persisted.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer persistedConn.Close()
	var saved auth.User
	if err := persisted.First(&saved, user.ID).Error; err != nil || saved.Email != user.Email {
		t.Fatalf("Persistensi PostgreSQL gagal: %v", err)
	}
	if err := db.Delete(&user).Error; err != nil {
		t.Fatal(err)
	}
	request(t, "GET", "/api/me", "", "Bearer "+login.Data.AccessToken, http.StatusUnauthorized, "auth.me")
}
