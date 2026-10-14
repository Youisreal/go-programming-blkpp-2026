package summary

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/lib/pq"
)

type Repository struct{ DB *sql.DB }

func (r Repository) query(ctx context.Context, query string, args ...any) (json.RawMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	// PostgreSQL membentuk JSON untuk SELECT maupun INSERT/UPDATE/DELETE RETURNING.
	query = "WITH result AS (" + query + ") SELECT COALESCE(json_agg(result), '[]'::json) FROM result"
	var output []byte
	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&output)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return nil, ErrForeignKey
			case "23505":
				return nil, ErrConflict
			case "23514":
				return nil, ErrConstraint
			}
		}
		return nil, err
	}
	if !json.Valid(output) {
		return nil, errors.New("PostgreSQL menghasilkan JSON tidak valid")
	}
	return output, nil
}

func (r Repository) Summary(ctx context.Context) (json.RawMessage, error) {
	return r.query(ctx, "SELECT COUNT(*) FROM products as product_count")
}
