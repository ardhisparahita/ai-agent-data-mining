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

func TestCartUpdateSuccess(t *testing.T) {

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)

	service := CartServiceImpl{
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}

	product := domain.Product{
		ID:       1,
		Name:     "Keyboard",
		Price:    750000,
		Stock:    10,
		ImageURL: "keyboard.jpg",
	}

	cart := &domain.Cart{
		ID:        1,
		UserID:    1,
		ProductID: 1,
		Quantity:  2,
		Product:   product,
	}

	cartRepo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(cart, nil).Once()

	cartRepo.On(
		"Update",
		mock.Anything,
		cart,
	).Return(nil)

	updatedCart := &domain.Cart{
		ID:        1,
		UserID:    1,
		ProductID: 1,
		Quantity:  5,
		Product:   product,
	}

	cartRepo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(updatedCart, nil).Once()

	result, err := service.Update(
		context.Background(),
		1,
		1,
		request.UpdateCartRequest{
			Quantity: 5,
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, 5, result.Quantity)

	cartRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestCartUpdateNotFound(t *testing.T) {

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
	).Return(nil, gorm.ErrRecordNotFound)

	result, err := service.Update(
		context.Background(),
		99,
		1,
		request.UpdateCartRequest{
			Quantity: 5,
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"cart not found",
	)

	cartRepo.AssertExpectations(t)
}

func TestCartUpdateInvalidQuantity(t *testing.T) {

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)

	service := CartServiceImpl{
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}

	result, err := service.Update(
		context.Background(),
		1,
		1,
		request.UpdateCartRequest{
			Quantity: 0,
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"quantity must be greater than 0",
	)

	cartRepo.AssertNotCalled(
		t,
		"FindByIDAndUserID",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	)
}

func TestCartUpdateInsufficientStock(t *testing.T) {

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
		Product: domain.Product{
			ID:    1,
			Stock: 2,
		},
	}

	cartRepo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(cart, nil)

	result, err := service.Update(
		context.Background(),
		1,
		1,
		request.UpdateCartRequest{
			Quantity: 5,
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"insufficient stock",
	)

	cartRepo.AssertExpectations(t)
}

func TestCartUpdateRepositoryError(t *testing.T) {

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
		Product: domain.Product{
			ID:    1,
			Stock: 10,
		},
	}

	cartRepo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(cart, nil)

	cartRepo.On(
		"Update",
		mock.Anything,
		cart,
	).Return(
		errors.New("database error"),
	)

	result, err := service.Update(
		context.Background(),
		1,
		1,
		request.UpdateCartRequest{
			Quantity: 5,
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"database error",
	)

	cartRepo.AssertExpectations(t)
}
