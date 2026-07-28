package repository

import (
	"context"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreateTx(ctx context.Context, tx *gorm.DB, payment *domain.Payment) error
	FIndByOrderID(ctx context.Context, orderID uint64) (*domain.Payment, error)
	UpdateTx(ctx context.Context, tx *gorm.DB, payment *domain.Payment) error
}
