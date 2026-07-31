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

func TestProductCreateSuccess(t *testing.T) {

	repo := new(MockProductRepository)

	service := ProductServiceImpl{
		Repo: repo,
	}

	repo.On(
		"Create",
		mock.Anything,
		mock.Anything,
	).Run(func(args mock.Arguments) {

		product := args.Get(1).(*domain.Product)
		product.ID = 1

	}).Return(nil)

	repo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(
		&domain.Product{
			ID:          1,
			CategoryID:  1,
			Name:        "Mechanical Keyboard",
			Description: "RGB Keyboard",
			Price:       750000,
			Stock:       10,
		},
		nil,
	)

	result, err := service.Create(
		context.Background(),
		request.CreateProductRequest{
			CategoryID:  1,
			Name:        "Mechanical Keyboard",
			Description: "RGB Keyboard",
			Price:       750000,
			Stock:       10,
		},
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
		"Mechanical Keyboard",
		result.Name,
	)

	assert.Equal(
		t,
		750000.0,
		result.Price,
	)

	assert.Equal(
		t,
		10,
		result.Stock,
	)

	repo.AssertExpectations(t)
}

func TestProductCreateRepositoryError(t *testing.T) {

	repo := new(MockProductRepository)

	service := ProductServiceImpl{
		Repo: repo,
	}

	repo.On(
		"Create",
		mock.Anything,
		mock.Anything,
	).Return(
		errors.New("database error"),
	)

	result, err := service.Create(
		context.Background(),
		request.CreateProductRequest{
			CategoryID:  1,
			Name:        "Mechanical Keyboard",
			Description: "RGB Keyboard",
			Price:       750000,
			Stock:       10,
		},
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

func TestProductCreateFindByIDError(t *testing.T) {

	repo := new(MockProductRepository)

	service := ProductServiceImpl{
		Repo: repo,
	}

	repo.On(
		"Create",
		mock.Anything,
		mock.Anything,
	).Run(func(args mock.Arguments) {

		product := args.Get(1).(*domain.Product)
		product.ID = 1

	}).Return(nil)

	repo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(
		nil,
		gorm.ErrRecordNotFound,
	)

	result, err := service.Create(
		context.Background(),
		request.CreateProductRequest{
			CategoryID:  1,
			Name:        "Mechanical Keyboard",
			Description: "RGB Keyboard",
			Price:       750000,
			Stock:       10,
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.True(
		t,
		errors.Is(err, gorm.ErrRecordNotFound),
	)

	repo.AssertExpectations(t)
}
