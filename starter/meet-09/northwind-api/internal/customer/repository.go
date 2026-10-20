package customer

import (
	"context"
	"database/sql"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (repo *Repository) GetAll(ctx context.Context) ([]Customer, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT customer_id, company_name, contact_name, contact_title FROM customers ORDER BY customer_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make([]Customer, 0)
	for rows.Next() {
		var value Customer
		err := rows.Scan(
			&value.CustomerID,
			&value.CustomerName,
			&value.ContactName,
			&value.ContactTitle,
		)

		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}
	return values, rows.Err()
}

func (repo *Repository) GetByID(ctx context.Context, id string) (Customer, error) {
	var value Customer
	err := repo.db.QueryRowContext(ctx, "SELECT customer_id, company_name, contact_name, contact_title FROM customers WHERE customer_id=$1", id).Scan(
		&value.CustomerID,
		&value.CustomerName,
		&value.ContactName,
		&value.ContactTitle,
	)
	return value, err
}

// func (repo *Repository) GetByID(ctx context.Context, id int) (Category, error) {
// 	var value Category
// 	err := repo.db.QueryRowContext(ctx, "SELECT category_id, category_name, description FROM public.categories WHERE category_id=$1", id).Scan(&value.CategoryID, &value.CategoryName, &value.Description)
// 	return value, err
// }

// func (repo *Repository) GetByID(ctx context.Context, id int) (Category, error) {
// 	var value Category
// 	err := repo.db.QueryRowContext(ctx, "SELECT category_id, category_name, description FROM public.categories WHERE category_id=$1", id).Scan(&value.CategoryID, &value.CategoryName, &value.Description)
// 	return value, err
// }

// func (repo *Repository) Create(ctx context.Context, value Category) error {
// 	_, err := repo.db.ExecContext(ctx, "INSERT INTO public.categories (category_id, category_name, description) VALUES ($1, $2, $3)", value.CategoryID, value.CategoryName, value.Description)
// 	return err
// }

// func (repo *Repository) Update(ctx context.Context, value Category) error {
// 	result, err := repo.db.ExecContext(ctx, "UPDATE public.categories SET category_name=$2, description=$3 WHERE category_id=$1", value.CategoryID, value.CategoryName, value.Description)
// 	if err != nil {
// 		return err
// 	}
// 	count, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}
// 	if count == 0 {
// 		return sql.ErrNoRows
// 	}
// 	return nil
// }

// func (repo *Repository) Delete(ctx context.Context, id int) error {
// 	result, err := repo.db.ExecContext(ctx, "DELETE FROM public.categories WHERE category_id=$1", id)
// 	if err != nil {
// 		return err
// 	}
// 	count, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}
// 	if count == 0 {
// 		return sql.ErrNoRows
// 	}
// 	return nil
// }
