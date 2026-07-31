package domain

import "time"

type Order struct {
	ID            uint64 `gorm:"primaryKey"`
	UserID        uint64
	AddressID     uint64
	RecipientName string
	Phone         string
	Address       string
	City          string
	PostalCode    string
	TotalAmount   float64
	Status        string

	User        User        `gorm:"foreignKey:UserID"`
	AddressData Address     `gorm:"foreignKey:AddressID"`
	OrderItems  []OrderItem `gorm:"foreignKey:OrderID"`
	Payment     Payment     `gorm:"foreignKey:OrderID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (o *Order) TableName() string {
	return "orders"
}
