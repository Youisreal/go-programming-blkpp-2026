package category

import "time"

type Category struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"nor null;size:100"`
	CategoryAt time.Time `json:"created_at"`
}

var DefaultCategories = []Category{
	{Name: "Elektronik"},
	{Name: "Fashion"},
	{Name: "Makanan"},
	{Name: "Olahraga"},
	{Name: "Rumah Tangga"},
	{Name: "Perawatan Diri"},
}
