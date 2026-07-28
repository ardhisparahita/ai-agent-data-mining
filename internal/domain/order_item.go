package domain

type OrderItem struct {
	ID          uint64 `gorm:"primaryKey"`
	OrderID     uint64
	ProductID   uint64
	ProductName string
	Price       float64
	Quantity    int
	Subtotal    float64
	Product     Product `gorm:"foreignKey:ProductID"`
}

func (oi *OrderItem) TableName() string {
	return "order_items"
}
