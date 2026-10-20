package product

// ponytail: CRUD diteruskan ke repository; tambah aturan bisnis saat use case membutuhkannya.
type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() []Product {
	return s.repo.GetAll()
}

func (s *Service) GetByID(id int) (Product, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Create(p Product) Product {
	return s.repo.Create(p)
}

func (s *Service) Update(id int, p Product) (Product, error) {
	return s.repo.Update(id, p)
}

func (s *Service) Delete(id int) error {
	return s.repo.Delete(id)
}
