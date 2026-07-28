package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestRegisterSuccess(t *testing.T) {

	repo := new(MockUserRepository)

	service := AuthServiceImpl{
		Repo: repo,
	}

	req := request.RegisterRequest{
		Name:     "Ardhis",
		Email:    "ardhis@gmail.com",
		Password: "password123",
	}

	repo.On(
		"FindByEmail",
		mock.Anything,
		req.Email,
	).Return(
		nil,
		gorm.ErrRecordNotFound,
	)

	repo.On(
		"Create",
		mock.Anything,
		mock.AnythingOfType("*domain.User"),
	).Run(func(args mock.Arguments) {

		user := args.Get(1).(*domain.User)
		user.ID = 1

	}).Return(nil)

	result, err := service.Register(
		context.Background(),
		req,
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(
		t,
		uint64(1),
		result.ID,
	)

	assert.Equal(
		t,
		req.Name,
		result.Name,
	)

	assert.Equal(
		t,
		req.Email,
		result.Email,
	)

	assert.Equal(
		t,
		"CUSTOMER",
		result.Role,
	)

	repo.AssertExpectations(t)
}

func TestRegisterDuplicateEmail(t *testing.T) {

	repo := new(MockUserRepository)

	service := AuthServiceImpl{
		Repo: repo,
	}

	req := request.RegisterRequest{
		Name:     "Ardhis",
		Email:    "ardhis@gmail.com",
		Password: "password123",
	}

	repo.On(
		"FindByEmail",
		mock.Anything,
		req.Email,
	).Return(
		&domain.User{
			ID:    1,
			Email: req.Email,
		},
		nil,
	)

	result, err := service.Register(
		context.Background(),
		req,
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"email already registered",
	)

	repo.AssertExpectations(t)
}

func TestRegisterFindByEmailRepositoryError(t *testing.T) {

	repo := new(MockUserRepository)

	service := AuthServiceImpl{
		Repo: repo,
	}

	req := request.RegisterRequest{
		Name:     "Ardhis",
		Email:    "ardhis@gmail.com",
		Password: "password123",
	}

	repo.On(
		"FindByEmail",
		mock.Anything,
		req.Email,
	).Return(
		nil,
		errors.New("database error"),
	)

	result, err := service.Register(
		context.Background(),
		req,
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"database error",
	)

	repo.AssertExpectations(t)
}

func TestRegisterCreateRepositoryError(t *testing.T) {

	repo := new(MockUserRepository)

	service := AuthServiceImpl{
		Repo: repo,
	}

	req := request.RegisterRequest{
		Name:     "Ardhis",
		Email:    "ardhis@gmail.com",
		Password: "password123",
	}

	repo.On(
		"FindByEmail",
		mock.Anything,
		req.Email,
	).Return(
		nil,
		gorm.ErrRecordNotFound,
	)

	repo.On(
		"Create",
		mock.Anything,
		mock.AnythingOfType("*domain.User"),
	).Return(
		errors.New("database error"),
	)

	result, err := service.Register(
		context.Background(),
		req,
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"database error",
	)

	repo.AssertExpectations(t)
}