package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/icl00ud/goban/internal/utils"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	return utils.Success(c, fiber.Map{
		"status":  "healthy",
		"service": "goban",
	})
}
