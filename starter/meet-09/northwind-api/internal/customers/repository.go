package customers

import (
	"context"
	"database/sql"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (repo *Repository) GetAll(ctx context.Context) ([]Customer, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT customer_id, company_name, contact_name, contact_title from customers ORDER BY customer_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	values := make([]Customer, 0)
	for rows.Next() {
		var value Customer
		if err := rows.Scan(&value.CustomerID, &value.CompanyName, &value.CustomerName, &value.ContactTitle); err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, rows.Err()
}

func (repo *Repository) GetByID(ctx context.Context, id string) (Customer, error) {
	var value Customer
	err := repo.db.QueryRowContext(ctx, "SELECT customer_id, company_name, contact_name, contact_title FROM customers WHERE customer_id=$1", id).Scan(&value.CustomerID, &value.CompanyName, &value.CustomerName, &value.ContactTitle)
	return value, err
}

func (repo *Repository) Create(ctx context.Context, value Customer) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO customers (customer_id, company_name, contact_name, contact_title) VALUES ($1, $2, $3, $4)", value.CustomerID, value.CompanyName, value.CustomerName, value.ContactTitle)
	return err
}

func (repo *Repository) Update(ctx context.Context, value Customer) error {
	result, err := repo.db.ExecContext(ctx, "UPDATE customers SET company_name=$2, contact_name=$3, contact_title=$4 WHERE customer_id=$1", value.CustomerID, value.CompanyName, value.CustomerName, value.ContactTitle)
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

func (repo *Repository) Delete(ctx context.Context, id string) error {
	result, err := repo.db.ExecContext(ctx, "DELETE FROM customers WHERE customer_id=$1", id)
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
