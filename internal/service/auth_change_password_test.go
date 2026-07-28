package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestChangePasswordSuccess(t *testing.T) {

	hashedPassword, _ := bcrypt.GenerateFromPassword(
		[]byte("oldpassword"),
		bcrypt.DefaultCost,
	)

	repo := new(MockUserRepository)

	service := AuthServiceImpl{
		Repo: repo,
	}

	user := &domain.User{
		ID:       1,
		Name:     "Ardhis",
		Email:    "ardhis@gmail.com",
		Password: string(hashedPassword),
		Role:     "CUSTOMER",
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

	err := service.ChangePassword(
		context.Background(),
		1,
		request.ChangePasswordRequest{
			OldPassword: "oldpassword",
			NewPassword: "newpassword123",
		},
	)

	assert.NoError(t, err)

	assert.NotEqual(
		t,
		string(hashedPassword),
		user.Password,
	)

	assert.NoError(
		t,
		bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte("newpassword123"),
		),
	)

	repo.AssertExpectations(t)
}

func TestChangePasswordUserNotFound(t *testing.T) {

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

	err := service.ChangePassword(
		context.Background(),
		99,
		request.ChangePasswordRequest{
			OldPassword: "oldpassword",
			NewPassword: "newpassword",
		},
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"user not found",
	)

	repo.AssertExpectations(t)
}

func TestChangePasswordWrongOldPassword(t *testing.T) {

	hashedPassword, _ := bcrypt.GenerateFromPassword(
		[]byte("oldpassword"),
		bcrypt.DefaultCost,
	)

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
			ID:       1,
			Password: string(hashedPassword),
		},
		nil,
	)

	err := service.ChangePassword(
		context.Background(),
		1,
		request.ChangePasswordRequest{
			OldPassword: "wrongpassword",
			NewPassword: "newpassword",
		},
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"old password is incorrect",
	)

	repo.AssertExpectations(t)
}

func TestChangePasswordRepositoryError(t *testing.T) {

	hashedPassword, _ := bcrypt.GenerateFromPassword(
		[]byte("oldpassword"),
		bcrypt.DefaultCost,
	)

	repo := new(MockUserRepository)

	service := AuthServiceImpl{
		Repo: repo,
	}

	user := &domain.User{
		ID:       1,
		Password: string(hashedPassword),
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
	).Return(
		errors.New("database error"),
	)

	err := service.ChangePassword(
		context.Background(),
		1,
		request.ChangePasswordRequest{
			OldPassword: "oldpassword",
			NewPassword: "newpassword",
		},
	)

	assert.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"database error",
	)

	repo.AssertExpectations(t)
}
