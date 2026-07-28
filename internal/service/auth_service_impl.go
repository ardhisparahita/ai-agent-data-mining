package service

import (
	"context"
	"errors"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
	"github.com/ardhisparahita/ecommerce-api/internal/repository"
	"github.com/ardhisparahita/ecommerce-api/pkg/config"
	"github.com/ardhisparahita/ecommerce-api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthServiceImpl struct {
	Repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &AuthServiceImpl{
		Repo: repo,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, req request.RegisterRequest) (*response.UserResponse, error) {
	_, err := s.Repo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, utils.Conflict("email already registered")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user := domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "CUSTOMER",
	}

	if err := s.Repo.Create(ctx, &user); err != nil {
		return nil, err
	}

	return &response.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, req request.LoginRequest) (*response.AuthResponse, error) {
	user, err := s.Repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, utils.Unauthorized("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		return nil, utils.Unauthorized("invalid email or password")
	}

	token, err := utils.GenerateToken(
		user.ID,
		user.Role,
		config.Get("JWT_SECRET"),
	)
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{Token: token}, nil
}

func (s *AuthServiceImpl) GetProfile(ctx context.Context, userID uint64) (*response.UserResponse, error) {
	user, err := s.Repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("user not found")
		}
		return nil, err
	}

	return &response.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *AuthServiceImpl) UpdateProfile(ctx context.Context, userID uint64, req request.UpdateProfileRequest) (*response.UserResponse, error) {
	user, err := s.Repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("user not found")
		}
		return nil, err
	}

	user.Name = req.Name
	if err := s.Repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return &response.UserResponse{
		ID:    user.ID,
		Name:  req.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *AuthServiceImpl) ChangePassword(ctx context.Context, UserID uint64, req request.ChangePasswordRequest) error {
	user, err := s.Repo.FindByID(ctx, UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("user not found")
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.OldPassword),
	)

	if err != nil {
		return utils.BadRequest("old password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return s.Repo.Update(ctx, user)
}
