package service

import (
	"context"
	"errors"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
	"github.com/ardhisparahita/ecommerce-api/internal/mapper"
	"github.com/ardhisparahita/ecommerce-api/internal/repository"
	"github.com/ardhisparahita/ecommerce-api/pkg/utils"
	"gorm.io/gorm"
)

type CartServiceImpl struct {
	CartRepo    repository.CartRepository
	ProductRepo repository.ProductRepository
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository) CartService {
	return &CartServiceImpl{
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}
}

func (s *CartServiceImpl) AddToCart(ctx context.Context, userID uint64, req request.AddToCartRequest) (*response.CartResponse, error) {
	product, err := s.ProductRepo.FindByID(ctx, req.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("product not found")
		}

		return nil, err
	}

	if product.Stock < req.Quantity {
		return nil, utils.BadRequest("insufficient stock")
	}

	cart, err := s.CartRepo.FindByUserIDAndProductID(ctx, userID, req.ProductID)

	if err == nil {
		newQuantity := cart.Quantity + req.Quantity

		if product.Stock < newQuantity {
			return nil, utils.BadRequest("insufficient stock")
		}

		cart.Quantity = newQuantity

		if err := s.CartRepo.Update(ctx, cart); err != nil {
			return nil, err
		}

		cart, err := s.CartRepo.FindByIDAndUserID(ctx, cart.ID, userID)
		if err != nil {
			return nil, err
		}

		return mapper.ToCartResponse(cart), nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	newCart := domain.Cart{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := s.CartRepo.Create(ctx, &newCart); err != nil {
		return nil, err
	}

	cart, err = s.CartRepo.FindByIDAndUserID(ctx, newCart.ID, userID)
	if err != nil {
		return nil, err
	}

	return mapper.ToCartResponse(cart), nil
}

func (s *CartServiceImpl) FindAll(ctx context.Context, userID uint64) (*response.CartListResponse, error) {
	carts, err := s.CartRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	items, grandTotal := mapper.ToCartResponses(carts)

	return &response.CartListResponse{
		Items:      items,
		TotalItems: len(items),
		GrandTotal: grandTotal,
	}, nil
}

func (s *CartServiceImpl) Update(ctx context.Context, id uint64, userID uint64, req request.UpdateCartRequest) (*response.CartResponse, error) {
	if req.Quantity <= 0 {
		return nil, utils.BadRequest("quantity must be greater than 0")
	}

	cart, err := s.CartRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("cart not found")
		}
		return nil, err
	}

	if cart.Product.Stock < req.Quantity {
		return nil, utils.BadRequest("insufficient stock")
	}

	cart.Quantity = req.Quantity

	if err := s.CartRepo.Update(ctx, cart); err != nil {
		return nil, err
	}

	cart, err = s.CartRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return mapper.ToCartResponse(cart), nil

}

func (s *CartServiceImpl) Delete(ctx context.Context, id uint64, userID uint64) error {
	_, err := s.CartRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("cart not found")
		}
		return err
	}

	return s.CartRepo.Delete(ctx, id, userID)
}
