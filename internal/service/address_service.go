package service

import (
	"context"

	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
)

type AddressService interface {
	Create(ctx context.Context, userID uint64, req request.CreateAddressRequest) (*response.AddressResponse, error)
	FindAllByUserID(ctx context.Context, userID uint64) ([]response.AddressResponse, error)
	FindByID(ctx context.Context, id uint64, userID uint64) (*response.AddressResponse, error)
	Update(ctx context.Context, id uint64, userID uint64, req request.UpdateAddressRequest) (*response.AddressResponse, error)
	Delete(ctx context.Context, id uint64, userID uint64) error
}
