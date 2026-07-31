package handler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	_ "github.com/ardhisparahita/ecommerce-api/internal/dto/response"
	"github.com/ardhisparahita/ecommerce-api/internal/service"
	"github.com/ardhisparahita/ecommerce-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	Service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{
		Service: service,
	}
}

// Create godoc
//
// @Summary Create product
// @Description Create a new product
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.CreateProductRequest true "Product Request"
// @Success 201 {object} response.ProductSwaggerResponse
// @Failure 400 {object} response.ErrorSwaggerResponse
// @Failure 401 {object} response.ErrorSwaggerResponse
// @Failure 422 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /products [post]
func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var req request.CreateProductRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := utils.ValidationStruct(req); err != nil {
		return utils.ResponseError(c, err)
	}

	data, err := h.Service.Create(c.UserContext(), req)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusCreated,
		"product created",
		data,
	)
}

// FindAll godoc
//
// @Summary Get all products
// @Description Get products with pagination, search and category filter
// @Tags Products
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page Number"
// @Param limit query int false "Items Per Page"
// @Param search query string false "Search Product"
// @Param category_id query int false "Category ID"
// @Success 200 {object} response.ProductListSwaggerResponse
// @Failure 401 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /products [get]
func (h *ProductHandler) FindAll(c *fiber.Ctx) error {
	var req request.ProductQueryRequest

	if err := c.QueryParser(&req); err != nil {
		return err
	}

	if err := utils.ValidationStruct(req); err != nil {
		return utils.ResponseError(c, err)
	}

	data, err := h.Service.FindAll(c.UserContext(), req)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"get products",
		data,
	)
}

// FindByID godoc
//
// @Summary Get product detail
// @Description Get product by id
// @Tags Products
// @Security BearerAuth
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.ProductSwaggerResponse
// @Failure 401 {object} response.ErrorSwaggerResponse
// @Failure 404 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /products/{id} [get]
func (h *ProductHandler) FindByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid product Id")
	}

	res, err := h.Service.FindByID(c.UserContext(), id)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"get one product",
		res,
	)
}

// Update godoc
//
// @Summary Update product
// @Description Update product by id
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body request.UpdateProductRequest true "Update Product Request"
// @Success 200 {object} response.ProductSwaggerResponse
// @Failure 400 {object} response.ErrorSwaggerResponse
// @Failure 401 {object} response.ErrorSwaggerResponse
// @Failure 404 {object} response.ErrorSwaggerResponse
// @Failure 422 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /products/{id} [put]
func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid product Id")
	}

	var req request.UpdateProductRequest

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := utils.ValidationStruct(req); err != nil {
		return utils.ResponseError(c, err)
	}

	data, err := h.Service.Update(c.UserContext(), id, req)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"product updated",
		data,
	)
}

// Delete godoc
//
// @Summary Delete product
// @Description Delete product by id
// @Tags Products
// @Security BearerAuth
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.MessageSwaggerResponse
// @Failure 401 {object} response.ErrorSwaggerResponse
// @Failure 404 {object} response.ErrorSwaggerResponse
// @Failure 500 {object} response.ErrorSwaggerResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid product Id")
	}

	err = h.Service.Delete(c.UserContext(), id)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"product deleted",
		nil,
	)
}

func (h *ProductHandler) UploadImage(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid category id")
	}

	file, err := c.FormFile("image")
	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"image is required",
		)
	}

	fileName := fmt.Sprintf("%s_%s",
		time.Now().Format("20060102150405"),
		file.Filename,
	)
	path := "./uploads/products/" + fileName

	if err := c.SaveFile(file, path); err != nil {
		return err
	}

	data, err := h.Service.UploadImage(c.UserContext(), id, "/uploads/products/"+fileName)

	if err != nil {
		return err
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"image uploaded successfully",
		data,
	)
}
