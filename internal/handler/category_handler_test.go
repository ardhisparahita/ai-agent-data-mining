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

func setupCategoryHandler() (*fiber.App, *MockCategoryService) {

	app := fiber.New()

	service := new(MockCategoryService)

	handler := CategoryHandler{
		Service: service,
	}

	app.Post("/categories", handler.Create)

	return app, service
}

func TestCategoryCreateSuccess(t *testing.T) {

	app, service := setupCategoryHandler()

	reqBody := request.CreateCategoryRequest{
		Name: "Keyboard",
	}

	service.On(
		"Create",
		mock.Anything,
		reqBody,
	).Return(
		&response.CategoryResponse{
			ID:   1,
			Name: "Keyboard",
		},
		nil,
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/categories",
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
		fiber.StatusCreated,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}

func TestCategoryCreateValidationError(t *testing.T) {

	app, _ := setupCategoryHandler()

	reqBody := map[string]any{
		"name": "",
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/categories",
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

func TestCategoryCreateServiceError(t *testing.T) {

	app, service := setupCategoryHandler()

	reqBody := request.CreateCategoryRequest{
		Name: "Keyboard",
	}

	service.On(
		"Create",
		mock.Anything,
		reqBody,
	).Return(
		nil,
		utils.Conflict("category already exists"),
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/categories",
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

func TestCategoryFindAllSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockCategoryService)

	handler := CategoryHandler{
		Service: service,
	}

	app.Get("/categories", handler.FindAll)

	service.On(
		"FindAll",
		mock.Anything,
	).Return(
		[]response.CategoryResponse{
			{
				ID:   1,
				Name: "Keyboard",
			},
			{
				ID:   2,
				Name: "Mouse",
			},
		},
		nil,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/categories",
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

func TestCategoryFindAllServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockCategoryService)

	handler := CategoryHandler{
		Service: service,
	}

	app.Get("/categories", handler.FindAll)

	service.On(
		"FindAll",
		mock.Anything,
	).Return(
		nil,
		utils.BadRequest("failed to get categories"),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/categories",
		nil,
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

func TestCategoryFindByIDSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockCategoryService)

	handler := CategoryHandler{
		Service: service,
	}

	app.Get("/categories/:id", handler.FindByID)

	service.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(
		&response.CategoryResponse{
			ID:   1,
			Name: "Keyboard",
		},
		nil,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/categories/1",
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

func TestCategoryFindByIDInvalidID(t *testing.T) {

	app := fiber.New()

	service := new(MockCategoryService)

	handler := CategoryHandler{
		Service: service,
	}

	app.Get("/categories/:id", handler.FindByID)

	req := httptest.NewRequest(
		http.MethodGet,
		"/categories/abc",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusBadRequest,
		resp.StatusCode,
	)
}

func TestCategoryFindByIDServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockCategoryService)

	handler := CategoryHandler{
		Service: service,
	}

	app.Get("/categories/:id", handler.FindByID)

	service.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(
		nil,
		utils.NotFound("category not found"),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/categories/1",
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

func TestCategoryUpdateSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockCategoryService)

	handler := CategoryHandler{
		Service: service,
	}

	app.Put("/categories/:id", handler.Update)

	reqBody := request.UpdateCategoryRequest{
		Name: "Keyboard Gaming",
	}

	service.On(
		"Update",
		mock.Anything,
		uint64(1),
		reqBody,
	).Return(
		&response.CategoryResponse{
			ID:   1,
			Name: "Keyboard Gaming",
		},
		nil,
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/categories/1",
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

func TestCategoryUpdateInvalidID(t *testing.T) {

	app := fiber.New()

	service := new(MockCategoryService)

	handler := CategoryHandler{
		Service: service,
	}

	app.Put("/categories/:id", handler.Update)

	reqBody := request.UpdateCategoryRequest{
		Name: "Keyboard Gaming",
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/categories/abc",
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

func TestCategoryUpdateValidationError(t *testing.T) {

	app := fiber.New()

	service := new(MockCategoryService)

	handler := CategoryHandler{
		Service: service,
	}

	app.Put("/categories/:id", handler.Update)

	reqBody := map[string]any{
		"name": "",
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/categories/1",
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

func TestCategoryUpdateServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockCategoryService)

	handler := CategoryHandler{
		Service: service,
	}

	app.Put("/categories/:id", handler.Update)

	reqBody := request.UpdateCategoryRequest{
		Name: "Keyboard Gaming",
	}

	service.On(
		"Update",
		mock.Anything,
		uint64(1),
		reqBody,
	).Return(
		nil,
		utils.NotFound("category not found"),
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/categories/1",
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
