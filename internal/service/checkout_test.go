package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	"github.com/ardhisparahita/ecommerce-api/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCheckoutSuccess(t *testing.T) {

	db, sqlMock := testutils.NewMockDB(t)

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)
	addressRepo := new(MockAddressRepository)
	orderRepo := new(MockOrderRepository)
	orderItemRepo := new(MockOrderItemRepository)
	paymentRepo := new(MockPaymentRepository)

	service := CheckoutServiceImpl{
		DB:            db,
		CartRepo:      cartRepo,
		ProductRepo:   productRepo,
		AddressRepo:   addressRepo,
		OrderRepo:     orderRepo,
		OrderItemRepo: orderItemRepo,
		PaymentRepo:   paymentRepo,
	}

	address := &domain.Address{
		ID:            1,
		UserID:        1,
		RecipientName: "John Doe",
		Phone:         "08123456789",
		Address:       "Jl. Malioboro",
		City:          "Yogyakarta",
		PostalCode:    "55213",
	}

	carts := []domain.Cart{
		{
			ID:        1,
			UserID:    1,
			ProductID: 1,
			Quantity:  2,
			Product: domain.Product{
				ID:    1,
				Name:  "Keyboard",
				Price: 750000,
				Stock: 10,
			},
		},
		{
			ID:        2,
			UserID:    1,
			ProductID: 2,
			Quantity:  1,
			Product: domain.Product{
				ID:    2,
				Name:  "Mouse",
				Price: 250000,
				Stock: 5,
			},
		},
	}

	addressRepo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(address, nil)

	cartRepo.On(
		"FindAllByUserID",
		mock.Anything,
		uint64(1),
	).Return(carts, nil)

	sqlMock.ExpectBegin()

	orderRepo.On(
		"CreateTx",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*domain.Order"),
	).Run(func(args mock.Arguments) {

		order := args.Get(2).(*domain.Order)
		order.ID = 1

	}).Return(nil)

	orderItemRepo.On(
		"CreateTx",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*domain.OrderItem"),
	).Return(nil).Twice()

	productRepo.On(
		"UpdateTx",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*domain.Product"),
	).Return(nil).Twice()

	paymentRepo.On(
		"CreateTx",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*domain.Payment"),
	).Return(nil)

	cartRepo.On(
		"DeleteAllByUserIDTx",
		mock.Anything,
		mock.Anything,
		uint64(1),
	).Return(nil)

	sqlMock.ExpectCommit()

	result, err := service.Checkout(
		context.Background(),
		1,
		request.CheckoutRequest{
			AddressID:     1,
			PaymentMethod: "BANK_TRANSFER",
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
		float64(1750000),
		result.TotalAmount,
	)

	assert.Equal(
		t,
		domain.OrderPending,
		result.Status,
	)

	assert.Equal(
		t,
		"BANK_TRANSFER",
		result.PaymentMethod,
	)

	assert.NoError(
		t,
		sqlMock.ExpectationsWereMet(),
	)

	addressRepo.AssertExpectations(t)
	cartRepo.AssertExpectations(t)
	orderRepo.AssertExpectations(t)
	orderItemRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
	paymentRepo.AssertExpectations(t)
}

func TestCheckoutAddressNotFound(t *testing.T) {

	db, _ := testutils.NewMockDB(t)

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)
	addressRepo := new(MockAddressRepository)
	orderRepo := new(MockOrderRepository)
	orderItemRepo := new(MockOrderItemRepository)
	paymentRepo := new(MockPaymentRepository)

	service := CheckoutServiceImpl{
		DB:            db,
		CartRepo:      cartRepo,
		ProductRepo:   productRepo,
		AddressRepo:   addressRepo,
		OrderRepo:     orderRepo,
		OrderItemRepo: orderItemRepo,
		PaymentRepo:   paymentRepo,
	}

	addressRepo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(99),
		uint64(1),
	).Return(
		nil,
		gorm.ErrRecordNotFound,
	)

	result, err := service.Checkout(
		context.Background(),
		1,
		request.CheckoutRequest{
			AddressID:     99,
			PaymentMethod: "BANK_TRANSFER",
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"address not found",
	)

	addressRepo.AssertExpectations(t)
}

