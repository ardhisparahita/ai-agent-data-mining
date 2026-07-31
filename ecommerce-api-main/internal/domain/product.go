package domain

import "time"

type Product struct {
	ID          uint64 `gorm:"primaryKey"`
	CategoryID  uint64
	Name        string
	Description string
	Price       float64
	Stock       int
	ImageURL    string
	Category    Category `gorm:"foreignKey:categoryID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Product) TableName() string {
	return "products"
}
