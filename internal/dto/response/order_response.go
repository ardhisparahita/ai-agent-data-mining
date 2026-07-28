package response

import "time"

type OrderResponse struct {
	ID            uint64  `json:"id"`
	TotalAmount   float64 `json:"total_amount"`
	Status        string  `json:"status"`
	PaymentMethod string  `json:"payment_method"`
}

type OrderListResponse struct {
	ID          uint64    `json:"id"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrderDetailResponse struct {
	ID            uint64                `json:"id"`
	RecipientName string                `json:"recipient_name"`
	Phone         string                `json:"phone"`
	Address       string                `json:"address"`
	City          string                `json:"city"`
	PostalCode    string                `json:"postal_code"`
	TotalAmount   float64               `json:"total_amount"`
	Status        string                `json:"status"`
	Items         []OrderItemResponse   `json:"items"`
	Payment       PaymentDetailResponse `json:"payments"`
}

type OrderItemResponse struct {
	ProductID   uint64  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
}

type PaymentDetailResponse struct {
	Method string  `json:"method"`
	Status string  `json:"status"`
	Amount float64 `json:"amount"`
}
