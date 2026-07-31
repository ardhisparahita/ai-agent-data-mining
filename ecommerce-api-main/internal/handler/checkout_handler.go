package handler

import (
	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	_ "github.com/ardhisparahita/ecommerce-api/internal/dto/response"
	"github.com/ardhisparahita/ecommerce-api/internal/service"
	"github.com/ardhisparahita/ecommerce-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type CheckoutHandler struct {
	Service service.CheckoutService
}

func NewCheckoutHandler(service service.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{
		Service: service,
	}
}

// Checkout godoc
//
// @Summary Checkout cart
// @Description Create order from current user's cart
// @Tags Checkout
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.CheckoutRequest true "Checkout Request"
// @Success 201 {object} response.CheckoutSwaggerResponse
// @Failure 400 {object} response.ErrorSwaggerResponse
// @Failure 401 {object} response.ErrorSwaggerResponse
// @Failure 404 {object} response.ErrorSwaggerResponse
// @Failure 422 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /checkouts [post]
func (h *CheckoutHandler) Checkout(c *fiber.Ctx) error {
	var req request.CheckoutRequest

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := utils.ValidationStruct(req); err != nil {
		return utils.ResponseError(c, err)
	}

	userID := utils.GetUserID(c)

	data, err := h.Service.Checkout(c.UserContext(), userID, req)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"checkout success",
		data,
	)
}
