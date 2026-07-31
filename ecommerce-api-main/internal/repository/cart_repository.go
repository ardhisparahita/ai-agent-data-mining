package repository

import (
	"context"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"gorm.io/gorm"
)

type CartRepository interface {
	Create(ctx context.Context, cart *domain.Cart) error
	FindAllByUserID(ctx context.Context, userID uint64) ([]domain.Cart, error)
	FindByIDAndUserID(ctx context.Context, id uint64, userID uint64) (*domain.Cart, error)
	FindByUserIDAndProductID(ctx context.Context, userID uint64, productID uint64) (*domain.Cart, error)
	Update(ctx context.Context, cart *domain.Cart) error
	Delete(ctx context.Context, id uint64, userID uint64) error
	DeleteAllByUserID(ctx context.Context, userID uint64) error
	DeleteAllByUserIDTx(ctx context.Context, tx *gorm.DB, userID uint64) error
}
