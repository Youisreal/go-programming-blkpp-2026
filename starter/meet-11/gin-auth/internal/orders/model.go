package orders

import (
	"time"

	"gin-auth/internal/auth"
	"gin-auth/internal/product"
)

type Order struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	User      auth.User `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CreatedAt time.Time `json:"created_at"`

	Total   int           `json:"total" gorm:"not null;default:0"`
	Details []OrderDetail `json:"order_details" gorm:"foreignKey:OrderID"`
}

type OrderDetail struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	OrderID   uint            `json:"order_id" gorm:"not null;index"`
	Order     *Order          `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	ProductID uint            `json:"product_id" gorm:"not null;index"`
	Product   product.Product `json:"product" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Qty       int             `json:"qty" gorm:"not null;default:0"`
	Price     int             `json:"price" gorm:"not null;default:0"`
}
