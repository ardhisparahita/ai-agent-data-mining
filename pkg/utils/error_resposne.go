package utils

import "github.com/gofiber/fiber/v2"

func BadRequest(message string) error {
	return fiber.NewError(
		fiber.StatusBadRequest,
		message,
	)
}

func Unauthorized(message string) error {
	return fiber.NewError(
		fiber.StatusUnauthorized,
		message,
	)
}

func Forbidden(message string) error {
	return fiber.NewError(
		fiber.StatusForbidden,
		message,
	)
}

func NotFound(message string) error {
	return fiber.NewError(
		fiber.StatusNotFound,
		message,
	)
}

func Conflict(message string) error {
	return fiber.NewError(
		fiber.StatusConflict,
		message,
	)
}
