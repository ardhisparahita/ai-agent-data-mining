package domain

import "time"

type Cart struct {
	ID        uint64 `gorm:"primaryKey"`
	UserID    uint64
	ProductID uint64
	Quantity  int
	User      User      `gorm:"foreignKey:UserID"`
	Product   Product   `gorm:"foreignKey:ProductID"`
	CreatedAt time.Time `gorm:"primaryKey"`
	UpdatedAt time.Time `gorm:"primaryKey"`
}
