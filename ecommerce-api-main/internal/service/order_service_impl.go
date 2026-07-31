package service

import (
	"context"
	"errors"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
	"github.com/ardhisparahita/ecommerce-api/internal/mapper"
	"github.com/ardhisparahita/ecommerce-api/internal/repository"
	"github.com/ardhisparahita/ecommerce-api/pkg/utils"
	"gorm.io/gorm"
)

type OrderServiceImpl struct {
	DB          *gorm.DB
	OrderRepo   repository.OrderRepository
	ProductRepo repository.ProductRepository
	PaymentRepo repository.PaymentRepository
}

func NewOrderService(
	db *gorm.DB,
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	paymentRepo repository.PaymentRepository,
) OrderService {
	return &OrderServiceImpl{
		DB:          db,
		OrderRepo:   orderRepo,
		ProductRepo: productRepo,
		PaymentRepo: paymentRepo,
	}
}

func (s *OrderServiceImpl) FindAll(ctx context.Context, userID uint64) ([]response.OrderListResponse, error) {
	orders, err := s.OrderRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return mapper.ToOrderListResponse(orders), nil
}

func (s *OrderServiceImpl) FindByID(ctx context.Context, id uint64, userID uint64) (*response.OrderDetailResponse, error) {
	order, err := s.OrderRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("order not found")
		}
		return nil, err
	}

	return mapper.ToOrderDetailResponse(order), nil
}

func (s *OrderServiceImpl) MarkAsPaid(ctx context.Context, id uint64) error {
	order, err := s.OrderRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("order not found")
		}

		return err
	}

	switch order.Status {

	case domain.OrderPending:

	case domain.OrderPaid:
		return utils.BadRequest("order already paid")

	case domain.OrderShipped:
		return utils.BadRequest("order already shipped")

	case domain.OrderCompleted:
		return utils.BadRequest("order already completed")

	case domain.OrderCancelled:
		return utils.BadRequest("cancelled order cannot be paid")

	default:
		return utils.BadRequest("invalid order status")
	}

	payment, err := s.PaymentRepo.FIndByOrderID(ctx, order.ID)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("payment not found")
		}

		return err
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {

		payment.Status = domain.PaymentPaid

		if err := s.PaymentRepo.UpdateTx(
			ctx,
			tx,
			payment,
		); err != nil {
			return err
		}

		order.Status = domain.OrderPaid

		if err := s.OrderRepo.UpdateTx(
			ctx,
			tx,
			order,
		); err != nil {
			return err
		}

		return nil
	})
}

func (s *OrderServiceImpl) MarkAsFailed(ctx context.Context, id uint64) error {

	order, err := s.OrderRepo.FindByIDWithItems(ctx, id)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("order not found")
		}

		return err
	}

	switch order.Status {

	case domain.OrderPending:

	case domain.OrderPaid:
		return utils.BadRequest("paid order cannot be marked as failed")

	case domain.OrderShipped:
		return utils.BadRequest("shipped order cannot be marked as failed")

	case domain.OrderCompleted:
		return utils.BadRequest("completed order cannot be marked as failed")

	case domain.OrderCancelled:
		return utils.BadRequest("order already cancelled")

	default:
		return utils.BadRequest("invalid order status")
	}

	payment, err := s.PaymentRepo.FIndByOrderID(ctx, order.ID)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("payment not found")
		}

		return err
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {

		for _, item := range order.OrderItems {

			product, err := s.ProductRepo.FindByID(ctx, item.ProductID)
			if err != nil {
				return err
			}

			product.Stock += item.Quantity

			if err := s.ProductRepo.UpdateTx(
				ctx,
				tx,
				product,
			); err != nil {
				return err
			}
		}

		payment.Status = domain.PaymentFailed

		if err := s.PaymentRepo.UpdateTx(
			ctx,
			tx,
			payment,
		); err != nil {
			return err
		}

		order.Status = domain.OrderCancelled

		if err := s.OrderRepo.UpdateTx(
			ctx,
			tx,
			order,
		); err != nil {
			return err
		}

		return nil
	})
}

func (s *OrderServiceImpl) Cancel(ctx context.Context, id uint64, userID uint64) error {
	order, err := s.OrderRepo.FindByIDAndUserIDWithItems(
		ctx,
		id,
		userID,
	)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("order not found")
		}

		return err
	}

	switch order.Status {

	case domain.OrderPending:

	case domain.OrderCancelled:
		return utils.BadRequest("order already cancelled")

	case domain.OrderPaid:
		return utils.BadRequest("paid order cannot be cancelled")

	case domain.OrderShipped:
		return utils.BadRequest("shipped order cannot be cancelled")

	case domain.OrderCompleted:
		return utils.BadRequest("completed order cannot be cancelled")

	default:
		return utils.BadRequest("invalid order status")
	}

	payment, err := s.PaymentRepo.FIndByOrderID(ctx, order.ID)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("payment not found")
		}

		return err
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {

		for _, item := range order.OrderItems {

			product, err := s.ProductRepo.FindByID(
				ctx,
				item.ProductID,
			)
			if err != nil {
				return err
			}

			product.Stock += item.Quantity

			if err := s.ProductRepo.UpdateTx(
				ctx,
				tx,
				product,
			); err != nil {
				return err
			}
		}

		payment.Status = domain.PaymentFailed

		if err := s.PaymentRepo.UpdateTx(
			ctx,
			tx,
			payment,
		); err != nil {
			return err
		}

		order.Status = domain.OrderCancelled

		if err := s.OrderRepo.UpdateTx(
			ctx,
			tx,
			order,
		); err != nil {
			return err
		}

		return nil
	})
}

func (s *OrderServiceImpl) MarkAsShipped(ctx context.Context, id uint64) error {

	order, err := s.OrderRepo.FindByID(ctx, id)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("order not found")
		}

		return err
	}

	switch order.Status {

	case domain.OrderPaid:

	case domain.OrderPending:
		return utils.BadRequest("order must be paid before shipping")

	case domain.OrderShipped:
		return utils.BadRequest("order already shipped")

	case domain.OrderCompleted:
		return utils.BadRequest("completed order cannot be shipped")

	case domain.OrderCancelled:
		return utils.BadRequest("cancelled order cannot be shipped")

	default:
		return utils.BadRequest("invalid order status")
	}

	order.Status = domain.OrderShipped

	if err := s.OrderRepo.Update(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *OrderServiceImpl) MarkAsCompleted(ctx context.Context, id uint64, userID uint64, role string) error {
	order, err := s.OrderRepo.FindByID(ctx, id)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("order not found")
		}

		return err
	}

	if role != domain.RoleAdmin && order.UserID != userID {
		return utils.Forbidden(
			"you are not allowed to complete this order",
		)
	}

	switch order.Status {

	case domain.OrderShipped:

	case domain.OrderPending:
		return utils.BadRequest(
			"order must be shipped first",
		)

	case domain.OrderPaid:
		return utils.BadRequest(
			"order must be shipped first",
		)

	case domain.OrderCompleted:
		return utils.BadRequest(
			"order already completed",
		)

	case domain.OrderCancelled:
		return utils.BadRequest(
			"cancelled order cannot be completed",
		)

	default:
		return utils.BadRequest(
			"invalid order status",
		)
	}

	order.Status = domain.OrderCompleted

	if err := s.OrderRepo.Update(ctx, order); err != nil {
		return err
	}

	return nil
}
