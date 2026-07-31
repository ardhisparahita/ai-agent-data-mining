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

func TestProductUpdateSuccess(t *testing.T) {

	repo := new(MockProductRepository)

	service := ProductServiceImpl{
		Repo: repo,
	}

	product := &domain.Product{
		ID:          1,
		CategoryID:  1,
		Name:        "Old Keyboard",
		Description: "Old Description",
		Price:       500000,
		Stock:       5,
	}

	repo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(
		product,
		nil,
	)

	repo.On(
		"Update",
		mock.Anything,
		product,
	).Return(nil)

	result, err := service.Update(
		context.Background(),
		1,
		request.UpdateProductRequest{
			CategoryID:  2,
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
		uint64(2),
		result.CategoryID,
	)

	assert.Equal(
		t,
		"Mechanical Keyboard",
		result.Name,
	)

	assert.Equal(
		t,
		"RGB Keyboard",
		result.Description,
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

func TestProductUpdateNotFound(t *testing.T) {

	repo := new(MockProductRepository)

	service := ProductServiceImpl{
		Repo: repo,
	}

	repo.On(
		"FindByID",
		mock.Anything,
		uint64(99),
	).Return(
		nil,
		gorm.ErrRecordNotFound,
	)

	result, err := service.Update(
		context.Background(),
		99,
		request.UpdateProductRequest{
			CategoryID:  1,
			Name:        "Keyboard",
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
		"product not found",
	)

	repo.AssertExpectations(t)
}

func TestProductUpdateRepositoryError(t *testing.T) {

	repo := new(MockProductRepository)

	service := ProductServiceImpl{
		Repo: repo,
	}

	product := &domain.Product{
		ID:          1,
		CategoryID:  1,
		Name:        "Keyboard",
		Description: "RGB Keyboard",
		Price:       750000,
		Stock:       10,
	}

	repo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(
		product,
		nil,
	)

	repo.On(
		"Update",
		mock.Anything,
		product,
	).Return(
		errors.New("database error"),
	)

	result, err := service.Update(
		context.Background(),
		1,
		request.UpdateProductRequest{
			CategoryID:  2,
			Name:        "Mechanical Keyboard",
			Description: "Updated Description",
			Price:       800000,
			Stock:       20,
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
