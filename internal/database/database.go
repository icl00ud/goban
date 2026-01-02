package database

import (
	"fmt"
	"log"

	"github.com/icl00ud/goban/internal/config"
	"github.com/icl00ud/goban/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.DBDriver {
	case "postgres":
		dialector = postgres.Open(cfg.DatabaseURL)
	case "sqlite":
		dialector = sqlite.Open(cfg.DatabaseURL)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.DBDriver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Printf("Connected to %s database successfully", cfg.DBDriver)
	return db, nil
}

func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Board{},
		&models.Column{},
		&models.Card{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
