package main

import (
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	goban "github.com/icl00ud/goban"
	"github.com/icl00ud/goban/internal/config"
	"github.com/icl00ud/goban/internal/database"
	"github.com/icl00ud/goban/internal/router"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173,http://localhost:8080",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Setup API routes
	router.Setup(app, db, cfg)

	// Setup static file serving with SPA fallback
	setupStaticServing(app)

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupStaticServing configures static file serving from embedded files with SPA fallback
func setupStaticServing(app *fiber.App) {
	// Get the embedded filesystem, stripping the "web/dist" prefix
	distFS, err := fs.Sub(goban.StaticFiles, "web/dist")
	if err != nil {
		log.Printf("Warning: Could not load embedded static files: %v", err)
		return
	}

	// Serve static assets (js, css, images, etc.)
	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:   http.FS(distFS),
		Browse: false,
	}))

	// SPA fallback: serve index.html for all non-API routes
	app.Use(func(c *fiber.Ctx) error {
		path := c.Path()

		// Skip API routes
		if strings.HasPrefix(path, "/api") {
			return c.Next()
		}

		// Try to serve the requested file
		file, err := distFS.Open(strings.TrimPrefix(path, "/"))
		if err == nil {
			file.Close()
			return c.Next()
		}

		// Fallback to index.html for SPA routing
		indexHTML, err := fs.ReadFile(distFS, "index.html")
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Not found")
		}

		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.Send(indexHTML)
	})

	// Serve static files
	app.Use(filesystem.New(filesystem.Config{
		Root:   http.FS(distFS),
		Browse: false,
	}))
}
