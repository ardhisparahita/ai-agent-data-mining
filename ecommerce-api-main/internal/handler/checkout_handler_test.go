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

func TestCheckoutSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockCheckoutService)

	handler := CheckoutHandler{
		Service: service,
	}

	app.Post("/checkouts", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.Checkout(c)
	})

	reqBody := request.CheckoutRequest{
		AddressID:     1,
		PaymentMethod: "BANK_TRANSFER",
	}

	service.On(
		"Checkout",
		mock.Anything,
		uint64(1),
		reqBody,
	).Return(
		&response.OrderResponse{
			ID:            1,
			TotalAmount:   1500000,
			Status:        "PENDING",
			PaymentMethod: "BANK_TRANSFER",
		},
		nil,
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/checkouts",
		bytes.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	service.AssertExpectations(t)
}

func TestCheckoutValidationError(t *testing.T) {

	app := fiber.New()

	service := new(MockCheckoutService)

	handler := CheckoutHandler{
		Service: service,
	}

	app.Post("/checkouts", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.Checkout(c)
	})

	reqBody := map[string]any{
		"address_id": 0,
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/checkouts",
		bytes.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCheckoutServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockCheckoutService)

	handler := CheckoutHandler{
		Service: service,
	}

	app.Post("/checkouts", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.Checkout(c)
	})

	reqBody := request.CheckoutRequest{
		AddressID:     1,
		PaymentMethod: "BANK_TRANSFER",
	}

	service.On(
		"Checkout",
		mock.Anything,
		uint64(1),
		reqBody,
	).Return(
		nil,
		utils.BadRequest("cart is empty"),
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/checkouts",
		bytes.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	service.AssertExpectations(t)
}
