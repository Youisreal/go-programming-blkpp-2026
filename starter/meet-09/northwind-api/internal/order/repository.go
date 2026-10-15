package order

import (
	"context"
	"database/sql"
	"math"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (repo *Repository) GetAll(ctx context.Context) ([]Order, error) {
	rows, err := repo.db.QueryContext(ctx, `SELECT o.order_id, o.customer_id, o.order_date,
 COALESCE((SELECT SUM(d.unit_price::double precision * d.quantity * (1-d.discount::double precision))
 FROM public.order_details d WHERE d.order_id=o.order_id), 0)
 FROM public.orders o ORDER BY o.order_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := make([]Order, 0)
	for rows.Next() {
		var value Order
		if err := rows.Scan(&value.OrderID, &value.CustomerID, &value.OrderDate, &value.Total); err != nil {
			return nil, err
		}
		orders = append(orders, value)
	}
	return orders, rows.Err()
}

func (repo *Repository) GetByID(ctx context.Context, id int) (Order, error) {
	var value Order
	err := repo.db.QueryRowContext(ctx, `SELECT order_id, customer_id, order_date FROM public.orders WHERE order_id=$1`, id).Scan(&value.OrderID, &value.CustomerID, &value.OrderDate)
	if err != nil {
		return value, err
	}
	rows, err := repo.db.QueryContext(ctx, `SELECT product_id, unit_price::double precision, quantity, discount::double precision FROM public.order_details WHERE order_id=$1 ORDER BY product_id`, id)
	if err != nil {
		return value, err
	}
	defer rows.Close()
	value.Items = make([]Item, 0)
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ProductID, &item.UnitPrice, &item.Quantity, &item.Discount); err != nil {
			return value, err
		}
		item.Subtotal = Subtotal(item.UnitPrice, item.Quantity, item.Discount)
		value.Total += item.Subtotal
		value.Items = append(value.Items, item)
	}
	return value, rows.Err()
}

func (repo *Repository) CustomerExists(ctx context.Context, tx *sql.Tx, id string) error {
	var customerID string
	return tx.QueryRowContext(ctx, `SELECT customer_id FROM public.customers WHERE customer_id=$1 FOR KEY SHARE`, id).Scan(&customerID)
}

func (repo *Repository) LockProduct(ctx context.Context, tx *sql.Tx, id int) (float64, int, int, error) {
	var price sql.NullFloat64
	var stock sql.NullInt64
	var discontinued int
	err := tx.QueryRowContext(ctx, `SELECT unit_price::double precision, units_in_stock, discontinued FROM public.products WHERE product_id=$1 FOR UPDATE`, id).Scan(&price, &stock, &discontinued)
	if err != nil {
		return 0, 0, 0, err
	}
	if !price.Valid || price.Float64 < 0 || math.IsNaN(price.Float64) || math.IsInf(price.Float64, 0) {
		return 0, 0, 0, ErrPrice
	}
	return price.Float64, int(stock.Int64), discontinued, nil
}

func (repo *Repository) InsertOrder(ctx context.Context, tx *sql.Tx, input CreateRequest) (Order, error) {
	var value Order
	err := tx.QueryRowContext(ctx, `INSERT INTO public.orders (order_id, customer_id, order_date) VALUES ($1, $2, CURRENT_TIMESTAMP) RETURNING order_id, customer_id, order_date`, input.OrderID, input.CustomerID).Scan(&value.OrderID, &value.CustomerID, &value.OrderDate)
	return value, err
}

func (repo *Repository) InsertItem(ctx context.Context, tx *sql.Tx, orderID int, item Item) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO public.order_details (order_id, product_id, unit_price, quantity, discount) VALUES ($1, $2, $3, $4, $5)`, orderID, item.ProductID, item.UnitPrice, item.Quantity, item.Discount)
	return err
}

func (repo *Repository) ReduceStock(ctx context.Context, tx *sql.Tx, item Item) error {
	result, err := tx.ExecContext(ctx, `UPDATE public.products SET units_in_stock=units_in_stock-$2 WHERE product_id=$1 AND units_in_stock >= $2`, item.ProductID, item.Quantity)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return ErrStock
	}
	return nil
}
