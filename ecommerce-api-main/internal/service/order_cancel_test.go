package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCancelSuccess(t *testing.T) {

	db, sqlMock := testutils.NewMockDB(t)

	orderRepo := new(MockOrderRepository)
	productRepo := new(MockProductRepository)
	paymentRepo := new(MockPaymentRepository)

	service := OrderServiceImpl{
		DB:          db,
		OrderRepo:   orderRepo,
		ProductRepo: productRepo,
		PaymentRepo: paymentRepo,
	}

	order := &domain.Order{
		ID:     1,
		UserID: 1,
		Status: domain.OrderPending,
		OrderItems: []domain.OrderItem{
			{
				ProductID: 1,
				Quantity:  2,
			},
			{
				ProductID: 2,
				Quantity:  3,
			},
		},
	}

	product1 := &domain.Product{
		ID:    1,
		Stock: 5,
	}

	product2 := &domain.Product{
		ID:    2,
		Stock: 10,
	}

	payment := &domain.Payment{
		OrderID: 1,
		Status:  domain.PaymentPending,
	}

	orderRepo.On(
		"FindByIDAndUserIDWithItems",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(order, nil)

	paymentRepo.On(
		"FIndByOrderID",
		mock.Anything,
		uint64(1),
	).Return(payment, nil)

	productRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(product1, nil)

	productRepo.On(
		"FindByID",
		mock.Anything,
		uint64(2),
	).Return(product2, nil)

	sqlMock.ExpectBegin()

	productRepo.On(
		"UpdateTx",
		mock.Anything,
		mock.Anything,
		product1,
	).Return(nil)

	productRepo.On(
		"UpdateTx",
		mock.Anything,
		mock.Anything,
		product2,
	).Return(nil)

	paymentRepo.On(
		"UpdateTx",
		mock.Anything,
		mock.Anything,
		payment,
	).Return(nil)

	orderRepo.On(
		"UpdateTx",
		mock.Anything,
		mock.Anything,
		order,
	).Return(nil)

	sqlMock.ExpectCommit()

	err := service.Cancel(
		context.Background(),
		1,
		1,
	)

	assert.NoError(t, err)

	assert.Equal(
		t,
		domain.OrderCancelled,
		order.Status,
	)

	assert.Equal(
		t,
		domain.PaymentFailed,
		payment.Status,
	)

	assert.Equal(
		t,
		7,
		product1.Stock,
	)

	assert.Equal(
		t,
		13,
		product2.Stock,
	)

	assert.NoError(
		t,
		sqlMock.ExpectationsWereMet(),
	)

	orderRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
	paymentRepo.AssertExpectations(t)
}

func TestCancelOrderNotFound(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	orderRepo.On(
		"FindByIDAndUserIDWithItems",
		mock.Anything,
		uint64(99),
		uint64(1),
	).Return(
		nil,
		gorm.ErrRecordNotFound,
	)

	err := service.Cancel(
		context.Background(),
		99,
		1,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"order not found",
	)

	orderRepo.AssertExpectations(t)
}

func TestCancelPaymentNotFound(t *testing.T) {

	orderRepo := new(MockOrderRepository)
	paymentRepo := new(MockPaymentRepository)

	service := OrderServiceImpl{
		OrderRepo:   orderRepo,
		PaymentRepo: paymentRepo,
	}

	order := &domain.Order{
		ID:     1,
		UserID: 1,
		Status: domain.OrderPending,
	}

	orderRepo.On(
		"FindByIDAndUserIDWithItems",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(
		order,
		nil,
	)

	paymentRepo.On(
		"FIndByOrderID",
		mock.Anything,
		uint64(1),
	).Return(
		nil,
		gorm.ErrRecordNotFound,
	)

	err := service.Cancel(
		context.Background(),
		1,
		1,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"payment not found",
	)

	orderRepo.AssertExpectations(t)
	paymentRepo.AssertExpectations(t)
}

func TestCancelProductNotFound(t *testing.T) {

	db, sqlMock := testutils.NewMockDB(t)

	orderRepo := new(MockOrderRepository)
	productRepo := new(MockProductRepository)
	paymentRepo := new(MockPaymentRepository)

	service := OrderServiceImpl{
		DB:          db,
		OrderRepo:   orderRepo,
		ProductRepo: productRepo,
		PaymentRepo: paymentRepo,
	}

	order := &domain.Order{
		ID:     1,
		UserID: 1,
		Status: domain.OrderPending,
		OrderItems: []domain.OrderItem{
			{
				ProductID: 1,
				Quantity:  2,
			},
		},
	}

	payment := &domain.Payment{
		OrderID: 1,
		Status:  domain.PaymentPending,
	}

	orderRepo.On(
		"FindByIDAndUserIDWithItems",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(order, nil)

	paymentRepo.On(
		"FIndByOrderID",
		mock.Anything,
		uint64(1),
	).Return(payment, nil)

	sqlMock.ExpectBegin()

	productRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(
		nil,
		gorm.ErrRecordNotFound,
	)

	sqlMock.ExpectRollback()

	err := service.Cancel(
		context.Background(),
		1,
		1,
	)

	assert.Error(t, err)

	assert.True(
		t,
		errors.Is(err, gorm.ErrRecordNotFound),
	)

	assert.NoError(
		t,
		sqlMock.ExpectationsWereMet(),
	)

	orderRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
	paymentRepo.AssertExpectations(t)
}

func TestCancelRollbackWhenUpdateProductFailed(t *testing.T) {

	db, sqlMock := testutils.NewMockDB(t)

	orderRepo := new(MockOrderRepository)
	productRepo := new(MockProductRepository)
	paymentRepo := new(MockPaymentRepository)

	service := OrderServiceImpl{
		DB:          db,
		OrderRepo:   orderRepo,
		ProductRepo: productRepo,
		PaymentRepo: paymentRepo,
	}

	order := &domain.Order{
		ID:     1,
		UserID: 1,
		Status: domain.OrderPending,
		OrderItems: []domain.OrderItem{
			{
				ProductID: 1,
				Quantity:  2,
			},
		},
	}

	product := &domain.Product{
		ID:    1,
		Stock: 5,
	}

	payment := &domain.Payment{
		OrderID: 1,
		Status:  domain.PaymentPending,
	}

	orderRepo.On(
		"FindByIDAndUserIDWithItems",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(order, nil)

	paymentRepo.On(
		"FIndByOrderID",
		mock.Anything,
		uint64(1),
	).Return(payment, nil)

	productRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(product, nil)

	sqlMock.ExpectBegin()

	productRepo.On(
		"UpdateTx",
		mock.Anything,
		mock.Anything,
		product,
	).Return(assert.AnError)

	sqlMock.ExpectRollback()

	err := service.Cancel(
		context.Background(),
		1,
		1,
	)

	assert.Error(t, err)

	assert.NoError(t, sqlMock.ExpectationsWereMet())

	orderRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
	paymentRepo.AssertExpectations(t)
}

func TestCancelRollbackWhenUpdatePaymentFailed(t *testing.T) {

	db, sqlMock := testutils.NewMockDB(t)

	orderRepo := new(MockOrderRepository)
	productRepo := new(MockProductRepository)
	paymentRepo := new(MockPaymentRepository)

	service := OrderServiceImpl{
		DB:          db,
		OrderRepo:   orderRepo,
		ProductRepo: productRepo,
		PaymentRepo: paymentRepo,
	}

	order := &domain.Order{
		ID:     1,
		UserID: 1,
		Status: domain.OrderPending,
		OrderItems: []domain.OrderItem{
			{
				ProductID: 1,
				Quantity:  2,
			},
		},
	}

	product := &domain.Product{
		ID:    1,
		Stock: 5,
	}

	payment := &domain.Payment{
		OrderID: 1,
		Status:  domain.PaymentPending,
	}

	orderRepo.On(
		"FindByIDAndUserIDWithItems",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(order, nil)

	paymentRepo.On(
		"FIndByOrderID",
		mock.Anything,
		uint64(1),
	).Return(payment, nil)

	productRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(product, nil)

	sqlMock.ExpectBegin()

	productRepo.On(
		"UpdateTx",
		mock.Anything,
		mock.Anything,
		product,
	).Return(nil)

	paymentRepo.On(
		"UpdateTx",
		mock.Anything,
		mock.Anything,
		payment,
	).Return(assert.AnError)

	sqlMock.ExpectRollback()

	err := service.Cancel(
		context.Background(),
		1,
		1,
	)

	assert.Error(t, err)

	assert.NoError(t, sqlMock.ExpectationsWereMet())

	orderRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
	paymentRepo.AssertExpectations(t)
}

func TestCancelRollbackWhenUpdateOrderFailed(t *testing.T) {

	db, sqlMock := testutils.NewMockDB(t)

	orderRepo := new(MockOrderRepository)
	productRepo := new(MockProductRepository)
	paymentRepo := new(MockPaymentRepository)

	service := OrderServiceImpl{
		DB:          db,
		OrderRepo:   orderRepo,
		ProductRepo: productRepo,
		PaymentRepo: paymentRepo,
	}

	order := &domain.Order{
		ID:     1,
		UserID: 1,
		Status: domain.OrderPending,
		OrderItems: []domain.OrderItem{
			{
				ProductID: 1,
				Quantity:  2,
			},
		},
	}

	product := &domain.Product{
		ID:    1,
		Stock: 5,
	}

	payment := &domain.Payment{
		OrderID: 1,
		Status:  domain.PaymentPending,
	}

	orderRepo.On(
		"FindByIDAndUserIDWithItems",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(order, nil)

	paymentRepo.On(
		"FIndByOrderID",
		mock.Anything,
		uint64(1),
	).Return(payment, nil)

	productRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(product, nil)

	sqlMock.ExpectBegin()

	productRepo.On(
		"UpdateTx",
		mock.Anything,
		mock.Anything,
		product,
	).Return(nil)

	paymentRepo.On(
		"UpdateTx",
		mock.Anything,
		mock.Anything,
		payment,
	).Return(nil)

	orderRepo.On(
		"UpdateTx",
		mock.Anything,
		mock.Anything,
		order,
	).Return(assert.AnError)

	sqlMock.ExpectRollback()

	err := service.Cancel(
		context.Background(),
		1,
		1,
	)

	assert.Error(t, err)

	assert.NoError(t, sqlMock.ExpectationsWereMet())

	orderRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
	paymentRepo.AssertExpectations(t)
}

func TestCancelAlreadyCancelled(t *testing.T) {

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
		"FindByIDAndUserIDWithItems",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(order, nil)

	err := service.Cancel(
		context.Background(),
		1,
		1,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "order already cancelled")

	orderRepo.AssertExpectations(t)
}

func TestCancelPaid(t *testing.T) {

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
		"FindByIDAndUserIDWithItems",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(order, nil)

	err := service.Cancel(
		context.Background(),
		1,
		1,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "paid order cannot be cancelled")

	orderRepo.AssertExpectations(t)
}

func TestCancelShipped(t *testing.T) {

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
		"FindByIDAndUserIDWithItems",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(order, nil)

	err := service.Cancel(
		context.Background(),
		1,
		1,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "shipped order cannot be cancelled")

	orderRepo.AssertExpectations(t)
}

func TestCancelCompleted(t *testing.T) {

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
		"FindByIDAndUserIDWithItems",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(order, nil)

	err := service.Cancel(
		context.Background(),
		1,
		1,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "completed order cannot be cancelled")

	orderRepo.AssertExpectations(t)
}

func TestCancelInvalidStatus(t *testing.T) {

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
		"FindByIDAndUserIDWithItems",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(order, nil)

	err := service.Cancel(
		context.Background(),
		1,
		1,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid order status")

	orderRepo.AssertExpectations(t)
}
