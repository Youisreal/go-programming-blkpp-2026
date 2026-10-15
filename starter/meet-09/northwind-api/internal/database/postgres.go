package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func ConnectDatabase(dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL wajib diisi")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return db, nil
}
