package repository

import (
	"context"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uint64) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}
