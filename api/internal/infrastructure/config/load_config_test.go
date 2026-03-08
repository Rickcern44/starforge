package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Scenarios(t *testing.T) {
	// Setup: Clear environment variables before each test
	clearEnv := func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("JWT_REFRESH_SECRET")
		os.Unsetenv("JWT_TTL")
	}

	t.Run("DefaultValues", func(t *testing.T) {
		clearEnv()
		config, err := LoadConfig()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if config.Server.Port != 3000 {
			t.Errorf("Expected default Port 3000, got %d", config.Server.Port)
		}
		if config.Database.ConnectionString != "postgres://postgres:password@localhost:5432/bouncy?sslmode=disable" {
			t.Errorf("Expected default ConnectionString, got '%s'", config.Database.ConnectionString)
		}
		if config.Auth.Ttl != 3600 {
			t.Errorf("Expected default Ttl 3600, got %d", config.Auth.Ttl)
		}
	})

	t.Run("EnvOverride", func(t *testing.T) {
		clearEnv()
		os.Setenv("PORT", "8080")
		os.Setenv("DATABASE_URL", "postgres://user:pass@host:5432/db")
		os.Setenv("JWT_TTL", "7200")
		defer clearEnv()

		config, err := LoadConfig()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if config.Server.Port != 8080 {
			t.Errorf("Expected overridden Port 8080, got %d", config.Server.Port)
		}
		if config.Database.ConnectionString != "postgres://user:pass@host:5432/db" {
			t.Errorf("Expected overridden ConnectionString 'postgres://user:pass@host:5432/db', got '%s'", config.Database.ConnectionString)
		}
		if config.Auth.Ttl != 7200 {
			t.Errorf("Expected overridden Ttl 7200, got %d", config.Auth.Ttl)
		}
	})

	t.Run("AppEnvMapping", func(t *testing.T) {
		clearEnv()
		os.Setenv("APP_ENV", "development")
		config, err := LoadConfig()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		if !config.IsDevelopment {
			t.Errorf("Expected IsDevelopment to be true for APP_ENV=development, got false")
		}

		clearEnv()
		os.Setenv("APP_ENV", "production")
		config, err = LoadConfig()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		if config.IsDevelopment {
			t.Errorf("Expected IsDevelopment to be false for APP_ENV=production, got true")
		}
		os.Unsetenv("APP_ENV")
	})
}
