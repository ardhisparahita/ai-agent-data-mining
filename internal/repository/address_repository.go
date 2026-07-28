package repository

import (
	"context"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
)

type AddressRepository interface {
	Create(ctx context.Context, address *domain.Address) error
	FindAllByUserID(ctx context.Context, userID uint64) ([]domain.Address, error)
	FindByIDAndUserID(ctx context.Context, id uint64, userID uint64) (*domain.Address, error)
	Update(ctx context.Context, address *domain.Address) error
	Delete(ctx context.Context, id uint64) error
}
