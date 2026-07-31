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

func TestCartAddSuccessNewCart(t *testing.T) {

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)

	service := CartServiceImpl{
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}

	product := &domain.Product{
		ID:    1,
		Name:  "Keyboard",
		Price: 750000,
		Stock: 10,
	}

	productRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(product, nil)

	cartRepo.On(
		"FindByUserIDAndProductID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(nil, gorm.ErrRecordNotFound)

	cartRepo.On(
		"Create",
		mock.Anything,
		mock.Anything,
	).Run(func(args mock.Arguments) {

		cart := args.Get(1).(*domain.Cart)
		cart.ID = 1

	}).Return(nil)

	cartRepo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(
		&domain.Cart{
			ID:        1,
			UserID:    1,
			ProductID: 1,
			Quantity:  2,
			Product:   *product,
		},
		nil,
	)

	result, err := service.AddToCart(
		context.Background(),
		1,
		request.AddToCartRequest{
			ProductID: 1,
			Quantity:  2,
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, uint64(1), result.ID)
	assert.Equal(t, 2, result.Quantity)

	cartRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestCartAddSuccessExistingCart(t *testing.T) {

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)

	service := CartServiceImpl{
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}

	product := &domain.Product{
		ID:    1,
		Price: 750000,
		Stock: 10,
	}

	cart := &domain.Cart{
		ID:        1,
		UserID:    1,
		ProductID: 1,
		Quantity:  2,
		Product:   *product,
	}

	productRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(product, nil)

	cartRepo.On(
		"FindByUserIDAndProductID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(cart, nil)

	cartRepo.On(
		"Update",
		mock.Anything,
		cart,
	).Return(nil)

	cartRepo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(cart, nil)

	result, err := service.AddToCart(
		context.Background(),
		1,
		request.AddToCartRequest{
			ProductID: 1,
			Quantity:  3,
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, 5, result.Quantity)

	cartRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestCartAddProductNotFound(t *testing.T) {

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)

	service := CartServiceImpl{
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}

	productRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(
		nil,
		gorm.ErrRecordNotFound,
	)

	result, err := service.AddToCart(
		context.Background(),
		1,
		request.AddToCartRequest{
			ProductID: 1,
			Quantity:  2,
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"product not found",
	)

	productRepo.AssertExpectations(t)
}

func TestCartAddInsufficientStock(t *testing.T) {

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)

	service := CartServiceImpl{
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}

	product := &domain.Product{
		ID:    1,
		Price: 750000,
		Stock: 2,
	}

	productRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(product, nil)

	result, err := service.AddToCart(
		context.Background(),
		1,
		request.AddToCartRequest{
			ProductID: 1,
			Quantity:  5,
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"insufficient stock",
	)

	productRepo.AssertExpectations(t)
	cartRepo.AssertNotCalled(
		t,
		"FindByUserIDAndProductID",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	)
}

func TestCartAddCreateRepositoryError(t *testing.T) {

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)

	service := CartServiceImpl{
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}

	product := &domain.Product{
		ID:    1,
		Price: 750000,
		Stock: 10,
	}

	productRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(product, nil)

	cartRepo.On(
		"FindByUserIDAndProductID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(nil, gorm.ErrRecordNotFound)

	cartRepo.On(
		"Create",
		mock.Anything,
		mock.Anything,
	).Return(
		errors.New("database error"),
	)

	result, err := service.AddToCart(
		context.Background(),
		1,
		request.AddToCartRequest{
			ProductID: 1,
			Quantity:  2,
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"database error",
	)

	productRepo.AssertExpectations(t)
	cartRepo.AssertExpectations(t)
}

func TestCartAddFindByIDRepositoryError(t *testing.T) {

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)

	service := CartServiceImpl{
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}

	product := &domain.Product{
		ID:    1,
		Price: 750000,
		Stock: 10,
	}

	productRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(product, nil)

	cartRepo.On(
		"FindByUserIDAndProductID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(nil, gorm.ErrRecordNotFound)

	cartRepo.On(
		"Create",
		mock.Anything,
		mock.Anything,
	).Run(func(args mock.Arguments) {

		cart := args.Get(1).(*domain.Cart)
		cart.ID = 1

	}).Return(nil)

	cartRepo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(
		nil,
		errors.New("database error"),
	)

	result, err := service.AddToCart(
		context.Background(),
		1,
		request.AddToCartRequest{
			ProductID: 1,
			Quantity:  2,
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"database error",
	)

	productRepo.AssertExpectations(t)
	cartRepo.AssertExpectations(t)
}
