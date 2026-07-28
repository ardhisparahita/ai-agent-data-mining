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

type AddressServiceImpl struct {
	Repo repository.AddressRepository
}

func NewAddressService(repo repository.AddressRepository) AddressService {
	return &AddressServiceImpl{
		Repo: repo,
	}
}

func (s *AddressServiceImpl) Create(ctx context.Context, userID uint64, req request.CreateAddressRequest) (*response.AddressResponse, error) {
	address := domain.Address{
		UserID:        userID,
		RecipientName: req.RecipientName,
		Phone:         req.Phone,
		Address:       req.Address,
		City:          req.City,
		PostalCode:    req.PostalCode,
	}

	err := s.Repo.Create(ctx, &address)
	if err != nil {
		return nil, err
	}

	return mapper.ToAddressResponse(&address), nil
}

func (s *AddressServiceImpl) FindAllByUserID(ctx context.Context, userID uint64) ([]response.AddressResponse, error) {
	addresses, err := s.Repo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return mapper.ToAddressResponses(addresses), nil
}

func (s *AddressServiceImpl) FindByID(ctx context.Context, id uint64, userID uint64) (*response.AddressResponse, error) {
	address, err := s.Repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("address not found")
		}
		return nil, err
	}

	return mapper.ToAddressResponse(address), nil
}

func (s *AddressServiceImpl) Update(ctx context.Context, id uint64, userID uint64, req request.UpdateAddressRequest) (*response.AddressResponse, error) {
	address, err := s.Repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("address not found")
		}
		return nil, err
	}

	address.RecipientName = req.RecipientName
	address.Phone = req.Phone
	address.Address = req.Address
	address.City = req.City
	address.PostalCode = req.PostalCode

	err = s.Repo.Update(ctx, address)
	if err != nil {
		return nil, err
	}

	return mapper.ToAddressResponse(address), nil
}

func (s *AddressServiceImpl) Delete(ctx context.Context, id uint64, userID uint64) error {
	_, err := s.Repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("address not found")
		}
		return err
	}

	return s.Repo.Delete(ctx, id)
}
