package category

import (
	"context"
	"database/sql"
	"encoding/json"
)

// Repository handles category persistence.
type Repository struct{ DB *sql.DB }

func (r Repository) List(ctx context.Context) (json.RawMessage, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT category_id, category_name, description
		FROM categories
		ORDER BY category_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Category, 0)
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.CategoryID, &c.CategoryName, &c.Description); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	payload, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(payload), nil
}
