package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the CLI/bot configuration stored in ~/.config/goban/config.json
type Config struct {
	ServerURL             string  `json:"server_url"`
	Token                 string  `json:"token"`
	TelegramBotToken      string  `json:"telegram_bot_token"`
	TelegramAllowedUsers  []int64 `json:"telegram_allowed_users"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "goban", "config.json"), nil
}

// Load reads config from disk, overlaying environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		ServerURL: "http://localhost:8080",
	}

	path, err := configPath()
	if err == nil {
		if data, err := os.ReadFile(path); err == nil {
			_ = json.Unmarshal(data, cfg)
		}
	}

	// Environment variables override file values
	if v := os.Getenv("GOBAN_SERVER_URL"); v != "" {
		cfg.ServerURL = v
	}
	if v := os.Getenv("GOBAN_TOKEN"); v != "" {
		cfg.Token = v
	}
	if v := os.Getenv("TELEGRAM_BOT_TOKEN"); v != "" {
		cfg.TelegramBotToken = v
	}

	return cfg, nil
}

// Save persists the config to disk.
func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
