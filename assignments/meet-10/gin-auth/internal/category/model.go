package category

import (
	"time"
)

type Category struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"nor null;size:100"`
	CategoryAt time.Time `json:"created_at"`
}

type UpdateCategoryRequest struct {
	Name string `json:"category_name" gorm:"not null;size:100"`
}
