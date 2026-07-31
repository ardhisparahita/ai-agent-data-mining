package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCartFindAllSuccess(t *testing.T) {

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)

	service := CartServiceImpl{
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}

	carts := []domain.Cart{
		{
			ID:        1,
			UserID:    1,
			ProductID: 1,
			Quantity:  2,
			Product: domain.Product{
				ID:       1,
				Name:     "Keyboard",
				Price:    750000,
				ImageURL: "keyboard.jpg",
			},
		},
		{
			ID:        2,
			UserID:    1,
			ProductID: 2,
			Quantity:  1,
			Product: domain.Product{
				ID:       2,
				Name:     "Mouse",
				Price:    250000,
				ImageURL: "mouse.jpg",
			},
		},
	}

	cartRepo.On(
		"FindAllByUserID",
		mock.Anything,
		uint64(1),
	).Return(carts, nil)

	result, err := service.FindAll(
		context.Background(),
		1,
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Len(t, result.Items, 2)

	assert.Equal(t, 2, result.TotalItems)
	assert.Equal(t, float64(1750000), result.GrandTotal)

	assert.Equal(t, uint64(1), result.Items[0].ID)
	assert.Equal(t, "Keyboard", result.Items[0].Product.Name)

	assert.Equal(t, uint64(2), result.Items[1].ID)
	assert.Equal(t, "Mouse", result.Items[1].Product.Name)

	cartRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestCartFindAllEmpty(t *testing.T) {

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)

	service := CartServiceImpl{
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}

	cartRepo.On(
		"FindAllByUserID",
		mock.Anything,
		uint64(1),
	).Return([]domain.Cart{}, nil)

	result, err := service.FindAll(
		context.Background(),
		1,
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Len(t, result.Items, 0)
	assert.Equal(t, 0, result.TotalItems)
	assert.Equal(t, float64(0), result.GrandTotal)

	cartRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestCartFindAllRepositoryError(t *testing.T) {

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)

	service := CartServiceImpl{
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}

	cartRepo.On(
		"FindAllByUserID",
		mock.Anything,
		uint64(1),
	).Return(nil, errors.New("database error"))

	result, err := service.FindAll(
		context.Background(),
		1,
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"database error",
	)

	cartRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}
