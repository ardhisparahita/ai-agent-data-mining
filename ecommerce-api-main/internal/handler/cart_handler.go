package handler

import (
	"strconv"

	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	_ "github.com/ardhisparahita/ecommerce-api/internal/dto/response"
	"github.com/ardhisparahita/ecommerce-api/internal/service"
	"github.com/ardhisparahita/ecommerce-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type CartHandler struct {
	Service service.CartService
}

func NewCartHandler(service service.CartService) *CartHandler {
	return &CartHandler{Service: service}
}

// AddToCart godoc
//
// @Summary Add product to cart
// @Description Add a product to shopping cart
// @Tags Carts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.AddToCartRequest true "Add To Cart Request"
// @Success 201 {object} response.CartSwaggerResponse
// @Failure 400 {object} response.ErrorSwaggerResponse
// @Failure 401 {object} response.ErrorSwaggerResponse
// @Failure 404 {object} response.ErrorSwaggerResponse
// @Failure 422 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /carts [post]
func (h *CartHandler) AddToCart(c *fiber.Ctx) error {
	var req request.AddToCartRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := utils.ValidationStruct(req); err != nil {
		return utils.ResponseError(c, err)
	}

	userID := utils.GetUserID(c)

	data, err := h.Service.AddToCart(c.UserContext(), userID, req)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusCreated,
		"product added to cart",
		data,
	)
}

// FindAll godoc
//
// @Summary Get cart items
// @Description Get all cart items of current user
// @Tags Carts
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.CartListSwaggerResponse
// @Failure 401 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /carts [get]
func (h *CartHandler) FindAll(c *fiber.Ctx) error {
	userID := utils.GetUserID(c)

	data, err := h.Service.FindAll(c.UserContext(), userID)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"get all carts",
		data,
	)
}

// Update godoc
//
// @Summary Update cart quantity
// @Description Update quantity of cart item
// @Tags Carts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Cart ID"
// @Param request body request.UpdateCartRequest true "Update Cart Request"
// @Success 200 {object} response.CartSwaggerResponse
// @Failure 400 {object} response.ErrorSwaggerResponse
// @Failure 401 {object} response.ErrorSwaggerResponse
// @Failure 404 {object} response.ErrorSwaggerResponse
// @Failure 422 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /carts/{id} [put]
func (h *CartHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid cart id")
	}

	var req request.UpdateCartRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := utils.ValidationStruct(req); err != nil {
		return utils.ResponseError(c, err)
	}

	userID := utils.GetUserID(c)

	data, err := h.Service.Update(c.UserContext(), id, userID, req)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"cart updated",
		data,
	)
}

// Delete godoc
//
// @Summary Delete cart item
// @Description Remove item from cart
// @Tags Carts
// @Security BearerAuth
// @Produce json
// @Param id path int true "Cart ID"
// @Success 200 {object} response.MessageSwaggerResponse
// @Failure 401 {object} response.ErrorSwaggerResponse
// @Failure 404 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /carts/{id} [delete]
func (h *CartHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid cart id")
	}

	userID := utils.GetUserID(c)
	if err := h.Service.Delete(c.UserContext(), id, userID); err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"cart deleted",
		nil,
	)
}
