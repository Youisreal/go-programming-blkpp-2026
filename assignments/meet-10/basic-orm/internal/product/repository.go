package product

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("product not found")

// ponytail: in-memory map cukup demo. Upgrade ke SQL DB saat butuh persistensi disk.
type Repository struct {
	mu     sync.RWMutex
	items  map[int]Product
	nextID int
}

func NewRepository() *Repository {
	return &Repository{
		items:  make(map[int]Product),
		nextID: 1,
	}
}

func (r *Repository) GetAll() []Product {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]Product, 0, len(r.items))
	for _, v := range r.items {
		list = append(list, v)
	}
	return list
}

func (r *Repository) GetByID(id int) (Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.items[id]
	if !ok {
		return Product{}, ErrNotFound
	}
	return p, nil
}

func (r *Repository) Create(p Product) Product {
	r.mu.Lock()
	defer r.mu.Unlock()
	p.ID = r.nextID
	r.items[r.nextID] = p
	r.nextID++
	return p
}

func (r *Repository) Update(id int, p Product) (Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[id]; !ok {
		return Product{}, ErrNotFound
	}
	p.ID = id
	r.items[id] = p
	return p, nil
}

func (r *Repository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[id]; !ok {
		return ErrNotFound
	}
	delete(r.items, id)
	return nil
}
