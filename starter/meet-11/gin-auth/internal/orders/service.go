package orders

import (
	"context"
	"errors"
	"math"

	"gin-auth/internal/auth"
)

var (
	ErrInvalidOrder      = errors.New("product_id dan quantity harus lebih dari 0")
	ErrInsufficientStock = errors.New("stok product tidak cukup")
)

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo: repo} }

func (s *Service) Create(ctx context.Context, productID uint, quantity int, user auth.User) (Order, error) {
	if productID == 0 || quantity <= 0 || user.ID == 0 {
		return Order{}, ErrInvalidOrder
	}
	return s.repo.Create(ctx, productID, quantity, user.ID)
}

func orderTotal(price, quantity int) (int, error) {
	if price < 0 || quantity <= 0 || (price > 0 && quantity > math.MaxInt/price) {
		return 0, ErrInvalidOrder
	}
	return price * quantity, nil
}

func (s *Service) List(ctx context.Context, userID uint) ([]Order, error) {
	return s.repo.List(ctx, userID)
}
func (s *Service) GetByID(ctx context.Context, id, userID uint) (*Order, error) {
	return s.repo.GetByID(ctx, id, userID)
}