func TestCheckoutCartEmpty(t *testing.T) {

	db, _ := testutils.NewMockDB(t)

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)
	addressRepo := new(MockAddressRepository)
	orderRepo := new(MockOrderRepository)
	orderItemRepo := new(MockOrderItemRepository)
	paymentRepo := new(MockPaymentRepository)

	service := CheckoutServiceImpl{
		DB:            db,
		CartRepo:      cartRepo,
		ProductRepo:   productRepo,
		AddressRepo:   addressRepo,
		OrderRepo:     orderRepo,
		OrderItemRepo: orderItemRepo,
		PaymentRepo:   paymentRepo,
	}

	address := &domain.Address{
		ID:            1,
		UserID:        1,
		RecipientName: "John Doe",
		Phone:         "08123456789",
		Address:       "Jl. Malioboro",
		City:          "Yogyakarta",
		PostalCode:    "55213",
	}

	addressRepo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(address, nil)

	cartRepo.On(
		"FindAllByUserID",
		mock.Anything,
		uint64(1),
	).Return([]domain.Cart{}, nil)

	result, err := service.Checkout(
		context.Background(),
		1,
		request.CheckoutRequest{
			AddressID:     1,
			PaymentMethod: "BANK_TRANSFER",
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"cart is empty",
	)

	addressRepo.AssertExpectations(t)
	cartRepo.AssertExpectations(t)
}

func TestCheckoutInsufficientStock(t *testing.T) {

	db, _ := testutils.NewMockDB(t)

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)
	addressRepo := new(MockAddressRepository)
	orderRepo := new(MockOrderRepository)
	orderItemRepo := new(MockOrderItemRepository)
	paymentRepo := new(MockPaymentRepository)

	service := CheckoutServiceImpl{
		DB:            db,
		CartRepo:      cartRepo,
		ProductRepo:   productRepo,
		AddressRepo:   addressRepo,
		OrderRepo:     orderRepo,
		OrderItemRepo: orderItemRepo,
		PaymentRepo:   paymentRepo,
	}

	address := &domain.Address{
		ID:            1,
		UserID:        1,
		RecipientName: "John Doe",
		Phone:         "08123456789",
		Address:       "Jl. Malioboro",
		City:          "Yogyakarta",
		PostalCode:    "55213",
	}

	addressRepo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(address, nil)

	carts := []domain.Cart{
		{
			ID:        1,
			UserID:    1,
			ProductID: 1,
			Quantity:  5,
			Product: domain.Product{
				ID:    1,
				Name:  "Keyboard",
				Price: 750000,
				Stock: 2,
			},
		},
	}

	cartRepo.On(
		"FindAllByUserID",
		mock.Anything,
		uint64(1),
	).Return(carts, nil)

	result, err := service.Checkout(
		context.Background(),
		1,
		request.CheckoutRequest{
			AddressID:     1,
			PaymentMethod: "BANK_TRANSFER",
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"Keyboard stock is insufficient",
	)

	addressRepo.AssertExpectations(t)
	cartRepo.AssertExpectations(t)
}

