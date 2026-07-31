package domain

import "time"

type Payment struct {
	ID      uint64 `gorm:"primaryKey"`
	OrderID uint64
	Method  string
	Status  string
	Amount  float64
	PaidAt  *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
