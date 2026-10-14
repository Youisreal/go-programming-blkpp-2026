package config

import (
	"net/url"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:password@localhost/postgres")
	t.Setenv("ADDR", "")
	cfg, err := Load()
	if err != nil || cfg.Addr != "127.0.0.1:8080" {
		t.Fatalf("default config: %+v, err=%v", cfg, err)
	}
	t.Setenv("ADDR", "127.0.0.1:8081")
	cfg, err = Load()
	if err != nil || cfg.Addr != "127.0.0.1:8081" {
		t.Fatalf("custom config: %+v, err=%v", cfg, err)
	}
	t.Setenv("DATABASE_URL", "")
	if _, err := Load(); err == nil {
		t.Fatal("DATABASE_URL kosong harus ditolak")
	}
}

func TestDatabaseURL(t *testing.T) {
	for _, value := range []string{"", "not a URL", "https://example.com/db", "postgres://user@/db", "postgres://localhost/db", "postgres://user@localhost", "postgres://user@localhost/db?sslmode=disable", "postgres://user@localhost/db?sslmode=prefer"} {
		if _, err := DatabaseURL(value); err == nil {
			t.Fatalf("connection string tidak valid diterima: %q", value)
		}
	}
	for _, mode := range []string{"", "require", "verify-ca", "verify-full"} {
		dsn, err := DatabaseURL("postgresql://user:p%40ss@localhost:5432/postgres?sslmode=" + mode)
		if err != nil {
			t.Fatal(err)
		}
		u, err := url.Parse(dsn)
		if err != nil {
			t.Fatal(err)
		}
		want := mode
		if want == "" {
			want = "require"
		}
		password, _ := u.User.Password()
		if u.Query().Get("sslmode") != want || password != "p@ss" {
			t.Fatal("SSL/password berubah saat normalisasi")
		}
	}
}
