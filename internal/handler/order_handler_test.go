package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
	"github.com/ardhisparahita/ecommerce-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderFindAllSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockOrderService)

	handler := OrderHandler{
		Service: service,
	}

	app.Get("/orders", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.FindAll(c)
	})

	service.On(
		"FindAll",
		mock.Anything,
		uint64(1),
	).Return(
		[]response.OrderListResponse{
			{
				ID:          1,
				TotalAmount: 1500000,
				Status:      "PENDING",
				CreatedAt:   time.Now(),
			},
			{
				ID:          2,
				TotalAmount: 2500000,
				Status:      "COMPLETED",
				CreatedAt:   time.Now(),
			},
		},
		nil,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/orders",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	service.AssertExpectations(t)
}

func TestOrderFindAllServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockOrderService)

	handler := OrderHandler{
		Service: service,
	}

	app.Get("/orders", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.FindAll(c)
	})

	service.On(
		"FindAll",
		mock.Anything,
		uint64(1),
	).Return(
		nil,
		utils.BadRequest("failed get orders"),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/orders",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	service.AssertExpectations(t)
}
