package service

import (
	"context"

	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
)

type OrderService interface {
	FindAll(ctx context.Context, userID uint64) ([]response.OrderListResponse, error)
	FindByID(ctx context.Context, id uint64, userID uint64) (*response.OrderDetailResponse, error)

	MarkAsPaid(ctx context.Context, id uint64) error
	MarkAsFailed(ctx context.Context, id uint64) error

	Cancel(ctx context.Context, id uint64, userID uint64) error
	MarkAsShipped(ctx context.Context, id uint64) error
	MarkAsCompleted(ctx context.Context, id uint64, userID uint64, role string) error
}
