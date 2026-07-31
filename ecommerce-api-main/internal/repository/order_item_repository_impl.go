package repository

import (
	"context"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"gorm.io/gorm"
)

type OrderItemRepositoryImpl struct {
	DB *gorm.DB
}

func NewOrderItemRepository(db *gorm.DB) OrderItemRepository {
	return &OrderItemRepositoryImpl{
		DB: db,
	}
}

func (r *OrderItemRepositoryImpl) CreateTx(ctx context.Context, tx *gorm.DB, orderItem *domain.OrderItem) error {
	return tx.WithContext(ctx).Create(orderItem).Error
}
