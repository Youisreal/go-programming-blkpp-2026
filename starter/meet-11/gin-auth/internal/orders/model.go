package orders

import (
	"gin-auth/internal/auth"
	"gin-auth/internal/product"
	"time"
)

type Orders struct {
	ID        uint      `json:"order_id" gorm:"column:id;primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	User      auth.User `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Total     uint      `json:"total" gorm:"not null;index"`
}

type Order_detail struct {
	OrderID   uint            `json:"order_id" gorm:"index;primaryKey"`
	Order     Orders          `json:"order" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	ProductID uint            `json:"product_id" gorm:"index;primaryKey"`
	Product   product.Product `json:"product" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Qty       uint            `json:"qty" gorm:"not null;index"`
	Price     uint            `json:"price" gorm:"not null;index"`
}
