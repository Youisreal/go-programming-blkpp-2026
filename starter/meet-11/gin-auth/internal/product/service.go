package product

import (
	"context"
	"gin-auth/internal/auth"
	"strings"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, name string, stock int, price int, user auth.User) (Product, error) {
<<<<<<< HEAD
	item := Product{Name: strings.TrimSpace(name), Stock: stock, Price: price, UserID: user.ID}
=======
	item := Product{
		Name:   strings.TrimSpace(name),
		Stock:  stock,
		Price:  price,
		UserID: user.ID,
	}
>>>>>>> efe0f2aefaad664f682bc23304654db574d83c45
	if err := s.repo.Create(ctx, &item); err != nil {
		return Product{}, err
	}
	item.User = user
	return item, nil
}

func (s *Service) List(ctx context.Context, userID uint) ([]Product, error) {
	return s.repo.ListByUserID(ctx, userID)
}
