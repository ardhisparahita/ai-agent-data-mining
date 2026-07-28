package service

import (
	"context"
	"testing"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestMarkAsShippedSuccess(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	order := &domain.Order{
		ID:     1,
		Status: domain.OrderPaid,
	}

	orderRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	orderRepo.On(
		"Update",
		mock.Anything,
		order,
	).Return(nil)

	err := service.MarkAsShipped(
		context.Background(),
		1,
	)

	assert.NoError(t, err)

	assert.Equal(
		t,
		domain.OrderShipped,
		order.Status,
	)

	orderRepo.AssertExpectations(t)
}

func TestMarkAsShippedOrderNotFound(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	orderRepo.On(
		"FindByID",
		mock.Anything,
		uint64(99),
	).Return(
		nil,
		gorm.ErrRecordNotFound,
	)

	err := service.MarkAsShipped(
		context.Background(),
		99,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"order not found",
	)

	orderRepo.AssertExpectations(t)
}

func TestMarkAsShippedPending(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	order := &domain.Order{
		ID:     1,
		Status: domain.OrderPending,
	}

	orderRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	err := service.MarkAsShipped(
		context.Background(),
		1,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"order must be paid before shipping",
	)

	orderRepo.AssertExpectations(t)
}

func TestMarkAsShippedCancelled(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	order := &domain.Order{
		ID:     1,
		Status: domain.OrderCancelled,
	}

	orderRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	err := service.MarkAsShipped(
		context.Background(),
		1,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"cancelled order cannot be shipped",
	)

	orderRepo.AssertExpectations(t)
}

func TestMarkAsShippedAlreadyShipped(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	order := &domain.Order{
		ID:     1,
		Status: domain.OrderShipped,
	}

	orderRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	err := service.MarkAsShipped(
		context.Background(),
		1,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"order already shipped",
	)

	orderRepo.AssertExpectations(t)
}

func TestMarkAsShippedAlreadyCompleted(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	order := &domain.Order{
		ID:     1,
		Status: domain.OrderCompleted,
	}

	orderRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	err := service.MarkAsShipped(
		context.Background(),
		1,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"completed order cannot be shipped",
	)

	orderRepo.AssertExpectations(t)
}

func TestMarkAsShippedInvalidStatus(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	order := &domain.Order{
		ID:     1,
		Status: "UNKNOWN_STATUS",
	}

	orderRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	err := service.MarkAsShipped(
		context.Background(),
		1,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"invalid order status",
	)

	orderRepo.AssertExpectations(t)
}

func TestMarkAsCompletedSuccessAsAdmin(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	order := &domain.Order{
		ID:     1,
		UserID: 999,
		Status: domain.OrderShipped,
	}

	orderRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	orderRepo.On(
		"Update",
		mock.Anything,
		order,
	).Return(nil)

	err := service.MarkAsCompleted(
		context.Background(),
		1,
		1,
		"ADMIN",
	)

	assert.NoError(t, err)

	assert.Equal(
		t,
		domain.OrderCompleted,
		order.Status,
	)

	orderRepo.AssertExpectations(t)
}

func TestMarkAsCompletedSuccessAsCustomer(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	order := &domain.Order{
		ID:     1,
		UserID: 1,
		Status: domain.OrderShipped,
	}

	orderRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	orderRepo.On(
		"Update",
		mock.Anything,
		order,
	).Return(nil)

	err := service.MarkAsCompleted(
		context.Background(),
		1,
		1,
		"CUSTOMER",
	)

	assert.NoError(t, err)

	assert.Equal(
		t,
		domain.OrderCompleted,
		order.Status,
	)

	orderRepo.AssertExpectations(t)
}

func TestMarkAsCompletedOrderNotFound(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	orderRepo.On(
		"FindByID",
		mock.Anything,
		uint64(99),
	).Return(
		nil,
		gorm.ErrRecordNotFound,
	)

	err := service.MarkAsCompleted(
		context.Background(),
		99,
		1,
		"ADMIN",
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"order not found",
	)

	orderRepo.AssertExpectations(t)
}

func TestMarkAsCompletedForbidden(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	order := &domain.Order{
		ID:     1,
		UserID: 10,
		Status: domain.OrderShipped,
	}

	orderRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	err := service.MarkAsCompleted(
		context.Background(),
		1,
		1,
		"CUSTOMER",
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"you are not allowed",
	)

	orderRepo.AssertExpectations(t)
}

func TestMarkAsCompletedAlreadyCompleted(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	order := &domain.Order{
		ID:     1,
		UserID: 1,
		Status: domain.OrderCompleted,
	}

	orderRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	err := service.MarkAsCompleted(
		context.Background(),
		1,
		1,
		"ADMIN",
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"order already completed",
	)

	orderRepo.AssertExpectations(t)
}

func TestMarkAsCompletedCancelled(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	order := &domain.Order{
		ID:     1,
		UserID: 1,
		Status: domain.OrderCancelled,
	}

	orderRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	err := service.MarkAsCompleted(
		context.Background(),
		1,
		1,
		"ADMIN",
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"cancelled order cannot be completed",
	)

	orderRepo.AssertExpectations(t)
}

func TestMarkAsCompletedPending(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	order := &domain.Order{
		ID:     1,
		UserID: 1,
		Status: domain.OrderPending,
	}

	orderRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	err := service.MarkAsCompleted(
		context.Background(),
		1,
		1,
		"ADMIN",
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"order must be shipped first",
	)

	orderRepo.AssertExpectations(t)
}

func TestMarkAsCompletedPaid(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	order := &domain.Order{
		ID:     1,
		UserID: 1,
		Status: domain.OrderPaid,
	}

	orderRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	err := service.MarkAsCompleted(
		context.Background(),
		1,
		1,
		"ADMIN",
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"order must be shipped first",
	)

	orderRepo.AssertExpectations(t)
}

func TestMarkAsCompletedInvalidStatus(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	order := &domain.Order{
		ID:     1,
		UserID: 1,
		Status: "UNKNOWN_STATUS",
	}

	orderRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	err := service.MarkAsCompleted(
		context.Background(),
		1,
		1,
		"ADMIN",
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"invalid order status",
	)

	orderRepo.AssertExpectations(t)
}
