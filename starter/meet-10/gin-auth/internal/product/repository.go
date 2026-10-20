package product

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, item *Product) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *Repository) ListByUserID(ctx context.Context, userID uint) ([]Product, error) {
	items := []Product{}
	// ponytail: daftar maksimal 100 product terbaru; tambah pagination saat data melebihi batas.
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Preload("User").Order("id DESC").Limit(100).Find(&items).Error
	return items, err
}

func (r *Repository) DetailProduct(ctx context.Context, userID uint, productID uint) (Product, error) {
	var p Product
	err := r.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, productID).Preload("User").Find(&p).Error
	return p, err
}

func (r *Repository) UpdateProduct(ctx context.Context, userID uint, productID uint, updatedData UpdateProductRequest) (Product, error) {
	var p Product
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", productID, userID).First(&p).Error
	if err != nil {
		return p, err
	}
	err = r.db.WithContext(ctx).Model(&p).Updates(updatedData).Error
	return p, err
}

func (r *Repository) DeleteProduct(ctx context.Context, userID uint, productID uint) error {
	var p Product
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", productID, userID).Delete(&p).Error
	return err
}
