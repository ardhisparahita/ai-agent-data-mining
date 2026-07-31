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

func TestProductDeleteSuccess(t *testing.T) {

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
		"Delete",
		mock.Anything,
		uint64(1),
	).Return(nil)

	err := service.Delete(
		context.Background(),
		1,
	)

	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestProductDeleteNotFound(t *testing.T) {

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

	err := service.Delete(
		context.Background(),
		99,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"product not found",
	)

	repo.AssertExpectations(t)
}

func TestProductDeleteRepositoryError(t *testing.T) {

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
		"Delete",
		mock.Anything,
		uint64(1),
	).Return(
		errors.New("database error"),
	)

	err := service.Delete(
		context.Background(),
		1,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"database error",
	)

	repo.AssertExpectations(t)
}
