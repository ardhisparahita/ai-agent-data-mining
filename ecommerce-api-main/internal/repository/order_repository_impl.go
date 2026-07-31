package repository

import (
	"context"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"gorm.io/gorm"
)

type OrderRepositoryImpl struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &OrderRepositoryImpl{
		DB: db,
	}
}

func (r *OrderRepositoryImpl) CreateTx(ctx context.Context, tx *gorm.DB, order *domain.Order) error {
	return tx.WithContext(ctx).Create(order).Error
}

func (r *OrderRepositoryImpl) FindAllByUserID(ctx context.Context, userID uint64) ([]domain.Order, error) {
	var orders []domain.Order

	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepositoryImpl) FindByIDAndUserID(ctx context.Context, id uint64, userID uint64) (*domain.Order, error) {
	var order domain.Order

	err := r.DB.WithContext(ctx).Preload("OrderItems").Preload("Payment").Where("id = ? AND user_id = ?", id, userID).First(&order).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepositoryImpl) UpdateTx(ctx context.Context, tx *gorm.DB, order *domain.Order) error {
	return tx.WithContext(ctx).Save(order).Error
}

func (r *OrderRepositoryImpl) FindByID(ctx context.Context, id uint64) (*domain.Order, error) {
	var order domain.Order
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepositoryImpl) Update(ctx context.Context, order *domain.Order) error {
	return r.DB.WithContext(ctx).Save(order).Error
}

func (r *OrderRepositoryImpl) FindByIDWithItems(ctx context.Context, id uint64) (*domain.Order, error) {
	var order domain.Order

	err := r.DB.WithContext(ctx).Preload("OrderItems").Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, err
}
func (r *OrderRepositoryImpl) FindByIDAndUserIDWithItems(ctx context.Context, id uint64, userID uint64) (*domain.Order, error) {
	var order domain.Order

	err := r.DB.WithContext(ctx).Preload("OrderItems").Where("id = ? AND user_id = ?", id, userID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, err
}
