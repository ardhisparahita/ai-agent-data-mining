package domain

import "time"

type Category struct {
	ID        uint64 `gorm:"primaryKey"`
	Name      string `gorm:"size:100"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Category) TableName() string {
	return "categories"
}
