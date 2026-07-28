package mapper

import (
	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
)

func ToOrderListResponse(orders []domain.Order) []response.OrderListResponse {
	result := make([]response.OrderListResponse, 0, len(orders))

	for _, order := range orders {
		result = append(result, response.OrderListResponse{
			ID:          order.ID,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt,
		})
	}
	return result
}

func ToOrderDetailResponse(order *domain.Order) *response.OrderDetailResponse {
	items := make([]response.OrderItemResponse, 0, len(order.OrderItems))

	for _, item := range order.OrderItems {
		items = append(items, response.OrderItemResponse{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Price:       item.Price,
			Quantity:    item.Quantity,
			Subtotal:    item.Subtotal,
		})
	}

	return &response.OrderDetailResponse{
		ID:            order.ID,
		RecipientName: order.RecipientName,
		Phone:         order.Phone,
		Address:       order.Address,
		City:          order.City,
		PostalCode:    order.PostalCode,
		TotalAmount:   order.TotalAmount,
		Status:        order.Status,
		Items:         items,
		Payment: response.PaymentDetailResponse{
			Method: order.Payment.Method,
			Status: order.Payment.Status,
			Amount: order.Payment.Amount,
		},
	}
}
