package service

import (
	"context"

	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
)

type AuthService interface {
	Register(ctx context.Context, req request.RegisterRequest) (*response.UserResponse, error)
	Login(ctx context.Context, req request.LoginRequest) (*response.AuthResponse, error)
	GetProfile(ctx context.Context, userID uint64) (*response.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uint64, req request.UpdateProfileRequest) (*response.UserResponse, error)
	ChangePassword(ctx context.Context, UserID uint64, req request.ChangePasswordRequest) error
}
