package category

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

func (r *Repository) Create(ctx context.Context, item *Category) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *Repository) ListByUserID(ctx context.Context, userID uint) ([]Category, error) {
	items := []Category{}
	// ponytail: daftar maksimal 100 product terbaru; tambah pagination saat data melebihi batas.
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Preload("User").Order("id DESC").Limit(100).Find(&items).Error
	return items, err
}

func (r *Repository) DetailCategory(ctx context.Context, userID uint, categoryID uint) (Category, error) {
	var p Category
	err := r.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, categoryID).Preload("User").Find(&p).Error
	return p, err
}

func (r *Repository) UpdateCategory(ctx context.Context, userID uint, categoryID uint, updatedData UpdateCategoryRequest) (Category, error) {
	var p Category
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", categoryID, userID).First(&p).Error
	if err != nil {
		return p, err
	}
	err = r.db.WithContext(ctx).Model(&p).Updates(updatedData).Error
	return p, err
}

func (r *Repository) DeleteCategory(ctx context.Context, userID uint, categoryID uint) error {
	var p Category
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", categoryID, userID).Delete(&p).Error
	return err
}
