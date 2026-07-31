package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddressCreateSuccess(t *testing.T) {

	repo := new(MockAddressRepository)

	service := AddressServiceImpl{
		Repo: repo,
	}

	repo.On(
		"Create",
		mock.Anything,
		mock.Anything,
	).Run(func(args mock.Arguments) {

		address := args.Get(1).(*domain.Address)

		address.ID = 1

	}).Return(nil)

	result, err := service.Create(
		context.Background(),
		1,
		request.CreateAddressRequest{
			RecipientName: "John Doe",
			Phone:         "08123456789",
			Address:       "Jl. Malioboro",
			City:          "Yogyakarta",
			PostalCode:    "55213",
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
		"John Doe",
		result.RecipientName,
	)

	assert.Equal(
		t,
		"08123456789",
		result.Phone,
	)

	assert.Equal(
		t,
		"Jl. Malioboro",
		result.Address,
	)

	assert.Equal(
		t,
		"Yogyakarta",
		result.City,
	)

	assert.Equal(
		t,
		"55213",
		result.PostalCode,
	)

	repo.AssertExpectations(t)
}

func TestAddressCreateRepositoryError(t *testing.T) {

	repo := new(MockAddressRepository)

	service := AddressServiceImpl{
		Repo: repo,
	}

	repo.On(
		"Create",
		mock.Anything,
		mock.Anything,
	).Return(
		errors.New("database error"),
	)

	result, err := service.Create(
		context.Background(),
		1,
		request.CreateAddressRequest{
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
