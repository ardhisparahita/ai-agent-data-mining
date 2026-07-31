package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
	"github.com/ardhisparahita/ecommerce-api/pkg/utils"
)

func setupAuthHandler() (*fiber.App, *MockAuthService) {

	app := fiber.New()

	service := new(MockAuthService)

	handler := AuthHandler{
		Service: service,
	}

	app.Post("/register", handler.Register)

	return app, service
}

func TestRegisterSuccess(t *testing.T) {

	app, service := setupAuthHandler()

	reqBody := request.RegisterRequest{
		Name:     "Ardhis",
		Email:    "ardhis@gmail.com",
		Password: "password123",
	}

	service.On(
		"Register",
		mock.Anything,
		reqBody,
	).Return(
		&response.UserResponse{
			ID:    1,
			Name:  "Ardhis",
			Email: "ardhis@gmail.com",
			Role:  "CUSTOMER",
		},
		nil,
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/register",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	service.AssertExpectations(t)
}

func TestRegisterValidationError(t *testing.T) {

	app, _ := setupAuthHandler()

	reqBody := map[string]any{
		"name": "",
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/register",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusBadRequest,
		resp.StatusCode,
	)
}

func TestRegisterServiceError(t *testing.T) {

	app, service := setupAuthHandler()

	reqBody := request.RegisterRequest{
		Name:     "Ardhis",
		Email:    "ardhis@gmail.com",
		Password: "password123",
	}

	service.On(
		"Register",
		mock.Anything,
		reqBody,
	).Return(
		nil,
		utils.Conflict("email already registered"),
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/register",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusConflict,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}

func TestLoginSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockAuthService)

	handler := AuthHandler{
		Service: service,
	}

	app.Post("/login", handler.Login)

	reqBody := request.LoginRequest{
		Email:    "ardhis@gmail.com",
		Password: "password123",
	}

	service.On(
		"Login",
		mock.Anything,
		reqBody,
	).Return(
		&response.AuthResponse{
			Token: "dummy-jwt-token",
		},
		nil,
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusOK,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}

func TestLoginValidationError(t *testing.T) {

	app := fiber.New()

	service := new(MockAuthService)

	handler := AuthHandler{
		Service: service,
	}

	app.Post("/login", handler.Login)

	reqBody := map[string]any{
		"email": "",
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusBadRequest,
		resp.StatusCode,
	)
}

func TestLoginUnauthorized(t *testing.T) {

	app := fiber.New()

	service := new(MockAuthService)

	handler := AuthHandler{
		Service: service,
	}

	app.Post("/login", handler.Login)

	reqBody := request.LoginRequest{
		Email:    "ardhis@gmail.com",
		Password: "wrongpassword",
	}

	service.On(
		"Login",
		mock.Anything,
		reqBody,
	).Return(
		nil,
		utils.Unauthorized("invalid email or password"),
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusUnauthorized,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}

func TestGetProfileSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockAuthService)

	handler := AuthHandler{
		Service: service,
	}

	app.Get("/profile", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.GetProfile(c)
	})

	service.On(
		"GetProfile",
		mock.Anything,
		uint64(1),
	).Return(
		&response.UserResponse{
			ID:    1,
			Name:  "Ardhis",
			Email: "ardhis@gmail.com",
			Role:  "CUSTOMER",
		},
		nil,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/profile",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusOK,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}

func TestGetProfileServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockAuthService)

	handler := AuthHandler{
		Service: service,
	}

	app.Get("/profile", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.GetProfile(c)
	})

	service.On(
		"GetProfile",
		mock.Anything,
		uint64(1),
	).Return(
		nil,
		utils.NotFound("user not found"),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/profile",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusNotFound,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}

func TestUpdateProfileSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockAuthService)

	handler := AuthHandler{
		Service: service,
	}

	app.Put("/profile", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.UpdateProfile(c)
	})

	reqBody := request.UpdateProfileRequest{
		Name: "Ardhis Updated",
	}

	service.On(
		"UpdateProfile",
		mock.Anything,
		uint64(1),
		reqBody,
	).Return(
		&response.UserResponse{
			ID:    1,
			Name:  "Ardhis Updated",
			Email: "ardhis@gmail.com",
			Role:  "CUSTOMER",
		},
		nil,
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/profile",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusOK,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}

func TestUpdateProfileValidationError(t *testing.T) {

	app := fiber.New()

	service := new(MockAuthService)

	handler := AuthHandler{
		Service: service,
	}

	app.Put("/profile", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.UpdateProfile(c)
	})

	reqBody := map[string]any{
		"name": "",
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/profile",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusBadRequest,
		resp.StatusCode,
	)
}

func TestUpdateProfileServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockAuthService)

	handler := AuthHandler{
		Service: service,
	}

	app.Put("/profile", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.UpdateProfile(c)
	})

	reqBody := request.UpdateProfileRequest{
		Name: "Ardhis Updated",
	}

	service.On(
		"UpdateProfile",
		mock.Anything,
		uint64(1),
		reqBody,
	).Return(
		nil,
		utils.NotFound("user not found"),
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/profile",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusNotFound,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}

func TestChangePasswordSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockAuthService)

	handler := AuthHandler{
		Service: service,
	}

	app.Put("/change-password", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.ChangePassword(c)
	})

	reqBody := request.ChangePasswordRequest{
		OldPassword: "oldpassword",
		NewPassword: "newpassword123",
	}

	service.On(
		"ChangePassword",
		mock.Anything,
		uint64(1),
		reqBody,
	).Return(nil)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/change-password",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusOK,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}

func TestChangePasswordValidationError(t *testing.T) {

	app := fiber.New()

	service := new(MockAuthService)

	handler := AuthHandler{
		Service: service,
	}

	app.Put("/change-password", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.ChangePassword(c)
	})

	reqBody := map[string]any{
		"old_password": "",
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/change-password",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusBadRequest,
		resp.StatusCode,
	)
}

func TestChangePasswordServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockAuthService)

	handler := AuthHandler{
		Service: service,
	}

	app.Put("/change-password", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.ChangePassword(c)
	})

	reqBody := request.ChangePasswordRequest{
		OldPassword: "wrongpassword",
		NewPassword: "newpassword123",
	}

	service.On(
		"ChangePassword",
		mock.Anything,
		uint64(1),
		reqBody,
	).Return(
		utils.BadRequest("old password is incorrect"),
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/change-password",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusBadRequest,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}
