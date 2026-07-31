package repository

import (
	"context"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateTx(ctx context.Context, tx *gorm.DB, order *domain.Order) error
	UpdateTx(ctx context.Context, tx *gorm.DB, order *domain.Order) error
	FindByID(ctx context.Context, id uint64) (*domain.Order, error)
	Update(ctx context.Context, order *domain.Order) error
	FindAllByUserID(ctx context.Context, userID uint64) ([]domain.Order, error)
	FindByIDAndUserID(ctx context.Context, id uint64, userID uint64) (*domain.Order, error)
	FindByIDWithItems(ctx context.Context, id uint64) (*domain.Order, error)
	FindByIDAndUserIDWithItems(ctx context.Context, id uint64, userID uint64) (*domain.Order, error)
}
