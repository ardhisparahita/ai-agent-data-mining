package handler

import (
	"strconv"

	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	_ "github.com/ardhisparahita/ecommerce-api/internal/dto/response"
	"github.com/ardhisparahita/ecommerce-api/internal/service"
	"github.com/ardhisparahita/ecommerce-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type AddressHandler struct {
	Service service.AddressService
}

func NewAddressHandler(service service.AddressService) *AddressHandler {
	return &AddressHandler{Service: service}
}

// Create godoc
//
// @Summary Create address
// @Description Create a new address for authenticated user
// @Tags Addresses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.CreateAddressRequest true "Address Request"
// @Success 201 {object} response.AddressSwaggerResponse
// @Failure 400 {object} response.ErrorSwaggerResponse
// @Failure 401 {object} response.ErrorSwaggerResponse
// @Failure 422 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /addresses [post]
func (h *AddressHandler) Create(c *fiber.Ctx) error {
	var req request.CreateAddressRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := utils.ValidationStruct(req); err != nil {
		return utils.ResponseError(c, err)
	}

	userID := utils.GetUserID(c)

	data, err := h.Service.Create(c.UserContext(), userID, req)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusCreated,
		"address created",
		data,
	)
}

// FindAll godoc
//
// @Summary Get all addresses
// @Description Get all addresses owned by user
// @Tags Addresses
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.AddressListSwaggerResponse
// @Failure 401 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /addresses [get]
func (h *AddressHandler) FindAll(c *fiber.Ctx) error {
	userID := utils.GetUserID(c)

	addresses, err := h.Service.FindAllByUserID(c.UserContext(), userID)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"get all addresses",
		addresses,
	)
}

// FindByID godoc
//
// @Summary Get address detail
// @Description Get address by id
// @Tags Addresses
// @Security BearerAuth
// @Produce json
// @Param id path int true "Address ID"
// @Success 200 {object} response.AddressSwaggerResponse
// @Failure 401 {object} response.ErrorSwaggerResponse
// @Failure 404 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /addresses/{id} [get]
func (h *AddressHandler) FindByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid address Id")
	}

	userID := utils.GetUserID(c)

	address, err := h.Service.FindByID(c.UserContext(), id, userID)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"get one address",
		address,
	)
}

// Update godoc
//
// @Summary Update address
// @Description Update address by id
// @Tags Addresses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Param request body request.UpdateAddressRequest true "Address Request"
// @Success 200 {object} response.AddressSwaggerResponse
// @Failure 400 {object} response.ErrorSwaggerResponse
// @Failure 401 {object} response.ErrorSwaggerResponse
// @Failure 404 {object} response.ErrorSwaggerResponse
// @Failure 422 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /addresses/{id} [put]
func (h *AddressHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid address Id")
	}

	var req request.UpdateAddressRequest

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
		"address updated",
		data,
	)
}

// Delete godoc
//
// @Summary Delete address
// @Description Delete address by id
// @Tags Addresses
// @Security BearerAuth
// @Produce json
// @Param id path int true "Address ID"
// @Success 200 {object} response.MessageSwaggerResponse
// @Failure 401 {object} response.ErrorSwaggerResponse
// @Failure 404 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /addresses/{id} [delete]
func (h *AddressHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid address Id")
	}

	userID := utils.GetUserID(c)

	err = h.Service.Delete(c.UserContext(), id, userID)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"address deleted",
		nil,
	)
}
