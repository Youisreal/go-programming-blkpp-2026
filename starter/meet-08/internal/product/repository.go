package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/lib/pq"
)

// ponytail: satu backend PostgreSQL; tambah interface jika ada backend kedua.
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

func (r Repository) List(ctx context.Context) (json.RawMessage, error) {
	return r.query(ctx, "select *, c.category_name from products p join categories c on p.category_id = c.category_id group by p.product_id , c.category_id order by p.product_id")
}
func (r Repository) Get(ctx context.Context, id int) (json.RawMessage, error) {
	return r.query(ctx, "SELECT * FROM products WHERE product_id = $1", id)
}
func (r Repository) Create(ctx context.Context, p Input) (json.RawMessage, error) {
	return r.query(ctx, `INSERT INTO products (product_name, supplier_id, category_id, quantity_per_unit,
        unit_price, units_in_stock, units_on_order, reorder_level, discontinued)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *`,
		p.ProductName, p.SupplierID, p.CategoryID, p.QuantityPerUnit, p.UnitPrice,
		p.UnitsInStock, p.UnitsOnOrder, p.ReorderLevel, p.Discontinued)
}
func (r Repository) Update(ctx context.Context, id int, p Input) (json.RawMessage, error) {
	return r.query(ctx, `UPDATE products SET product_name=$1, supplier_id=$2, category_id=$3,
        quantity_per_unit=$4, unit_price=$5, units_in_stock=$6, units_on_order=$7,
        reorder_level=$8, discontinued=$9 WHERE product_id=$10 RETURNING *`,
		p.ProductName, p.SupplierID, p.CategoryID, p.QuantityPerUnit, p.UnitPrice,
		p.UnitsInStock, p.UnitsOnOrder, p.ReorderLevel, p.Discontinued, id)
}
func (r Repository) Delete(ctx context.Context, id int) (json.RawMessage, error) {
	return r.query(ctx, "DELETE FROM products WHERE product_id = $1 RETURNING *", id)
}
func (r Repository) Categories(ctx context.Context) (json.RawMessage, error) {
	return r.query(ctx, "SELECT * FROM categories ORDER BY category_id")
}
func (r Repository) Suppliers(ctx context.Context) (json.RawMessage, error) {
	return r.query(ctx, "SELECT * FROM suppliers ORDER BY supplier_id")
}
