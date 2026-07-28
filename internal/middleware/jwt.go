package middleware

import (
	"strings"

	"github.com/ardhisparahita/ecommerce-api/pkg/config"
	"github.com/ardhisparahita/ecommerce-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return utils.Unauthorized("missing token")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := utils.ValidateToken(tokenString, config.Get("JWT_SECRET"))
		if err != nil || !token.Valid {
			return utils.Unauthorized("invalid token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return utils.Unauthorized("invalid token claims")
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return utils.Unauthorized("invalid user id")
		}

		role, ok := claims["role"].(string)
		if !ok {
			return utils.Unauthorized("invalid role")
		}

		c.Locals(
			"user_id",
			uint64(userIDFloat),
		)

		c.Locals(
			"role",
			role,
		)

		return c.Next()
	}
}
