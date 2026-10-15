package product

import (
	"context"
	"database/sql"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (repo *Repository) GetAll(ctx context.Context) ([]Product, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT product_id, product_name, supplier_id, category_id, COALESCE(unit_price, 0), COALESCE(units_in_stock, 0), discontinued FROM public.products ORDER BY product_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	values := make([]Product, 0)
	for rows.Next() {
		var value Product
		if err := rows.Scan(&value.ProductID, &value.ProductName, &value.SupplierID, &value.CategoryID, &value.UnitPrice, &value.UnitsInStock, &value.Discontinued); err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, rows.Err()
}

func (repo *Repository) GetByID(ctx context.Context, id int) (Product, error) {
	var value Product
	err := repo.db.QueryRowContext(ctx, "SELECT product_id, product_name, supplier_id, category_id, COALESCE(unit_price, 0), COALESCE(units_in_stock, 0), discontinued FROM public.products WHERE product_id=$1", id).Scan(&value.ProductID, &value.ProductName, &value.SupplierID, &value.CategoryID, &value.UnitPrice, &value.UnitsInStock, &value.Discontinued)
	return value, err
}

func (repo *Repository) Create(ctx context.Context, value Product) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO public.products (product_id, product_name, supplier_id, category_id, unit_price, units_in_stock, discontinued) VALUES ($1, $2, $3, $4, $5, $6, $7)", value.ProductID, value.ProductName, value.SupplierID, value.CategoryID, value.UnitPrice, value.UnitsInStock, value.Discontinued)
	return err
}

func (repo *Repository) Update(ctx context.Context, value Product) error {
	result, err := repo.db.ExecContext(ctx, "UPDATE public.products SET product_name=$2, supplier_id=$3, category_id=$4, unit_price=$5, units_in_stock=$6, discontinued=$7 WHERE product_id=$1", value.ProductID, value.ProductName, value.SupplierID, value.CategoryID, value.UnitPrice, value.UnitsInStock, value.Discontinued)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (repo *Repository) Delete(ctx context.Context, id int) error {
	result, err := repo.db.ExecContext(ctx, "DELETE FROM public.products WHERE product_id=$1", id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}
