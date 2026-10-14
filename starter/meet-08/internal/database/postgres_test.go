package database

import (
	"errors"
	"strings"
	"testing"
)

func TestConnectionError(t *testing.T) {
	dsn := "postgresql://postgres:p%40ss%3Aword@localhost/postgres?sslmode=require"
	err := connectionError(errors.New("password authentication failed; p@ss:word; p%40ss%3Aword; "+dsn), dsn)
	message := err.Error()
	if strings.Contains(message, "p@ss:word") || strings.Contains(message, "p%40ss%3Aword") || strings.Contains(message, dsn) {
		t.Fatal("error mengandung password atau connection string")
	}
	if !strings.Contains(message, "password authentication failed") {
		t.Fatal("penyebab error hilang")
	}
}
