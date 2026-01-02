package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/icl00ud/goban/internal/config"
	"github.com/icl00ud/goban/internal/handlers"
	"github.com/icl00ud/goban/internal/middleware"
	"github.com/icl00ud/goban/internal/repository"
	"github.com/icl00ud/goban/internal/services"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	boardRepo := repository.NewBoardRepository(db)
	columnRepo := repository.NewColumnRepository(db)
	cardRepo := repository.NewCardRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	boardService := services.NewBoardService(boardRepo, columnRepo)
	columnService := services.NewColumnService(columnRepo, boardRepo)
	cardService := services.NewCardService(cardRepo, columnRepo, boardRepo)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(authService)
	boardHandler := handlers.NewBoardHandler(boardService)
	columnHandler := handlers.NewColumnHandler(columnService)
	cardHandler := handlers.NewCardHandler(cardService)

	// API group
	api := app.Group("/api/v1")

	// Health check (public)
	api.Get("/health", healthHandler.Check)

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)

	// Protected auth routes
	auth.Get("/me", middleware.AuthMiddleware(cfg.JWTSecret), authHandler.Me)

	// Protected routes middleware
	protected := api.Group("", middleware.AuthMiddleware(cfg.JWTSecret))

	// Board routes
	protected.Get("/boards", boardHandler.List)
	protected.Post("/boards", boardHandler.Create)
	protected.Put("/boards/reorder", boardHandler.Reorder)
	protected.Get("/boards/:id", boardHandler.Get)
	protected.Put("/boards/:id", boardHandler.Update)
	protected.Delete("/boards/:id", boardHandler.Delete)

	// Column routes
	protected.Post("/boards/:boardId/columns", columnHandler.Create)
	protected.Put("/columns/:id", columnHandler.Update)
	protected.Delete("/columns/:id", columnHandler.Delete)
	protected.Put("/columns/reorder", columnHandler.Reorder)

	// Card routes
	protected.Post("/columns/:columnId/cards", cardHandler.Create)
	protected.Get("/cards/:id", cardHandler.Get)
	protected.Put("/cards/:id", cardHandler.Update)
	protected.Delete("/cards/:id", cardHandler.Delete)
	protected.Put("/cards/:id/move", cardHandler.Move)
	protected.Put("/cards/reorder", cardHandler.Reorder)
}
