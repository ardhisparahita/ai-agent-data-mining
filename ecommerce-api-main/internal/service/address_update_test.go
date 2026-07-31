package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestAddressUpdateSuccess(t *testing.T) {

	repo := new(MockAddressRepository)

	service := AddressServiceImpl{
		Repo: repo,
	}

	address := &domain.Address{
		ID:            1,
		UserID:        1,
		RecipientName: "Old Name",
		Phone:         "081111111111",
		Address:       "Old Address",
		City:          "Old City",
		PostalCode:    "11111",
	}

	repo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(
		address,
		nil,
	)

	repo.On(
		"Update",
		mock.Anything,
		address,
	).Return(nil)

	result, err := service.Update(
		context.Background(),
		1,
		1,
		request.UpdateAddressRequest{
			RecipientName: "John Doe",
			Phone:         "08123456789",
			Address:       "Jl. Malioboro",
			City:          "Yogyakarta",
			PostalCode:    "55213",
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, uint64(1), result.ID)
	assert.Equal(t, "John Doe", result.RecipientName)
	assert.Equal(t, "08123456789", result.Phone)
	assert.Equal(t, "Jl. Malioboro", result.Address)
	assert.Equal(t, "Yogyakarta", result.City)
	assert.Equal(t, "55213", result.PostalCode)

	repo.AssertExpectations(t)
}

func TestAddressUpdateNotFound(t *testing.T) {

	repo := new(MockAddressRepository)

	service := AddressServiceImpl{
		Repo: repo,
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

	result, err := service.Update(
		context.Background(),
		99,
		1,
		request.UpdateAddressRequest{
			RecipientName: "John Doe",
			Phone:         "08123456789",
			Address:       "Jl. Malioboro",
			City:          "Yogyakarta",
			PostalCode:    "55213",
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"address not found",
	)

	repo.AssertExpectations(t)
}

func TestAddressUpdateRepositoryError(t *testing.T) {

	repo := new(MockAddressRepository)

	service := AddressServiceImpl{
		Repo: repo,
	}

	address := &domain.Address{
		ID:            1,
		UserID:        1,
		RecipientName: "Old Name",
		Phone:         "081111111111",
		Address:       "Old Address",
		City:          "Old City",
		PostalCode:    "11111",
	}

	repo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(
		address,
		nil,
	)

	repo.On(
		"Update",
		mock.Anything,
		address,
	).Return(
		errors.New("database error"),
	)

	result, err := service.Update(
		context.Background(),
		1,
		1,
		request.UpdateAddressRequest{
			RecipientName: "John Doe",
			Phone:         "08123456789",
			Address:       "Jl. Malioboro",
			City:          "Yogyakarta",
			PostalCode:    "55213",
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"database error",
	)

	repo.AssertExpectations(t)
}
