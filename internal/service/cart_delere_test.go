package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCartDeleteSuccess(t *testing.T) {

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)

	service := CartServiceImpl{
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}

	cart := &domain.Cart{
		ID:        1,
		UserID:    1,
		ProductID: 1,
		Quantity:  2,
	}

	cartRepo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(
		cart,
		nil,
	)

	cartRepo.On(
		"Delete",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(nil)

	err := service.Delete(
		context.Background(),
		1,
		1,
	)

	assert.NoError(t, err)

	cartRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestCartDeleteNotFound(t *testing.T) {

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)

	service := CartServiceImpl{
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}

	cartRepo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(99),
		uint64(1),
	).Return(
		nil,
		gorm.ErrRecordNotFound,
	)

	err := service.Delete(
		context.Background(),
		99,
		1,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"cart not found",
	)

	cartRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestCartDeleteRepositoryError(t *testing.T) {

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)

	service := CartServiceImpl{
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}

	cart := &domain.Cart{
		ID:        1,
		UserID:    1,
		ProductID: 1,
		Quantity:  2,
	}

	cartRepo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(
		cart,
		nil,
	)

	cartRepo.On(
		"Delete",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(
		errors.New("database error"),
	)

	err := service.Delete(
		context.Background(),
		1,
		1,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"database error",
	)

	cartRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}
