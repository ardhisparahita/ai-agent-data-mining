package utils

import "github.com/gofiber/fiber/v2"

func GetRole(c *fiber.Ctx) string {
	return c.Locals("role").(string)
}
