package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
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
	projectRoot, _ := getRootDirectory()

	// 1. Load .env file if present (sets env vars for the process)
	if projectRoot != "" {
		_ = godotenv.Load(filepath.Join(projectRoot, ".env"))
	} else {
		_ = godotenv.Load()
	}

	// 2. Load Defaults
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

	// 3. Optional TOML override
	configPath := os.Getenv("APP_CONFIG_PATH")
	if configPath != "" {
		if !filepath.IsAbs(configPath) && projectRoot != "" {
			configPath = filepath.Join(projectRoot, configPath)
		}
		if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
			return nil, fmt.Errorf("failed to load config file at %s: %w", configPath, err)
		}
	}

	// 4. Environment variables (Highest priority)
	// We map flat env vars (e.g. DB_PORT) to nested struct fields (database.port)
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
		// Fallback: lowercase and replace underscores for others (e.g. SERVER_PORT -> server.port)
		return strings.ReplaceAll(strings.ToLower(s), "_", ".")
	}), nil)
	if err != nil {
		return nil, fmt.Errorf("error loading env vars: %w", err)
	}

	// Unmarshal into the struct
	config := &Config{}
	if err := k.Unmarshal("", config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return config, nil
}

func getRootDirectory() (string, error) {
	dir, err := os.Getwd()
	marker := "go.mod"

	if err != nil {
		slog.Error(err.Error())
	}
	for {
		// Check if the marker file exists in the current directory
		markerPath := filepath.Join(dir, marker)
		if _, err := os.Stat(markerPath); err == nil {
			return dir, nil // Found the project root
		}

		// Move one directory up
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root without finding marker
			return "", fmt.Errorf("project root not found (no %s)", marker)
		}
		dir = parent
	}
}
