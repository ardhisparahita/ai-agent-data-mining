package repository

import (
	"context"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		DB: db,
	}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id uint64) (*domain.User, error) {
	var user domain.User

	err := r.DB.WithContext(ctx).First(&user, id).Error

	if err != nil {
		return nil, err
	}

	return &user, err
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *domain.User) error {
	return r.DB.WithContext(ctx).Save(user).Error
}