func TestCheckoutCreateOrderFailed(t *testing.T) {

	db, sqlMock := testutils.NewMockDB(t)

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)
	addressRepo := new(MockAddressRepository)
	orderRepo := new(MockOrderRepository)
	orderItemRepo := new(MockOrderItemRepository)
	paymentRepo := new(MockPaymentRepository)

	service := CheckoutServiceImpl{
		DB:            db,
		CartRepo:      cartRepo,
		ProductRepo:   productRepo,
		AddressRepo:   addressRepo,
		OrderRepo:     orderRepo,
		OrderItemRepo: orderItemRepo,
		PaymentRepo:   paymentRepo,
	}

	address := &domain.Address{
		ID:            1,
		UserID:        1,
		RecipientName: "John Doe",
		Phone:         "08123456789",
		Address:       "Jl. Malioboro",
		City:          "Yogyakarta",
		PostalCode:    "55213",
	}

	carts := []domain.Cart{
		{
			ID:        1,
			UserID:    1,
			ProductID: 1,
			Quantity:  2,
			Product: domain.Product{
				ID:    1,
				Name:  "Keyboard",
				Price: 750000,
				Stock: 10,
			},
		},
	}

	addressRepo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(address, nil)

	cartRepo.On(
		"FindAllByUserID",
		mock.Anything,
		uint64(1),
	).Return(carts, nil)

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	orderRepo.On(
		"CreateTx",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*domain.Order"),
	).Return(errors.New("database error"))

	result, err := service.Checkout(
		context.Background(),
		1,
		request.CheckoutRequest{
			AddressID:     1,
			PaymentMethod: "BANK_TRANSFER",
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(t, err.Error(), "database error")

	assert.NoError(t, sqlMock.ExpectationsWereMet())

	addressRepo.AssertExpectations(t)
	cartRepo.AssertExpectations(t)
	orderRepo.AssertExpectations(t)
}

func TestCheckoutCreateOrderItemFailed(t *testing.T) {

	db, sqlMock := testutils.NewMockDB(t)

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)
	addressRepo := new(MockAddressRepository)
	orderRepo := new(MockOrderRepository)
	orderItemRepo := new(MockOrderItemRepository)
	paymentRepo := new(MockPaymentRepository)

	service := CheckoutServiceImpl{
		DB:            db,
		CartRepo:      cartRepo,
		ProductRepo:   productRepo,
		AddressRepo:   addressRepo,
		OrderRepo:     orderRepo,
		OrderItemRepo: orderItemRepo,
		PaymentRepo:   paymentRepo,
	}

	address := &domain.Address{
		ID:            1,
		UserID:        1,
		RecipientName: "John Doe",
		Phone:         "08123456789",
		Address:       "Jl. Malioboro",
		City:          "Yogyakarta",
		PostalCode:    "55213",
	}

	carts := []domain.Cart{
		{
			ID:        1,
			UserID:    1,
			ProductID: 1,
			Quantity:  2,
			Product: domain.Product{
				ID:    1,
				Name:  "Keyboard",
				Price: 750000,
				Stock: 10,
			},
		},
	}

	addressRepo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(address, nil)

	cartRepo.On(
		"FindAllByUserID",
		mock.Anything,
		uint64(1),
	).Return(carts, nil)

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	orderRepo.On(
		"CreateTx",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*domain.Order"),
	).Run(func(args mock.Arguments) {
		order := args.Get(2).(*domain.Order)
		order.ID = 1
	}).Return(nil)

	orderItemRepo.On(
		"CreateTx",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*domain.OrderItem"),
	).Return(errors.New("database error"))

	result, err := service.Checkout(
		context.Background(),
		1,
		request.CheckoutRequest{
			AddressID:     1,
			PaymentMethod: "BANK_TRANSFER",
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(t, err.Error(), "database error")

	assert.NoError(t, sqlMock.ExpectationsWereMet())

	addressRepo.AssertExpectations(t)
	cartRepo.AssertExpectations(t)
	orderRepo.AssertExpectations(t)
	orderItemRepo.AssertExpectations(t)
}

func TestCheckoutUpdateProductFailed(t *testing.T) {

	db, sqlMock := testutils.NewMockDB(t)

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)
	addressRepo := new(MockAddressRepository)
	orderRepo := new(MockOrderRepository)
	orderItemRepo := new(MockOrderItemRepository)
	paymentRepo := new(MockPaymentRepository)

	service := CheckoutServiceImpl{
		DB:            db,
		CartRepo:      cartRepo,
		ProductRepo:   productRepo,
		AddressRepo:   addressRepo,
		OrderRepo:     orderRepo,
		OrderItemRepo: orderItemRepo,
		PaymentRepo:   paymentRepo,
	}

	address := &domain.Address{
		ID:            1,
		UserID:        1,
		RecipientName: "John Doe",
		Phone:         "08123456789",
		Address:       "Jl. Malioboro",
		City:          "Yogyakarta",
		PostalCode:    "55213",
	}

	carts := []domain.Cart{
		{
			ID:        1,
			UserID:    1,
			ProductID: 1,
			Quantity:  2,
			Product: domain.Product{
				ID:    1,
				Name:  "Keyboard",
				Price: 750000,
				Stock: 10,
			},
		},
	}

	addressRepo.On("FindByIDAndUserID", mock.Anything, uint64(1), uint64(1)).
		Return(address, nil)

	cartRepo.On("FindAllByUserID", mock.Anything, uint64(1)).
		Return(carts, nil)

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	orderRepo.On(
		"CreateTx",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*domain.Order"),
	).Run(func(args mock.Arguments) {
		order := args.Get(2).(*domain.Order)
		order.ID = 1
	}).Return(nil)

	orderItemRepo.On(
		"CreateTx",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*domain.OrderItem"),
	).Return(nil)

	productRepo.On(
		"UpdateTx",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*domain.Product"),
	).Return(errors.New("database error"))

	result, err := service.Checkout(
		context.Background(),
		1,
		request.CheckoutRequest{
			AddressID:     1,
			PaymentMethod: "BANK_TRANSFER",
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")

	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestCheckoutCreatePaymentFailed(t *testing.T) {

	db, sqlMock := testutils.NewMockDB(t)

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)
	addressRepo := new(MockAddressRepository)
	orderRepo := new(MockOrderRepository)
	orderItemRepo := new(MockOrderItemRepository)
	paymentRepo := new(MockPaymentRepository)

	service := CheckoutServiceImpl{
		DB:            db,
		CartRepo:      cartRepo,
		ProductRepo:   productRepo,
		AddressRepo:   addressRepo,
		OrderRepo:     orderRepo,
		OrderItemRepo: orderItemRepo,
		PaymentRepo:   paymentRepo,
	}

	address := &domain.Address{
		ID:            1,
		UserID:        1,
		RecipientName: "John Doe",
		Phone:         "08123456789",
		Address:       "Jl. Malioboro",
		City:          "Yogyakarta",
		PostalCode:    "55213",
	}

	carts := []domain.Cart{
		{
			ID:        1,
			UserID:    1,
			ProductID: 1,
			Quantity:  2,
			Product: domain.Product{
				ID:    1,
				Name:  "Keyboard",
				Price: 750000,
				Stock: 10,
			},
		},
	}

	addressRepo.On("FindByIDAndUserID", mock.Anything, uint64(1), uint64(1)).
		Return(address, nil)

	cartRepo.On("FindAllByUserID", mock.Anything, uint64(1)).
		Return(carts, nil)

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	orderRepo.On(
		"CreateTx",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*domain.Order"),
	).Run(func(args mock.Arguments) {
		order := args.Get(2).(*domain.Order)
		order.ID = 1
	}).Return(nil)

	orderItemRepo.On(
		"CreateTx",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*domain.OrderItem"),
	).Return(nil)

	productRepo.On(
		"UpdateTx",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*domain.Product"),
	).Return(nil)

	paymentRepo.On(
		"CreateTx",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*domain.Payment"),
	).Return(errors.New("database error"))

	result, err := service.Checkout(
		context.Background(),
		1,
		request.CheckoutRequest{
			AddressID:     1,
			PaymentMethod: "BANK_TRANSFER",
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(t, err.Error(), "database error")

	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestCheckoutDeleteCartFailed(t *testing.T) {

	db, sqlMock := testutils.NewMockDB(t)

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)
	addressRepo := new(MockAddressRepository)
	orderRepo := new(MockOrderRepository)
	orderItemRepo := new(MockOrderItemRepository)
	paymentRepo := new(MockPaymentRepository)

	service := CheckoutServiceImpl{
		DB:            db,
		CartRepo:      cartRepo,
		ProductRepo:   productRepo,
		AddressRepo:   addressRepo,
		OrderRepo:     orderRepo,
		OrderItemRepo: orderItemRepo,
		PaymentRepo:   paymentRepo,
	}

	address := &domain.Address{
		ID:            1,
		UserID:        1,
		RecipientName: "John Doe",
		Phone:         "08123456789",
		Address:       "Jl. Malioboro",
		City:          "Yogyakarta",
		PostalCode:    "55213",
	}

	carts := []domain.Cart{
		{
			ID:        1,
			UserID:    1,
			ProductID: 1,
			Quantity:  2,
			Product: domain.Product{
				ID:    1,
				Name:  "Keyboard",
				Price: 750000,
				Stock: 10,
			},
		},
	}

	addressRepo.On("FindByIDAndUserID", mock.Anything, uint64(1), uint64(1)).
		Return(address, nil)

	cartRepo.On("FindAllByUserID", mock.Anything, uint64(1)).
		Return(carts, nil)

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	orderRepo.On(
		"CreateTx",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*domain.Order"),
	).Run(func(args mock.Arguments) {
		order := args.Get(2).(*domain.Order)
		order.ID = 1
	}).Return(nil)

	orderItemRepo.On(
		"CreateTx",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*domain.OrderItem"),
	).Return(nil)

	productRepo.On(
		"UpdateTx",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*domain.Product"),
	).Return(nil)

	paymentRepo.On(
		"CreateTx",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*domain.Payment"),
	).Return(nil)

	cartRepo.On(
		"DeleteAllByUserIDTx",
		mock.Anything,
		mock.Anything,
		uint64(1),
	).Return(errors.New("database error"))

	result, err := service.Checkout(
		context.Background(),
		1,
		request.CheckoutRequest{
			AddressID:     1,
			PaymentMethod: "BANK_TRANSFER",
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(t, err.Error(), "database error")

	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestCheckoutRepositoryError(t *testing.T) {

	db, _ := testutils.NewMockDB(t)

	cartRepo := new(MockCartRepository)
	productRepo := new(MockProductRepository)
	addressRepo := new(MockAddressRepository)
	orderRepo := new(MockOrderRepository)
	orderItemRepo := new(MockOrderItemRepository)
	paymentRepo := new(MockPaymentRepository)

	service := CheckoutServiceImpl{
		DB:            db,
		CartRepo:      cartRepo,
		ProductRepo:   productRepo,
		AddressRepo:   addressRepo,
		OrderRepo:     orderRepo,
		OrderItemRepo: orderItemRepo,
		PaymentRepo:   paymentRepo,
	}

	address := &domain.Address{
		ID:            1,
		UserID:        1,
		RecipientName: "John Doe",
		Phone:         "08123456789",
		Address:       "Jl. Malioboro",
		City:          "Yogyakarta",
		PostalCode:    "55213",
	}

	addressRepo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(address, nil)

	cartRepo.On(
		"FindAllByUserID",
		mock.Anything,
		uint64(1),
	).Return(
		nil,
		errors.New("database error"),
	)

	result, err := service.Checkout(
		context.Background(),
		1,
		request.CheckoutRequest{
			AddressID:     1,
			PaymentMethod: "BANK_TRANSFER",
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"database error",
	)

	addressRepo.AssertExpectations(t)
	cartRepo.AssertExpectations(t)
}
