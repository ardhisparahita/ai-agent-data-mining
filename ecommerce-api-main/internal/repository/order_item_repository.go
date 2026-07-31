package repository

import (
	"context"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"gorm.io/gorm"
)

type OrderItemRepository interface {
	CreateTx(ctx context.Context, tx *gorm.DB, orderItem *domain.OrderItem) error
}
