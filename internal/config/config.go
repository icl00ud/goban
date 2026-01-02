package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DBDriver    string
	DatabaseURL string
	JWTSecret   string
}

func Load() *Config {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DBDriver:    getEnv("DB_DRIVER", "sqlite"),
		DatabaseURL: getEnv("DATABASE_URL", "./goban.db"),
		JWTSecret:   getEnv("JWT_SECRET", "default-secret-change-me"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
