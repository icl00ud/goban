package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/icl00ud/goban/internal/handlers"
	"github.com/icl00ud/goban/internal/utils"
)

// AuthMiddleware validates JWT tokens from cookies
func AuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from cookie
		token := c.Cookies(handlers.CookieName)
		if token == "" {
			return utils.Unauthorized(c, "Authentication required")
		}

		// Validate token
		claims, err := utils.ValidateToken(token, jwtSecret)
		if err != nil {
			return utils.Unauthorized(c, "Invalid or expired token")
		}

		// Store user info in context
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}
