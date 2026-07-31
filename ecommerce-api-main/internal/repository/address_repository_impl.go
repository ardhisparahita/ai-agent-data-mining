package repository

import (
	"context"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"gorm.io/gorm"
)

type AddressRepositoryImpl struct {
	DB *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &AddressRepositoryImpl{
		DB: db,
	}
}

func (r *AddressRepositoryImpl) Create(ctx context.Context, address *domain.Address) error {
	return r.DB.WithContext(ctx).Create(address).Error
}

func (r *AddressRepositoryImpl) FindAllByUserID(ctx context.Context, userID uint64) ([]domain.Address, error) {
	var addresses []domain.Address

	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&addresses).Error

	return addresses, err
}

func (r *AddressRepositoryImpl) FindByIDAndUserID(ctx context.Context, id uint64, userID uint64) (*domain.Address, error) {
	var address domain.Address

	err := r.DB.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&address).Error
	if err != nil {
		return nil, err
	}

	return &address, nil
}

func (r *AddressRepositoryImpl) Update(ctx context.Context, address *domain.Address) error {
	return r.DB.WithContext(ctx).Save(address).Error
}

func (r *AddressRepositoryImpl) Delete(ctx context.Context, id uint64) error {
	return r.DB.WithContext(ctx).Delete(&domain.Address{}, id).Error
}
