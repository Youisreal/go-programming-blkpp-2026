package order

import (
	"context"
	"database/sql"
	"errors"
	"sort"
)

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo: repo} }

func (s *Service) GetAll(ctx context.Context) ([]Order, error) { return s.repo.GetAll(ctx) }

func (s *Service) GetByID(ctx context.Context, id int) (Order, error) { return s.repo.GetByID(ctx, id) }

func (s *Service) Create(ctx context.Context, input CreateRequest) (Order, error) {
	if err := Validate(input); err != nil {
		return Order{}, err
	}
	tx, err := s.repo.db.BeginTx(ctx, nil)
	if err != nil {
		return Order{}, err
	}
	defer tx.Rollback()
	if err := s.repo.CustomerExists(ctx, tx, input.CustomerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Order{}, ErrCustomerNotFound
		}
		return Order{}, err
	}
	// Lock selalu berurutan agar dua order dengan item terbalik tidak saling deadlock.
	items := append([]ItemRequest(nil), input.Items...)
	sort.Slice(items, func(i, j int) bool { return items[i].ProductID < items[j].ProductID })
	details := make([]Item, 0, len(items))
	for _, requested := range items {
		price, stock, discontinued, err := s.repo.LockProduct(ctx, tx, requested.ProductID)
		if errors.Is(err, sql.ErrNoRows) {
			return Order{}, ErrProductNotFound
		}
		if err != nil {
			return Order{}, err
		}
		if stock < requested.Quantity || discontinued != 0 {
			return Order{}, ErrStock
		}
		// Northwind menyimpan discount sebagai real; samakan perhitungan dengan nilai tersimpan.
		discount := float64(float32(requested.Discount))
		details = append(details, Item{
			ProductID: requested.ProductID,
			UnitPrice: price,
			Quantity:  requested.Quantity,
			Discount:  discount,
			Subtotal:  Subtotal(price, requested.Quantity, discount),
		})
	}
	value, err := s.repo.InsertOrder(ctx, tx, input)
	if err != nil {
		return Order{}, err
	}
	for _, item := range details {
		if err := s.repo.InsertItem(ctx, tx, value.OrderID, item); err != nil {
			return Order{}, err
		}
		if err := s.repo.ReduceStock(ctx, tx, item); err != nil {
			return Order{}, err
		}
		value.Total += item.Subtotal
	}
	value.Items = details
	if err := tx.Commit(); err != nil {
		return Order{}, err
	}
	return value, nil
}
