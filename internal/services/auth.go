package services

import (
	"errors"

	"github.com/icl00ud/goban/internal/dto"
	"github.com/icl00ud/goban/internal/models"
	"github.com/icl00ud/goban/internal/repository"
	"github.com/icl00ud/goban/internal/utils"
)

var (
	ErrUserExists       = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user account
func (s *AuthService) Register(req *dto.RegisterRequest) (*models.User, error) {
	// Check if user already exists
	if s.userRepo.ExistsByEmail(req.Email) {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Name:         req.Name,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(req *dto.LoginRequest) (*models.User, string, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// GetUserByID retrieves a user by their ID
func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}
