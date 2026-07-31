package utils

import (
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
	"github.com/gofiber/fiber/v2"
)

func ResponseSuccess(c *fiber.Ctx, code int, message string, data interface{}) error {
	return c.Status(code).JSON(
		response.WebResponse{
			Code:    code,
			Status:  "success",
			Message: message,
			Data:    data,
		},
	)
}

func ResponseError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"code":    fiber.StatusBadRequest,
		"status":  "error",
		"message": "validation failed",
		"errors":  ValidationErrors(err),
	})
}
