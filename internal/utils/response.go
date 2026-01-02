package utils

import "github.com/gofiber/fiber/v2"

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success sends a success response
func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(Response{
		Success: true,
		Data:    data,
	})
}

// SuccessWithMessage sends a success response with a message
func SuccessWithMessage(c *fiber.Ctx, message string) error {
	return c.JSON(Response{
		Success: true,
		Message: message,
	})
}

// Created sends a 201 created response
func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Data:    data,
	})
}

// Error sends an error response
func Error(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(Response{
		Success: false,
		Error:   message,
	})
}

// BadRequest sends a 400 error response
func BadRequest(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusBadRequest, message)
}

// Unauthorized sends a 401 error response
func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, message)
}

// Forbidden sends a 403 error response
func Forbidden(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusForbidden, message)
}

// NotFound sends a 404 error response
func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message)
}

// InternalError sends a 500 error response
func InternalError(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusInternalServerError, message)
}
