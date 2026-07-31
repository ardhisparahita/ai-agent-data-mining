package service

import (
	"context"
	"testing"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestLoginSuccess(t *testing.T) {

	hashedPassword, _ := bcrypt.GenerateFromPassword(
		[]byte("password123"),
		bcrypt.DefaultCost,
	)

	repo := new(MockUserRepository)

	service := AuthServiceImpl{
		Repo: repo,
	}

	repo.On(
		"FindByEmail",
		mock.Anything,
		"ardhis@gmail.com",
	).Return(
		&domain.User{
			ID:       1,
			Name:     "Ardhis",
			Email:    "ardhis@gmail.com",
			Password: string(hashedPassword),
			Role:     "CUSTOMER",
		},
		nil,
	)

	result, err := service.Login(
		context.Background(),
		request.LoginRequest{
			Email:    "ardhis@gmail.com",
			Password: "password123",
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.NotEmpty(
		t,
		result.Token,
	)

	repo.AssertExpectations(t)
}

func TestLoginUserNotFound(t *testing.T) {

	repo := new(MockUserRepository)

	service := AuthServiceImpl{
		Repo: repo,
	}

	repo.On(
		"FindByEmail",
		mock.Anything,
		"ardhis@gmail.com",
	).Return(
		nil,
		gorm.ErrRecordNotFound,
	)

	result, err := service.Login(
		context.Background(),
		request.LoginRequest{
			Email:    "ardhis@gmail.com",
			Password: "password123",
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"invalid email or password",
	)

	repo.AssertExpectations(t)
}

func TestLoginWrongPassword(t *testing.T) {

	hashedPassword, _ := bcrypt.GenerateFromPassword(
		[]byte("password123"),
		bcrypt.DefaultCost,
	)

	repo := new(MockUserRepository)

	service := AuthServiceImpl{
		Repo: repo,
	}

	repo.On(
		"FindByEmail",
		mock.Anything,
		"ardhis@gmail.com",
	).Return(
		&domain.User{
			ID:       1,
			Name:     "Ardhis",
			Email:    "ardhis@gmail.com",
			Password: string(hashedPassword),
			Role:     "CUSTOMER",
		},
		nil,
	)

	result, err := service.Login(
		context.Background(),
		request.LoginRequest{
			Email:    "ardhis@gmail.com",
			Password: "wrongpassword",
		},
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.Contains(
		t,
		err.Error(),
		"invalid email or password",
	)

	repo.AssertExpectations(t)
}
