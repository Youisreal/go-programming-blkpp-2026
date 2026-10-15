package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func Open(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, connectionError(err, dsn)
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(30 * time.Minute)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, connectionError(err, dsn)
	}
	return db, nil
}

func connectionError(err error, dsn string) error {
	message := strings.ReplaceAll(err.Error(), dsn, "[DATABASE_URL]")
	if u, parseErr := url.Parse(dsn); parseErr == nil && u.User != nil {
		password, _ := u.User.Password()
		_, encoded, _ := strings.Cut(u.User.String(), ":")
		for _, secret := range []string{password, encoded} {
			if secret != "" {
				message = strings.ReplaceAll(message, secret, "[PASSWORD]")
			}
		}
	}
	return fmt.Errorf("koneksi Supabase PostgreSQL gagal: %s", message)
}
