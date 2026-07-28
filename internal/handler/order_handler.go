package handler

import (
	"strconv"

	_ "github.com/ardhisparahita/ecommerce-api/internal/dto/response"
	"github.com/ardhisparahita/ecommerce-api/internal/service"
	"github.com/ardhisparahita/ecommerce-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
)
type OrderHandler struct {
	Service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{Service: service}
}

// FindAll godoc
//
// @Summary Get user orders
// @Description Get all orders owned by authenticated user
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.OrderListSwaggerResponse
// @Failure 401 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /orders [get]
func (h *OrderHandler) FindAll(c *fiber.Ctx) error {
	userID := utils.GetUserID(c)

	data, err := h.Service.FindAll(c.UserContext(), userID)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"get orders",
		data,
	)
}

// FindByID godoc
//
// @Summary Get order detail
// @Description Get order detail by id
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.OrderSwaggerResponse
// @Failure 401 {object} response.ErrorSwaggerResponse
// @Failure 404 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /orders/{id} [get]
func (h *OrderHandler) FindByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid order Id")
	}

	userID := utils.GetUserID(c)

	data, err := h.Service.FindByID(c.UserContext(), id, userID)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"get order details",
		data,
	)
}

// Cancel godoc
//
// @Summary Cancel order
// @Description Cancel pending order
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.MessageSwaggerResponse
// @Failure 400 {object} response.ErrorSwaggerResponse
// @Failure 404 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /orders/{id}/cancel [patch]
func (h *OrderHandler) Cancel(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid order Id")
	}

	userID := utils.GetUserID(c)

	err = h.Service.Cancel(c.UserContext(), id, userID)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"order cancelled",
		nil,
	)
}

// MarkAsShipped godoc
//
// @Summary Mark order as shipped
// @Description Change order status to SHIPPED
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.MessageSwaggerResponse
// @Failure 400 {object} response.ErrorSwaggerResponse
// @Failure 404 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /orders/{id}/ship [patch]
func (h *OrderHandler) MarkAsShipped(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid order Id")
	}

	err = h.Service.MarkAsShipped(c.UserContext(), id)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"order shipped",
		nil,
	)

}

// MarkAsCompleted godoc
//
// @Summary Complete order
// @Description Change order status to COMPLETED
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.MessageSwaggerResponse
// @Failure 400 {object} response.ErrorSwaggerResponse
// @Failure 401 {object} response.ErrorSwaggerResponse
// @Failure 403 {object} response.ErrorSwaggerResponse
// @Failure 404 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /orders/{id}/complete [patch]
func (h *OrderHandler) MarkAsCompleted(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid order id")
	}

	userID := c.Locals("user_id").(uint64)
	role := c.Locals("role").(string)

	err = h.Service.MarkAsCompleted(
		c.UserContext(),
		id,
		userID,
		role,
	)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"order completed",
		nil,
	)
}

// MarkAsPaid godoc
//
// @Summary Mark order as paid
// @Description Change order status to PAID
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.MessageSwaggerResponse
// @Failure 400 {object} response.ErrorSwaggerResponse
// @Failure 404 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /orders/{id}/pay [patch]
func (h *OrderHandler) MarkAsPaid(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid order Id")
	}

	err = h.Service.MarkAsPaid(c.UserContext(), id)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"payment success",
		nil,
	)
}

// MarkAsFailed godoc
//
// @Summary Mark order as failed
// @Description Change order status to FAILED
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.MessageSwaggerResponse
// @Failure 400 {object} response.ErrorSwaggerResponse
// @Failure 404 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /orders/{id}/fail [patch]
func (h *OrderHandler) MarkAsFailed(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid order Id")
	}

	err = h.Service.MarkAsFailed(c.UserContext(), id)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"payment failed",
		nil,
	)
}
