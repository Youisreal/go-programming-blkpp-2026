package category

import (
	"time"

	"gin-auth/internal/auth"
)

type Category struct {
	ID        uint      `json:"category_id" gorm:"primaryKey"`
	Name      string    `json:"category_name" gorm:"not null;size:100"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	User      auth.User `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateCategoryRequest struct {
	Name string `json:"category_name" gorm:"not null;size:100"`
}
