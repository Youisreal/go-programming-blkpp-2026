package customers

import "context"

// ponytail: CRUD diteruskan ke repository untuk pola workshop yang konsisten;
// tambahkan aturan bisnis di sini saat use case membutuhkannya.
type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo: repo} }

func (s *Service) GetAll(ctx context.Context) ([]Customer, error) {
	return s.repo.GetAll(ctx)
}
func (s *Service) GetByID(ctx context.Context, id string) (Customer, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *Service) Create(ctx context.Context, value Customer) error {
	return s.repo.Create(ctx, value)
}

func (s *Service) Update(ctx context.Context, value Customer) error {
	return s.repo.Update(ctx, value)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
