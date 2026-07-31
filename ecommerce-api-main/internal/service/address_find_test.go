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

func TestAddressFindAllSuccess(t *testing.T) {

	repo := new(MockAddressRepository)

	service := AddressServiceImpl{
		Repo: repo,
	}

	repo.On(
		"FindAllByUserID",
		mock.Anything,
		uint64(1),
	).Return(
		[]domain.Address{
			{
				ID:            1,
				UserID:        1,
				RecipientName: "John Doe",
				Phone:         "08123456789",
				Address:       "Jl. Malioboro",
				City:          "Yogyakarta",
				PostalCode:    "55213",
			},
			{
				ID:            2,
				UserID:        1,
				RecipientName: "Jane Doe",
				Phone:         "08987654321",
				Address:       "Jl. Solo",
				City:          "Solo",
				PostalCode:    "57111",
			},
		},
		nil,
	)

	result, err := service.FindAllByUserID(
		context.Background(),
		1,
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Len(t, result, 2)

	assert.Equal(t, uint64(1), result[0].ID)
	assert.Equal(t, "John Doe", result[0].RecipientName)

	assert.Equal(t, uint64(2), result[1].ID)
	assert.Equal(t, "Jane Doe", result[1].RecipientName)

	repo.AssertExpectations(t)
}

func TestAddressFindAllRepositoryError(t *testing.T) {

	repo := new(MockAddressRepository)

	service := AddressServiceImpl{
		Repo: repo,
	}

	repo.On(
		"FindAllByUserID",
		mock.Anything,
		uint64(1),
	).Return(
		nil,
		errors.New("database error"),
	)

	result, err := service.FindAllByUserID(
		context.Background(),
		1,
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

func TestAddressFindByIDSuccess(t *testing.T) {

	repo := new(MockAddressRepository)

	service := AddressServiceImpl{
		Repo: repo,
	}

	repo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(
		&domain.Address{
			ID:            1,
			UserID:        1,
			RecipientName: "John Doe",
			Phone:         "08123456789",
			Address:       "Jl. Malioboro",
			City:          "Yogyakarta",
			PostalCode:    "55213",
		},
		nil,
	)

	result, err := service.FindByID(
		context.Background(),
		1,
		1,
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

func TestAddressFindByIDNotFound(t *testing.T) {

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

	result, err := service.FindByID(
		context.Background(),
		99,
		1,
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

func TestAddressFindByIDRepositoryError(t *testing.T) {

	repo := new(MockAddressRepository)

	service := AddressServiceImpl{
		Repo: repo,
	}

	repo.On(
		"FindByIDAndUserID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(
		nil,
		errors.New("database error"),
	)

	result, err := service.FindByID(
		context.Background(),
		1,
		1,
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
