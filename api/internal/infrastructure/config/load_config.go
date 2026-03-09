package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// LoadConfig loads configuration using this hierarchy (last loaded wins):
// 1. .env file (optional, sets env vars for the process)
// 2. Environment variables (highest priority)
func LoadConfig() (*Config, error) {
	// Load .env file if present
	_ = godotenv.Load()

	// Initialize config struct with pointers to child structs
	config := &Config{
		Server:   &ServerConfig{},
		Database: &DatabaseConfig{},
		Auth:     &AuthConfig{},
	}

	// Parse environment variables into the struct
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("error parsing environment variables: %w", err)
	}

	config.IsDevelopment = config.AppEnv == "development"

	return config, nil
}
