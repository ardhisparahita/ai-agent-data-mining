package service

import (
	"context"

	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
)

type CheckoutService interface {
	Checkout(ctx context.Context, userID uint64, req request.CheckoutRequest) (*response.OrderResponse, error)
}
