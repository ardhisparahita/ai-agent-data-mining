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

func TestGetProfileSuccess(t *testing.T) {

	repo := new(MockUserRepository)

	service := AuthServiceImpl{
		Repo: repo,
	}

	repo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(
		&domain.User{
			ID:    1,
			Name:  "Ardhis",
			Email: "ardhis@gmail.com",
			Role:  "CUSTOMER",
		},
		nil,
	)

	result, err := service.GetProfile(
		context.Background(),
		1,
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
		"Ardhis",
		result.Name,
	)

	assert.Equal(
		t,
		"ardhis@gmail.com",
		result.Email,
	)

	assert.Equal(
		t,
		"CUSTOMER",
		result.Role,
	)

	repo.AssertExpectations(t)
}

func TestGetProfileNotFound(t *testing.T) {

	repo := new(MockUserRepository)

	service := AuthServiceImpl{
		Repo: repo,
	}

	repo.On(
		"FindByID",
		mock.Anything,
		uint64(99),
	).Return(
		nil,
		gorm.ErrRecordNotFound,
	)

	result, err := service.GetProfile(
		context.Background(),
		99,
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"user not found",
	)

	repo.AssertExpectations(t)
}

func TestGetProfileRepositoryError(t *testing.T) {

	repo := new(MockUserRepository)

	service := AuthServiceImpl{
		Repo: repo,
	}

	repo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(
		nil,
		errors.New("database error"),
	)

	result, err := service.GetProfile(
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

func TestUpdateProfileSuccess(t *testing.T) {

	repo := new(MockUserRepository)

	service := AuthServiceImpl{
		Repo: repo,
	}

	user := &domain.User{
		ID:    1,
		Name:  "Old Name",
		Email: "ardhis@gmail.com",
		Role:  "CUSTOMER",
	}

	repo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(
		user,
		nil,
	)

	repo.On(
		"Update",
		mock.Anything,
		mock.Anything,
	).Return(nil)

	result, err := service.UpdateProfile(
		context.Background(),
		1,
		request.UpdateProfileRequest{
			Name: "New Name",
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
		"New Name",
		result.Name,
	)

	assert.Equal(
		t,
		"ardhis@gmail.com",
		result.Email,
	)

	assert.Equal(
		t,
		"CUSTOMER",
		result.Role,
	)

	repo.AssertExpectations(t)
}

func TestUpdateProfileNotFound(t *testing.T) {

	repo := new(MockUserRepository)

	service := AuthServiceImpl{
		Repo: repo,
	}

	repo.On(
		"FindByID",
		mock.Anything,
		uint64(99),
	).Return(
		nil,
		gorm.ErrRecordNotFound,
	)

	result, err := service.UpdateProfile(
		context.Background(),
		99,
		request.UpdateProfileRequest{
			Name: "New Name",
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"user not found",
	)

	repo.AssertExpectations(t)
}

func TestUpdateProfileRepositoryError(t *testing.T) {

	repo := new(MockUserRepository)

	service := AuthServiceImpl{
		Repo: repo,
	}

	repo.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(
		nil,
		errors.New("database error"),
	)

	result, err := service.UpdateProfile(
		context.Background(),
		1,
		request.UpdateProfileRequest{
			Name: "New Name",
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
