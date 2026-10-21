package orders

import (
	"context"
	"errors"
	"gin-auth/internal/product"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, productID uint, quantity int, userID uint) (Order, error) {
	var order Order
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item product.Product
		// SELECT FOR UPDATE mengunci baris product sampai transaksi commit atau rollback.
		// Order bersamaan pada product yang sama harus menunggu, lalu membaca stok terbaru
		// agar stok yang sama tidak terjual dua kali.
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&item, productID).Error; err != nil {
			return err
		}

		if item.Stock < quantity {
			return ErrInsufficientStock
		}

		total, err := orderTotal(item.Price, quantity)
		if err != nil {
			return err
		}
		order = Order{UserID: userID, Total: total}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		detail := OrderDetail{OrderID: order.ID, ProductID: item.ID, Qty: quantity, Price: item.Price}
		if err := tx.Create(&detail).Error; err != nil {
			return err
		}

		if err := tx.Model(&item).Update("stock", gorm.Expr("stock - ?", quantity)).Error; err != nil {
			return err
		}

		return tx.Preload("User").Preload("Details.Product").First(&order, order.ID).Error
	})
	if err != nil {
		return Order{}, err
	}
	return order, nil
}

func (r *Repository) List(ctx context.Context, userID uint) ([]Order, error) {
	orders := []Order{}
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("User").Preload("Details.Product").Order("id DESC").Find(&orders).Error
	return orders, err
}

func (r *Repository) GetByID(ctx context.Context, id, userID uint) (*Order, error) {
	var order Order
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("User").Preload("Details.Product").First(&order, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &order, nil
}
