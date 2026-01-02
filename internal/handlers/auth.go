package handlers

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/icl00ud/goban/internal/dto"
	"github.com/icl00ud/goban/internal/services"
	"github.com/icl00ud/goban/internal/utils"
)

const CookieName = "goban_token"

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	// Validate input
	if err := validateRegisterRequest(&req); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	// Register user
	user, err := h.authService.Register(&req)
	if err != nil {
		if errors.Is(err, services.ErrUserExists) {
			return utils.BadRequest(c, err.Error())
		}
		return utils.InternalError(c, "Failed to create user")
	}

	return utils.Created(c, dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		return utils.BadRequest(c, "Email and password are required")
	}

	// Authenticate user
	user, token, err := h.authService.Login(&req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return utils.Unauthorized(c, err.Error())
		}
		return utils.InternalError(c, "Login failed")
	}

	// Set HTTPOnly cookie
	c.Cookie(&fiber.Cookie{
		Name:     CookieName,
		Value:    token,
		Expires:  time.Now().Add(utils.TokenExpiration),
		HTTPOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: "Lax",
		Path:     "/",
	})

	return utils.Success(c, dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Clear the cookie
	c.Cookie(&fiber.Cookie{
		Name:     CookieName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})

	return utils.SuccessWithMessage(c, "Logged out successfully")
}

// Me returns the current authenticated user
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		return utils.NotFound(c, "User not found")
	}

	return utils.Success(c, dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}

// validateRegisterRequest validates registration request fields
func validateRegisterRequest(req *dto.RegisterRequest) error {
	req.Email = strings.TrimSpace(req.Email)
	req.Name = strings.TrimSpace(req.Name)

	if req.Email == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(req.Email, "@") {
		return errors.New("invalid email format")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if req.Name == "" {
		return errors.New("name is required")
	}
	return nil
}
