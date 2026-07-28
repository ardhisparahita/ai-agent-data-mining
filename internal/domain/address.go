package domain

import "time"

type Address struct {
	ID            uint64 `gorm:"primaryKey"`
	UserID        uint64
	RecipientName string
	Phone         string
	Address       string
	City          string
	PostalCode    string
	User          User `gorm:"foreignKey:UserID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (a *Address) TableName() string {
	return "addresses"
}
