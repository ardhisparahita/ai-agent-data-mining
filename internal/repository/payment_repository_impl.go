package repository

import (
	"context"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"gorm.io/gorm"
)

type PaymentRepositoryImpl struct {
	DB *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &PaymentRepositoryImpl{
		DB: db,
	}
}

func (r *PaymentRepositoryImpl) CreateTx(ctx context.Context, tx *gorm.DB, payment *domain.Payment) error {
	return tx.WithContext(ctx).Create(payment).Error
}

func (r *PaymentRepositoryImpl) FIndByOrderID(ctx context.Context, orderID uint64) (*domain.Payment, error) {
	var payment domain.Payment

	err := r.DB.WithContext(ctx).Where("order_id = ?", orderID).First(&payment).Error
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *PaymentRepositoryImpl) UpdateTx(ctx context.Context, tx *gorm.DB, payment *domain.Payment) error {
	return tx.WithContext(ctx).Save(payment).Error
}
