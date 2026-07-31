package service

import (
	"context"
	"testing"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateCategorySuccess(t *testing.T) {
	repo := new(MockCategoryRepository)

	service := CategoryServiceImpl{
		Repo: repo,
	}

	repo.On(
		"Create",
		mock.Anything,
		mock.Anything,
	).Run(func(args mock.Arguments) {
		category := args.Get(1).(*domain.Category)

		category.ID = 1
	}).Return(nil)

	result, err := service.Create(context.Background(), request.CreateCategoryRequest{
		Name: "Electronics",
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Electronics", result.Name)

	repo.AssertExpectations(t)
}

func TestFindAllCategorySuccess(t *testing.T) {
	repo := new(MockCategoryRepository)

	service := CategoryServiceImpl{
		Repo: repo,
	}

	repo.On(
		"FindAll",
		mock.Anything,
	).Return(
		[]domain.Category{
			{
				ID:   1,
				Name: "Electronics",
			},
			{
				ID:   1,
				Name: "Fashion",
			},
		},
		nil,
	)

	result, err := service.FindAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 2)

	repo.AssertExpectations(t)
}

func TestUpdateCategorySuccess(t *testing.T) {
	repo := new(MockCategoryRepository)

	service := CategoryServiceImpl{
		Repo: repo,
	}

	repo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(
		&domain.Category{
			ID:   1,
			Name: "Old Name",
		},
		nil,
	)

	repo.On(
		"Update",
		mock.Anything,
		mock.Anything,
	).Return(nil)

	result, err := service.Update(context.Background(), 1, request.UpdateCategoryRequest{
		Name: "New Name",
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Name", result.Name)

	repo.AssertExpectations(t)
}

func TestUpdateCategoryNotFound(t *testing.T) {
	repo := new(MockCategoryRepository)

	service := CategoryServiceImpl{
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

	result, err := service.Update(context.Background(), 99, request.UpdateCategoryRequest{
		Name: "New Name",
	})

	assert.Error(t, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
}
