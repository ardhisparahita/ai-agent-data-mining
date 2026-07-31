package service

import (
	"context"
	"testing"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestFindAllSuccess(t *testing.T) {

	repo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: repo,
	}

	orders := []domain.Order{
		{
			ID:          1,
			UserID:      1,
			TotalAmount: 100000,
			Status:      domain.OrderPending,
		},
		{
			ID:          2,
			UserID:      1,
			TotalAmount: 200000,
			Status:      domain.OrderPaid,
		},
	}

	repo.On(
		"FindAllByUserID",
		mock.Anything,
		uint64(1),
	).Return(
		orders,
		nil,
	)

	result, err := service.FindAll(
		context.Background(),
		1,
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)

	assert.Equal(t, uint64(1), result[0].ID)
	assert.Equal(t, domain.OrderPending, result[0].Status)

	assert.Equal(t, uint64(2), result[1].ID)
	assert.Equal(t, domain.OrderPaid, result[1].Status)

	repo.AssertExpectations(t)
}

func TestFindAllRepositoryError(t *testing.T) {

	repo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: repo,
	}

	repo.On(
		"FindAllByUserID",
		mock.Anything,
		uint64(1),
	).Return(
		nil,
		gorm.ErrInvalidDB,
	)

	result, err := service.FindAll(
		context.Background(),
		1,
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
}

func TestFindByIDSuccess(t *testing.T) {
	repo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: repo,
	}

	order := &domain.Order{
		ID:          1,
		UserID:      1,
		TotalAmount: 150000,
		Status:      domain.OrderPending,
	}

	repo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(
		order,
		nil,
	)

	result, err := service.FindByID(
		context.Background(),
		1,
		1,
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, uint64(1), result.ID)
	assert.Equal(t, domain.OrderPending, result.Status)

	repo.AssertExpectations(t)
}

func TestFindByIDNotFound(t *testing.T) {

	repo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: repo,
	}

	repo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(99),
		uint64(1),
	).Return(
		nil,
		gorm.ErrRecordNotFound,
	)

	result, err := service.FindByID(
		context.Background(),
		99,
		1,
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"order not found",
	)

	repo.AssertExpectations(t)
}

func TestFindByIDRepositoryError(t *testing.T) {

	repo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: repo,
	}

	repo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(
		nil,
		gorm.ErrInvalidDB,
	)

	result, err := service.FindByID(
		context.Background(),
		1,
		1,
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
}
