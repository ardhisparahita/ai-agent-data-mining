package handler

import (
	"context"

	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"

	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(
	ctx context.Context,
	req request.RegisterRequest,
) (*response.UserResponse, error) {

	args := m.Called(
		ctx,
		req,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*response.UserResponse), args.Error(1)
}

func (m *MockAuthService) Login(
	ctx context.Context,
	req request.LoginRequest,
) (*response.AuthResponse, error) {

	args := m.Called(
		ctx,
		req,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*response.AuthResponse), args.Error(1)
}

func (m *MockAuthService) GetProfile(
	ctx context.Context,
	userID uint64,
) (*response.UserResponse, error) {

	args := m.Called(
		ctx,
		userID,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*response.UserResponse), args.Error(1)
}

func (m *MockAuthService) UpdateProfile(
	ctx context.Context,
	userID uint64,
	req request.UpdateProfileRequest,
) (*response.UserResponse, error) {

	args := m.Called(
		ctx,
		userID,
		req,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*response.UserResponse), args.Error(1)
}

func (m *MockAuthService) ChangePassword(
	ctx context.Context,
	userID uint64,
	req request.ChangePasswordRequest,
) error {

	args := m.Called(
		ctx,
		userID,
		req,
	)

	return args.Error(0)
}

type MockCategoryService struct {
	mock.Mock
}

func (m *MockCategoryService) Create(
	ctx context.Context,
	req request.CreateCategoryRequest,
) (*response.CategoryResponse, error) {

	args := m.Called(
		ctx,
		req,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*response.CategoryResponse), args.Error(1)
}

func (m *MockCategoryService) FindAll(
	ctx context.Context,
) ([]response.CategoryResponse, error) {

	args := m.Called(
		ctx,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]response.CategoryResponse), args.Error(1)
}

func (m *MockCategoryService) FindByID(
	ctx context.Context,
	id uint64,
) (*response.CategoryResponse, error) {

	args := m.Called(
		ctx,
		id,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*response.CategoryResponse), args.Error(1)
}

func (m *MockCategoryService) Update(
	ctx context.Context,
	id uint64,
	req request.UpdateCategoryRequest,
) (*response.CategoryResponse, error) {

	args := m.Called(
		ctx,
		id,
		req,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*response.CategoryResponse), args.Error(1)
}

func (m *MockCategoryService) Delete(
	ctx context.Context,
	id uint64,
) error {

	args := m.Called(
		ctx,
		id,
	)

	return args.Error(0)
}

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) Create(
	ctx context.Context,
	req request.CreateProductRequest,
) (*response.ProductResponse, error) {

	args := m.Called(
		ctx,
		req,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*response.ProductResponse), args.Error(1)
}

func (m *MockProductService) FindAll(
	ctx context.Context,
	req request.ProductQueryRequest,
) (*response.ProductListResponse, error) {

	args := m.Called(
		ctx,
		req,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*response.ProductListResponse), args.Error(1)
}

func (m *MockProductService) FindByID(
	ctx context.Context,
	id uint64,
) (*response.ProductResponse, error) {

	args := m.Called(
		ctx,
		id,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*response.ProductResponse), args.Error(1)
}

func (m *MockProductService) Update(
	ctx context.Context,
	id uint64,
	req request.UpdateProductRequest,
) (*response.ProductResponse, error) {

	args := m.Called(
		ctx,
		id,
		req,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*response.ProductResponse), args.Error(1)
}

func (m *MockProductService) Delete(
	ctx context.Context,
	id uint64,
) error {

	args := m.Called(
		ctx,
		id,
	)

	return args.Error(0)
}

func (m *MockProductService) UploadImage(
	ctx context.Context,
	id uint64,
	imageURL string,
) (*response.ProductResponse, error) {

	args := m.Called(
		ctx,
		id,
		imageURL,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*response.ProductResponse), args.Error(1)
}

type MockAddressService struct {
	mock.Mock
}

