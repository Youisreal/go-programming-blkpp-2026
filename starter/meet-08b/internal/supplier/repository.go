package supplier

import (
	"context"
	"database/sql"
	"encoding/json"
)

// Repository handles supplier persistence.
type Repository struct{ DB *sql.DB }

func (r Repository) List(ctx context.Context) (json.RawMessage, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT * FROM suppliers ORDER BY supplier_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Supplier, 0)
	for rows.Next() {
		var s Supplier
		if err := rows.Scan(&s.SupplierID, &s.SupplierName, &s.ContactName, &s.Address, &s.City, &s.PostalCode, &s.Country, &s.Phone); err != nil {
			return nil, err
		}
		items = append(items, s)
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
