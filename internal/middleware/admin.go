package middleware

import (
	"github.com/ardhisparahita/ecommerce-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if utils.GetRole(c) != "ADMIN" {
			return utils.Forbidden("admin access required")
		}
		return c.Next()
	}

}
