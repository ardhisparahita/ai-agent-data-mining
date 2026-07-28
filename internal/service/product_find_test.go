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

func TestProductFindAllSuccess(t *testing.T) {

	repo := new(MockProductRepository)

	service := ProductServiceImpl{
		Repo: repo,
	}

	req := request.ProductQueryRequest{
		Page:  1,
		Limit: 10,
	}

	repo.On(
		"FindAll",
		mock.Anything,
		req,
	).Return(
		[]domain.Product{
			{
				ID:          1,
				CategoryID:  1,
				Name:        "Keyboard",
				Description: "RGB Keyboard",
				Price:       750000,
				Stock:       10,
			},
			{
				ID:          2,
				CategoryID:  1,
				Name:        "Mouse",
				Description: "Gaming Mouse",
				Price:       350000,
				Stock:       20,
			},
		},
		int64(2),
		nil,
	)

	result, err := service.FindAll(
		context.Background(),
		req,
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Len(
		t,
		result.Items,
		2,
	)

	assert.Equal(
		t,
		int64(2),
		result.TotalRows,
	)

	assert.Equal(
		t,
		1,
		result.Page,
	)

	assert.Equal(
		t,
		10,
		result.Limit,
	)

	assert.Equal(
		t,
		1,
		result.TotalPages,
	)

	repo.AssertExpectations(t)
}

func TestProductFindAllRepositoryError(t *testing.T) {

	repo := new(MockProductRepository)

	service := ProductServiceImpl{
		Repo: repo,
	}

	req := request.ProductQueryRequest{
		Page:  1,
		Limit: 10,
	}

	repo.On(
		"FindAll",
		mock.Anything,
		req,
	).Return(
		[]domain.Product{},
		int64(0),
		errors.New("database error"),
	)

	result, err := service.FindAll(
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

func TestProductFindByIDSuccess(t *testing.T) {

	repo := new(MockProductRepository)

	service := ProductServiceImpl{
		Repo: repo,
	}

	repo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(
		&domain.Product{
			ID:          1,
			CategoryID:  1,
			Name:        "Keyboard",
			Description: "RGB Keyboard",
			Price:       750000,
			Stock:       10,
		},
		nil,
	)

	result, err := service.FindByID(
		context.Background(),
		1,
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
		"Keyboard",
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

func TestProductFindByIDNotFound(t *testing.T) {

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

	result, err := service.FindByID(
		context.Background(),
		99,
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

func TestProductFindByIDRepositoryError(t *testing.T) {

	repo := new(MockProductRepository)

	service := ProductServiceImpl{
		Repo: repo,
	}

	repo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(
		nil,
		errors.New("database error"),
	)

	result, err := service.FindByID(
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

	repo.AssertExpectations(t)
}
