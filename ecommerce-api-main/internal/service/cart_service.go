package service

import (
	"context"

	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
)

type CartService interface {
	AddToCart(ctx context.Context, userID uint64, req request.AddToCartRequest) (*response.CartResponse, error)
	FindAll(ctx context.Context, userID uint64) (*response.CartListResponse, error)
	Update(ctx context.Context, id uint64, userID uint64, req request.UpdateCartRequest) (*response.CartResponse, error)
	Delete(ctx context.Context, id uint64, userID uint64) error
}
