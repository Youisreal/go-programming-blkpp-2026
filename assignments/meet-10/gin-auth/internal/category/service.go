package category

import (
	"context"
	"strings"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, name string) (Category, error) {
	item := Category{Name: strings.TrimSpace(name)}
	if err := s.repo.Create(ctx, &item); err != nil {
		return Category{}, err
	}
	return item, nil
}

func (s *Service) List(ctx context.Context, userID uint) ([]Category, error) {
	return s.repo.ListByUserID(ctx, userID)
}

func (s *Service) Detail(ctx context.Context, userID uint, productID uint) (Category, error) {
	return s.repo.DetailCategory(ctx, userID, productID)
}

func (s *Service) Update(ctx context.Context, userID uint, productID uint, dataCategory UpdateCategoryRequest) (Category, error) {
	return s.repo.UpdateCategory(ctx, userID, productID, dataCategory)
}

func (s *Service) Delete(ctx context.Context, userID uint, productID uint) error {
	err := s.repo.DeleteCategory(ctx, userID, productID)
	return err
}
