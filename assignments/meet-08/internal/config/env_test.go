package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnvFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	content := "\ufeff# Windows UTF-8 BOM\r\n\r\nexport DATABASE_URL = 'postgres://user:p%40ss@localhost/db?sslmode=require'\r\nADDR=\"127.0.0.1:8081\"\r\nLITERAL=$(whoami) ${HOME} # literal\r\nEMPTY=\r\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	values, err := readEnvFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if values["ADDR"] != "127.0.0.1:8081" || values["LITERAL"] != "$(whoami) ${HOME} # literal" || values["EMPTY"] != "" {
		t.Fatal("nilai literal/CRLF/kutip salah")
	}
	if _, err := DatabaseURL(values["DATABASE_URL"]); err != nil {
		t.Fatal("connection string berubah saat membaca .env")
	}
	for _, content := range []string{"password-secret", "1KEY=password-secret", "KEY='password-secret", "KEY=password-secret\x00", "KEY=\xff"} {
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := readEnvFile(path); err == nil || strings.Contains(err.Error(), "password-secret") {
			t.Fatalf("input rusak harus ditolak tanpa membocorkan nilai: %v", err)
		}
	}
	if _, err := readEnvFile(filepath.Join(t.TempDir(), "missing.env")); err != nil {
		t.Fatal(".env opsional jika memakai environment OS")
	}
}

func TestLoadDotEnv(t *testing.T) {
	dir := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(cwd); err != nil {
			t.Error(err)
		}
	})
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"DATABASE_URL", "ADDR"} {
		t.Setenv(key, "")
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(".env", []byte("DATABASE_URL='postgres://file:password@localhost/file_db'\nADDR=127.0.0.1:8081\n"), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load()
	if err != nil || cfg.Addr != "127.0.0.1:8081" || !strings.Contains(cfg.DatabaseURL, "/file_db") {
		t.Fatal("config .env tidak dimuat")
	}
	t.Setenv("DATABASE_URL", "postgres://os:password@localhost/os_db")
	t.Setenv("ADDR", "127.0.0.1:8082")
	cfg, err = Load()
	if err != nil || cfg.Addr != "127.0.0.1:8082" || !strings.Contains(cfg.DatabaseURL, "/os_db") {
		t.Fatal("environment OS harus lebih diutamakan daripada .env")
	}
	t.Setenv("DATABASE_URL", "")
	if _, err := Load(); err == nil {
		t.Fatal("environment OS kosong tidak boleh ditimpa dari .env")
	}
}
