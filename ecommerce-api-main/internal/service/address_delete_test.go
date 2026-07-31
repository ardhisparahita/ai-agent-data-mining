package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestAddressDeleteSuccess(t *testing.T) {

	repo := new(MockAddressRepository)

	service := AddressServiceImpl{
		Repo: repo,
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
		"Delete",
		mock.Anything,
		uint64(1),
	).Return(nil)

	err := service.Delete(
		context.Background(),
		1,
		1,
	)

	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestAddressDeleteNotFound(t *testing.T) {

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

	err := service.Delete(
		context.Background(),
		99,
		1,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"address not found",
	)

	repo.AssertExpectations(t)
}

func TestAddressDeleteRepositoryError(t *testing.T) {

	repo := new(MockAddressRepository)

	service := AddressServiceImpl{
		Repo: repo,
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
		"Delete",
		mock.Anything,
		uint64(1),
	).Return(
		errors.New("database error"),
	)

	err := service.Delete(
		context.Background(),
		1,
		1,
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"database error",
	)

	repo.AssertExpectations(t)
}
