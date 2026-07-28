package repository

import (
	"context"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"gorm.io/gorm"
)

type CartRepositoryImpl struct {
	DB *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &CartRepositoryImpl{
		DB: db,
	}
}
func (r *CartRepositoryImpl) Create(ctx context.Context, cart *domain.Cart) error {
	return r.DB.WithContext(ctx).Create(cart).Error
}

func (r *CartRepositoryImpl) FindAllByUserID(ctx context.Context, userID uint64) ([]domain.Cart, error) {
	var carts []domain.Cart

	err := r.DB.WithContext(ctx).Preload("Product").Where("user_id = ?", userID).Find(&carts).Error
	if err != nil {
		return nil, err
	}

	return carts, err
}

func (r *CartRepositoryImpl) FindByIDAndUserID(ctx context.Context, id uint64, userID uint64) (*domain.Cart, error) {
	var cart domain.Cart

	err := r.DB.WithContext(ctx).Preload("Product").Where("id = ? AND user_id = ?", id, userID).First(&cart).Error
	if err != nil {
		return nil, err
	}

	return &cart, err
}

func (r *CartRepositoryImpl) FindByUserIDAndProductID(ctx context.Context, userID uint64, productID uint64) (*domain.Cart, error) {
	var cart domain.Cart

	err := r.DB.WithContext(ctx).Where("user_id = ? AND product_id = ?", userID, productID).First(&cart).Error
	if err != nil {
		return nil, err
	}

	return &cart, err
}

func (r *CartRepositoryImpl) Update(ctx context.Context, cart *domain.Cart) error {
	return r.DB.WithContext(ctx).Save(cart).Error
}

func (r *CartRepositoryImpl) Delete(ctx context.Context, id uint64, userID uint64) error {
	return r.DB.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domain.Cart{}).Error
}

func (r *CartRepositoryImpl) DeleteAllByUserIDTx(ctx context.Context, tx *gorm.DB, userID uint64) error {
	return tx.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.Cart{}).Error
}

func (r *CartRepositoryImpl) DeleteAllByUserID(ctx context.Context, userID uint64) error {
	return r.DB.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.Cart{}).Error
}
