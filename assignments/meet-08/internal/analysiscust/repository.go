package analysiscust

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

func (r Repository) Analysiscust(ctx context.Context) (json.RawMessage, error) {
	return r.query(ctx, "select c.customer_id, c.contact_name, c.country, SUM(od.quantity * od.unit_price * (1 - od.discount )) as total_revenue from customers c join orders o on o.customer_id = c.customer_id join order_details od on od.order_id = o.order_id group by c.customer_id, c.contact_name order by total_revenue desc")
}
