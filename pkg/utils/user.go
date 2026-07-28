package utils

import "github.com/gofiber/fiber/v2"

func GetUserID(c *fiber.Ctx) uint64 {
	return c.Locals("user_id").(uint64)
}
