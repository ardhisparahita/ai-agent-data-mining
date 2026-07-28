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

func TestProductUploadImageSuccess(t *testing.T) {

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
		ImageURL:    "",
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

	result, err := service.UploadImage(
		context.Background(),
		1,
		"keyboard.jpg",
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(
		t,
		"keyboard.jpg",
		result.ImageURL,
	)

	repo.AssertExpectations(t)
}

func TestProductUploadImageNotFound(t *testing.T) {

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

	result, err := service.UploadImage(
		context.Background(),
		99,
		"keyboard.jpg",
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

func TestProductUploadImageRepositoryError(t *testing.T) {

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
		ImageURL:    "",
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

	result, err := service.UploadImage(
		context.Background(),
		1,
		"keyboard.jpg",
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