func (m *MockAddressService) Create(
	ctx context.Context,
	userID uint64,
	req request.CreateAddressRequest,
) (*response.AddressResponse, error) {

	args := m.Called(
		ctx,
		userID,
		req,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*response.AddressResponse), args.Error(1)
}

func (m *MockAddressService) FindAllByUserID(
	ctx context.Context,
	userID uint64,
) ([]response.AddressResponse, error) {

	args := m.Called(
		ctx,
		userID,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]response.AddressResponse), args.Error(1)
}

func (m *MockAddressService) FindByID(
	ctx context.Context,
	id uint64,
	userID uint64,
) (*response.AddressResponse, error) {

	args := m.Called(
		ctx,
		id,
		userID,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*response.AddressResponse), args.Error(1)
}

func (m *MockAddressService) Update(
	ctx context.Context,
	id uint64,
	userID uint64,
	req request.UpdateAddressRequest,
) (*response.AddressResponse, error) {

	args := m.Called(
		ctx,
		id,
		userID,
		req,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*response.AddressResponse), args.Error(1)
}

func (m *MockAddressService) Delete(
	ctx context.Context,
	id uint64,
	userID uint64,
) error {

	args := m.Called(
		ctx,
		id,
		userID,
	)

	return args.Error(0)
}

type MockCartService struct {
	mock.Mock
}

func (m *MockCartService) AddToCart(
	ctx context.Context,
	userID uint64,
	req request.AddToCartRequest,
) (*response.CartResponse, error) {

	args := m.Called(
		ctx,
		userID,
		req,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*response.CartResponse), args.Error(1)
}

func (m *MockCartService) FindAll(
	ctx context.Context,
	userID uint64,
) (*response.CartListResponse, error) {

	args := m.Called(
		ctx,
		userID,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*response.CartListResponse), args.Error(1)
}

func (m *MockCartService) Update(
	ctx context.Context,
	id uint64,
	userID uint64,
	req request.UpdateCartRequest,
) (*response.CartResponse, error) {

	args := m.Called(
		ctx,
		id,
		userID,
		req,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*response.CartResponse), args.Error(1)
}

func (m *MockCartService) Delete(
	ctx context.Context,
	id uint64,
	userID uint64,
) error {

	args := m.Called(
		ctx,
		id,
		userID,
	)

	return args.Error(0)
}

type MockCheckoutService struct {
	mock.Mock
}

func (m *MockCheckoutService) Checkout(
	ctx context.Context,
	userID uint64,
	req request.CheckoutRequest,
) (*response.OrderResponse, error) {

	args := m.Called(
		ctx,
		userID,
		req,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*response.OrderResponse), args.Error(1)
}

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) FindAll(
	ctx context.Context,
	userID uint64,
) ([]response.OrderListResponse, error) {

	args := m.Called(
		ctx,
		userID,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]response.OrderListResponse), args.Error(1)
}

func (m *MockOrderService) FindByID(
	ctx context.Context,
	id uint64,
	userID uint64,
) (*response.OrderDetailResponse, error) {

	args := m.Called(
		ctx,
		id,
		userID,
	)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*response.OrderDetailResponse), args.Error(1)
}

func (m *MockOrderService) Cancel(
	ctx context.Context,
	id uint64,
	userID uint64,
) error {

	args := m.Called(
		ctx,
		id,
		userID,
	)

	return args.Error(0)
}

func (m *MockOrderService) MarkAsPaid(
	ctx context.Context,
	id uint64,
) error {

	args := m.Called(
		ctx,
		id,
	)

	return args.Error(0)
}

func (m *MockOrderService) MarkAsFailed(
	ctx context.Context,
	id uint64,
) error {

	args := m.Called(
		ctx,
		id,
	)

	return args.Error(0)
}

func (m *MockOrderService) MarkAsShipped(
	ctx context.Context,
	id uint64,
) error {

	args := m.Called(
		ctx,
		id,
	)

	return args.Error(0)
}

func (m *MockOrderService) MarkAsCompleted(
	ctx context.Context,
	id uint64,
	userID uint64,
	role string,
) error {
	args := m.Called(ctx, id, userID, role)
	return args.Error(0)
}
