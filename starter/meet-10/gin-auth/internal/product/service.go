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

func (s *Service) Create(ctx context.Context, name string, user auth.User) (Product, error) {
	item := Product{Name: strings.TrimSpace(name), UserID: user.ID}
	if err := s.repo.Create(ctx, &item); err != nil {
		return Product{}, err
	}
	item.User = user
	return item, nil
}

func (s *Service) List(ctx context.Context, userID uint) ([]Product, error) {
	return s.repo.ListByUserID(ctx, userID)
}
