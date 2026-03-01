package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

// LoadConfig loads configuration using this hierarchy (last loaded wins):
// 1. Code-defined defaults
// 2. TOML file (optional, only if path is provided)
// 3. Environment variables (highest priority)
func LoadConfig() (*Config, error) {
	// Load .env file if present (sets env vars for the process)
	_ = godotenv.Load()

	// 1. Load Defaults
	_ = k.Load(confmap.Provider(map[string]interface{}{
		"server.port":         3000,
		"database.host":       "localhost",
		"database.port":       "5432",
		"database.username":   "postgres",
		"database.password":   "password",
		"database.database":   "bouncy",
		"auth.jwt_secret":     "super-secret",
		"auth.refresh_secret": "refresh-secret",
		"auth.jwt_ttl":        3600,
	}, "."), nil)

	// 2. Optional TOML override
	configPath := os.Getenv("APP_CONFIG_PATH")
	if configPath != "" {
		if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
			return nil, fmt.Errorf("failed to load config file at %s: %w", configPath, err)
		}
	}

	// 3. Environment variables (Highest priority)
	err := k.Load(env.Provider("", ".", func(s string) string {
		// Custom mapping for legacy environment variables
		mapping := map[string]string{
			"PORT":               "server.port",
			"DB_HOST":            "database.host",
			"DB_PORT":            "database.port",
			"DB_USER":            "database.username",
			"POSTGRES_PASS":      "database.password",
			"DB_NAME":            "database.database",
			"JWT_SECRET":         "auth.jwt_secret",
			"JWT_REFRESH_SECRET": "auth.refresh_secret",
			"JWT_TTL":            "auth.jwt_ttl",
		}
		if mapped, ok := mapping[s]; ok {
			return mapped
		}
		return strings.ReplaceAll(strings.ToLower(s), "_", ".")
	}), nil)
	if err != nil {
		return nil, fmt.Errorf("error loading env vars: %w", err)
	}

	config := &Config{}
	if err := k.Unmarshal("", config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return config, nil
}
