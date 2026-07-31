package service

import (
	"context"
	"testing"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestMarkAsPaidSuccess(t *testing.T) {

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
		Status: domain.OrderPending,
	}

	payment := &domain.Payment{
		OrderID: 1,
		Status:  domain.PaymentPending,
	}

	orderRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	paymentRepo.On(
		"FIndByOrderID",
		mock.Anything,
		uint64(1),
	).Return(payment, nil)

	sqlMock.ExpectBegin()

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

	err := service.MarkAsPaid(
		context.Background(),
		1,
	)

	assert.NoError(t, err)

	assert.Equal(
		t,
		domain.OrderPaid,
		order.Status,
	)

	assert.Equal(
		t,
		domain.PaymentPaid,
		payment.Status,
	)

	assert.NoError(
		t,
		sqlMock.ExpectationsWereMet(),
	)

	orderRepo.AssertExpectations(t)
	paymentRepo.AssertExpectations(t)
}

func TestMarkAsPaidOrderNotFound(t *testing.T) {

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

	err := service.MarkAsPaid(
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

func TestMarkAsPaidPaymentNotFound(t *testing.T) {

	orderRepo := new(MockOrderRepository)
	paymentRepo := new(MockPaymentRepository)

	service := OrderServiceImpl{
		OrderRepo:   orderRepo,
		PaymentRepo: paymentRepo,
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

	paymentRepo.On(
		"FIndByOrderID",
		mock.Anything,
		uint64(1),
	).Return(
		nil,
		gorm.ErrRecordNotFound,
	)

	err := service.MarkAsPaid(
		context.Background(),
		1,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"payment not found",
	)
}

func TestMarkAsPaidAlreadyPaid(t *testing.T) {

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

	err := service.MarkAsPaid(
		context.Background(),
		1,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"already paid",
	)
}

func TestMarkAsPaidAlreadyShipped(t *testing.T) {

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

	err := service.MarkAsPaid(
		context.Background(),
		1,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"already shipped",
	)
}

func TestMarkAsPaidAlreadyCompleted(t *testing.T) {

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

	err := service.MarkAsPaid(
		context.Background(),
		1,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"already completed",
	)
}

func TestMarkAsPaidCancelled(t *testing.T) {

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

	err := service.MarkAsPaid(
		context.Background(),
		1,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"cannot be paid",
	)
}

func TestMarkAsPaidInvalidStatus(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	order := &domain.Order{
		ID:     1,
		Status: "UNKNOWN",
	}

	orderRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	err := service.MarkAsPaid(
		context.Background(),
		1,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"invalid order status",
	)
}

func TestMarkAsPaidRollbackWhenUpdatePaymentFailed(t *testing.T) {

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
		Status: domain.OrderPending,
	}

	payment := &domain.Payment{
		OrderID: 1,
		Status:  domain.PaymentPending,
	}

	orderRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	paymentRepo.On(
		"FIndByOrderID",
		mock.Anything,
		uint64(1),
	).Return(payment, nil)

	sqlMock.ExpectBegin()

	paymentRepo.On(
		"UpdateTx",
		mock.Anything,
		mock.Anything,
		payment,
	).Return(assert.AnError)

	sqlMock.ExpectRollback()

	err := service.MarkAsPaid(
		context.Background(),
		1,
	)

	assert.Error(t, err)

	assert.Equal(
		t,
		domain.PaymentPaid,
		payment.Status,
	)

	assert.Equal(
		t,
		domain.OrderPending,
		order.Status,
	)

	assert.NoError(
		t,
		sqlMock.ExpectationsWereMet(),
	)

	orderRepo.AssertExpectations(t)
	paymentRepo.AssertExpectations(t)
}

func TestMarkAsPaidRollbackWhenUpdateOrderFailed(t *testing.T) {

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
		Status: domain.OrderPending,
	}

	payment := &domain.Payment{
		OrderID: 1,
		Status:  domain.PaymentPending,
	}

	orderRepo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	paymentRepo.On(
		"FIndByOrderID",
		mock.Anything,
		uint64(1),
	).Return(payment, nil)

	sqlMock.ExpectBegin()

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

	err := service.MarkAsPaid(
		context.Background(),
		1,
	)

	assert.Error(t, err)

	assert.Equal(
		t,
		domain.PaymentPaid,
		payment.Status,
	)

	assert.Equal(
		t,
		domain.OrderPaid,
		order.Status,
	)

	assert.NoError(
		t,
		sqlMock.ExpectationsWereMet(),
	)

	orderRepo.AssertExpectations(t)
	paymentRepo.AssertExpectations(t)
}

func TestMarkAsFailedSuccess(t *testing.T) {

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
		"FindByIDWithItems",
		mock.Anything,
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

	err := service.MarkAsFailed(
		context.Background(),
		1,
	)

	assert.NoError(t, err)

	assert.Equal(
		t,
		domain.PaymentFailed,
		payment.Status,
	)

	assert.Equal(
		t,
		domain.OrderCancelled,
		order.Status,
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

func TestMarkAsFailedOrderNotFound(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	orderRepo.On(
		"FindByIDWithItems",
		mock.Anything,
		uint64(99),
	).Return(
		nil,
		gorm.ErrRecordNotFound,
	)

	err := service.MarkAsFailed(
		context.Background(),
		99,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"order not found",
	)
}

func TestMarkAsFailedPaymentNotFound(t *testing.T) {

	orderRepo := new(MockOrderRepository)
	paymentRepo := new(MockPaymentRepository)

	service := OrderServiceImpl{
		OrderRepo:   orderRepo,
		PaymentRepo: paymentRepo,
	}

	order := &domain.Order{
		ID:     1,
		Status: domain.OrderPending,
	}

	orderRepo.On(
		"FindByIDWithItems",
		mock.Anything,
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

	err := service.MarkAsFailed(
		context.Background(),
		1,
	)

	assert.Error(t, err)
}

func TestMarkAsFailedRollbackWhenUpdateProductFailed(t *testing.T) {

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
		"FindByIDWithItems",
		mock.Anything,
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

	err := service.MarkAsFailed(
		context.Background(),
		1,
	)

	assert.Error(t, err)

	assert.NoError(t, sqlMock.ExpectationsWereMet())

	orderRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
	paymentRepo.AssertExpectations(t)
}

func TestMarkAsFailedRollbackWhenUpdatePaymentFailed(t *testing.T) {

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
		"FindByIDWithItems",
		mock.Anything,
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

	err := service.MarkAsFailed(
		context.Background(),
		1,
	)

	assert.Error(t, err)

	assert.NoError(t, sqlMock.ExpectationsWereMet())

	orderRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
	paymentRepo.AssertExpectations(t)
}

func TestMarkAsFailedRollbackWhenUpdateOrderFailed(t *testing.T) {

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
		"FindByIDWithItems",
		mock.Anything,
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

	err := service.MarkAsFailed(
		context.Background(),
		1,
	)

	assert.Error(t, err)

	assert.NoError(t, sqlMock.ExpectationsWereMet())

	orderRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
	paymentRepo.AssertExpectations(t)
}

func TestMarkAsFailedAlreadyPaid(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	order := &domain.Order{
		ID:     1,
		Status: domain.OrderPaid,
	}

	orderRepo.On(
		"FindByIDWithItems",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	err := service.MarkAsFailed(
		context.Background(),
		1,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "paid order cannot be marked as failed")

	orderRepo.AssertExpectations(t)
}

func TestMarkAsFailedAlreadyShipped(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	order := &domain.Order{
		ID:     1,
		Status: domain.OrderShipped,
	}

	orderRepo.On(
		"FindByIDWithItems",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	err := service.MarkAsFailed(
		context.Background(),
		1,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "shipped order cannot be marked as failed")

	orderRepo.AssertExpectations(t)
}

func TestMarkAsFailedAlreadyCompleted(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	order := &domain.Order{
		ID:     1,
		Status: domain.OrderCompleted,
	}

	orderRepo.On(
		"FindByIDWithItems",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	err := service.MarkAsFailed(
		context.Background(),
		1,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "completed order cannot be marked as failed")

	orderRepo.AssertExpectations(t)
}

func TestMarkAsFailedAlreadyCancelled(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	order := &domain.Order{
		ID:     1,
		Status: domain.OrderCancelled,
	}

	orderRepo.On(
		"FindByIDWithItems",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	err := service.MarkAsFailed(
		context.Background(),
		1,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "order already cancelled")

	orderRepo.AssertExpectations(t)
}

func TestMarkAsFailedInvalidStatus(t *testing.T) {

	orderRepo := new(MockOrderRepository)

	service := OrderServiceImpl{
		OrderRepo: orderRepo,
	}

	order := &domain.Order{
		ID:     1,
		Status: "UNKNOWN_STATUS",
	}

	orderRepo.On(
		"FindByIDWithItems",
		mock.Anything,
		uint64(1),
	).Return(order, nil)

	err := service.MarkAsFailed(
		context.Background(),
		1,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid order status")

	orderRepo.AssertExpectations(t)
}
