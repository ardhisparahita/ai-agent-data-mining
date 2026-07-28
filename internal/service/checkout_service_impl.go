package service

import (
	"context"
	"errors"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
	"github.com/ardhisparahita/ecommerce-api/internal/repository"
	"github.com/ardhisparahita/ecommerce-api/pkg/utils"
	"gorm.io/gorm"
)

type CheckoutServiceImpl struct {
	DB            *gorm.DB
	CartRepo      repository.CartRepository
	ProductRepo   repository.ProductRepository
	AddressRepo   repository.AddressRepository
	OrderRepo     repository.OrderRepository
	OrderItemRepo repository.OrderItemRepository
	PaymentRepo   repository.PaymentRepository
}

func NewCheckoutService(
	dB *gorm.DB,
	cartRepo repository.CartRepository,
	productRepo repository.ProductRepository,
	addressRepo repository.AddressRepository,
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	paymentRepo repository.PaymentRepository,
) CheckoutService {
	return &CheckoutServiceImpl{
		DB:            dB,
		CartRepo:      cartRepo,
		ProductRepo:   productRepo,
		AddressRepo:   addressRepo,
		OrderRepo:     orderRepo,
		OrderItemRepo: orderItemRepo,
		PaymentRepo:   paymentRepo,
	}
}

func (s *CheckoutServiceImpl) Checkout(ctx context.Context, userID uint64, req request.CheckoutRequest) (*response.OrderResponse, error) {
	address, err := s.AddressRepo.FindByIDAndUserID(ctx, req.AddressID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("address not found")
		}
		return nil, err
	}

	carts, err := s.CartRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(carts) == 0 {
		return nil, utils.BadRequest("cart is empty")
	}

	var grandTotal float64

	for _, cart := range carts {
		if cart.Product.Stock < cart.Quantity {
			return nil, utils.BadRequest(cart.Product.Name + " stock is insufficient")
		}
		grandTotal += cart.Product.Price * float64(cart.Quantity)
	}

	var order domain.Order

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		order = domain.Order{
			UserID:        userID,
			AddressID:     address.ID,
			RecipientName: address.RecipientName,
			Phone:         address.Phone,
			Address:       address.Address,
			City:          address.City,
			PostalCode:    address.PostalCode,
			TotalAmount:   grandTotal,
			Status:        "PENDING",
		}
		if err := s.OrderRepo.CreateTx(ctx, tx, &order); err != nil {
			return err
		}

		for _, cart := range carts {
			orderItem := domain.OrderItem{
				OrderID:     order.ID,
				ProductID:   cart.ProductID,
				ProductName: cart.Product.Name,
				Price:       cart.Product.Price,
				Quantity:    cart.Quantity,
				Subtotal:    cart.Product.Price * float64(cart.Quantity),
			}
			if err := s.OrderItemRepo.CreateTx(ctx, tx, &orderItem); err != nil {
				return err
			}

			cart.Product.Stock -= cart.Quantity

			if err := s.ProductRepo.UpdateTx(ctx, tx, &cart.Product); err != nil {
				return err
			}
		}

		payment := domain.Payment{
			OrderID: order.ID,
			Method:  req.PaymentMethod,
			Status:  "PENDING",
			Amount:  grandTotal,
		}

		if err := s.PaymentRepo.CreateTx(ctx, tx, &payment); err != nil {
			return err
		}

		if err := s.CartRepo.DeleteAllByUserIDTx(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &response.OrderResponse{
		ID:            order.ID,
		TotalAmount:   order.TotalAmount,
		Status:        order.Status,
		PaymentMethod: req.PaymentMethod,
	}, nil
}
