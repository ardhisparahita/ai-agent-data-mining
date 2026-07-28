package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
	"github.com/ardhisparahita/ecommerce-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCartAddSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockCartService)

	handler := CartHandler{
		Service: service,
	}

	app.Post("/carts", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.AddToCart(c)
	})

	reqBody := request.AddToCartRequest{
		ProductID: 1,
		Quantity:  2,
	}

	service.On(
		"AddToCart",
		mock.Anything,
		uint64(1),
		reqBody,
	).Return(
		&response.CartResponse{
			ID:       1,
			Quantity: 2,
			Subtotal: 300000,
			Product: response.CartProductResponse{
				ID:       1,
				Name:     "Keyboard",
				Price:    150000,
				ImageURL: "keyboard.jpg",
			},
		},
		nil,
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/carts",
		bytes.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	service.AssertExpectations(t)
}

func TestCartAddValidationError(t *testing.T) {

	app := fiber.New()

	service := new(MockCartService)

	handler := CartHandler{
		Service: service,
	}

	app.Post("/carts", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.AddToCart(c)
	})

	reqBody := map[string]any{
		"product_id": 1,
		"quantity":   0,
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/carts",
		bytes.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCartAddServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockCartService)

	handler := CartHandler{
		Service: service,
	}

	app.Post("/carts", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.AddToCart(c)
	})

	reqBody := request.AddToCartRequest{
		ProductID: 1,
		Quantity:  2,
	}

	service.On(
		"AddToCart",
		mock.Anything,
		uint64(1),
		reqBody,
	).Return(
		nil,
		utils.BadRequest("insufficient stock"),
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/carts",
		bytes.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	service.AssertExpectations(t)
}

func TestCartFindAllSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockCartService)

	handler := CartHandler{
		Service: service,
	}

	app.Get("/carts", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.FindAll(c)
	})

	service.On(
		"FindAll",
		mock.Anything,
		uint64(1),
	).Return(
		&response.CartListResponse{
			Items: []response.CartResponse{
				{
					ID:       1,
					Quantity: 2,
					Subtotal: 300000,
					Product: response.CartProductResponse{
						ID:       1,
						Name:     "Keyboard",
						Price:    150000,
						ImageURL: "keyboard.jpg",
					},
				},
				{
					ID:       2,
					Quantity: 1,
					Subtotal: 500000,
					Product: response.CartProductResponse{
						ID:       2,
						Name:     "Mouse",
						Price:    500000,
						ImageURL: "mouse.jpg",
					},
				},
			},
			TotalItems: 2,
			GrandTotal: 800000,
		},
		nil,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/carts",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	service.AssertExpectations(t)
}

func TestCartFindAllServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockCartService)

	handler := CartHandler{
		Service: service,
	}

	app.Get("/carts", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.FindAll(c)
	})

	service.On(
		"FindAll",
		mock.Anything,
		uint64(1),
	).Return(
		nil,
		utils.BadRequest("failed get carts"),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/carts",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	service.AssertExpectations(t)
}

func TestCartUpdateSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockCartService)

	handler := CartHandler{
		Service: service,
	}

	app.Put("/carts/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.Update(c)
	})

	reqBody := request.UpdateCartRequest{
		Quantity: 5,
	}

	service.On(
		"Update",
		mock.Anything,
		uint64(1),
		uint64(1),
		reqBody,
	).Return(
		&response.CartResponse{
			ID:       1,
			Quantity: 5,
			Subtotal: 750000,
			Product: response.CartProductResponse{
				ID:       1,
				Name:     "Keyboard",
				Price:    150000,
				ImageURL: "keyboard.jpg",
			},
		},
		nil,
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/carts/1",
		bytes.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	service.AssertExpectations(t)
}

func TestCartUpdateInvalidID(t *testing.T) {

	app := fiber.New()

	service := new(MockCartService)

	handler := CartHandler{
		Service: service,
	}

	app.Put("/carts/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.Update(c)
	})

	reqBody := request.UpdateCartRequest{
		Quantity: 5,
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/carts/abc",
		bytes.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCartUpdateValidationError(t *testing.T) {

	app := fiber.New()

	service := new(MockCartService)

	handler := CartHandler{
		Service: service,
	}

	app.Put("/carts/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.Update(c)
	})

	reqBody := map[string]any{
		"quantity": 0,
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/carts/1",
		bytes.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCartUpdateServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockCartService)

	handler := CartHandler{
		Service: service,
	}

	app.Put("/carts/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.Update(c)
	})

	reqBody := request.UpdateCartRequest{
		Quantity: 5,
	}

	service.On(
		"Update",
		mock.Anything,
		uint64(1),
		uint64(1),
		reqBody,
	).Return(
		nil,
		utils.NotFound("cart not found"),
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/carts/1",
		bytes.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	service.AssertExpectations(t)
}

func TestCartDeleteSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockCartService)

	handler := CartHandler{
		Service: service,
	}

	app.Delete("/carts/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.Delete(c)
	})

	service.On(
		"Delete",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(nil)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/carts/1",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	service.AssertExpectations(t)
}

func TestCartDeleteInvalidID(t *testing.T) {

	app := fiber.New()

	service := new(MockCartService)

	handler := CartHandler{
		Service: service,
	}

	app.Delete("/carts/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.Delete(c)
	})

	req := httptest.NewRequest(
		http.MethodDelete,
		"/carts/abc",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCartDeleteServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockCartService)

	handler := CartHandler{
		Service: service,
	}

	app.Delete("/carts/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.Delete(c)
	})

	service.On(
		"Delete",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(
		utils.NotFound("cart not found"),
	)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/carts/1",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	service.AssertExpectations(t)
}
