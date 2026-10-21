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

func (s *Service) Create(ctx context.Context, name string, user auth.User, categoryID uint) (Product, error) {
	item := Product{Name: strings.TrimSpace(name), UserID: user.ID, CategoryID: categoryID}
	if err := s.repo.Create(ctx, &item); err != nil {
		return Product{}, err
	}
	item.User = user
	return item, nil
}

func (s *Service) List(ctx context.Context, userID uint) ([]Product, error) {
	return s.repo.ListByUserID(ctx, userID)
}

func (s *Service) Detail(ctx context.Context, userID uint, productID uint) (Product, error) {
	return s.repo.DetailProduct(ctx, userID, productID)
}

func (s *Service) Update(ctx context.Context, userID uint, productID uint, dataProduct UpdateProductRequest) (Product, error) {
	return s.repo.UpdateProduct(ctx, userID, productID, dataProduct)
}

func (s *Service) Delete(ctx context.Context, userID uint, productID uint) error {
	err := s.repo.DeleteProduct(ctx, userID, productID)
	return err
}
