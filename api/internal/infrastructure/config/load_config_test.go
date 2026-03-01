package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_Scenarios(t *testing.T) {
	// Setup: Clear environment variables before each test
	clearEnv := func() {
		os.Unsetenv("APP_CONFIG_PATH")
		os.Unsetenv("PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_USER")
		os.Unsetenv("POSTGRES_PASS")
		os.Unsetenv("DB_NAME")
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
		if config.Database.Host != "localhost" {
			t.Errorf("Expected default Host 'localhost', got '%s'", config.Database.Host)
		}
	})

	t.Run("EnvOverride", func(t *testing.T) {
		clearEnv()
		os.Setenv("PORT", "8080")
		os.Setenv("DB_HOST", "db.example.com")
		defer clearEnv()

		config, err := LoadConfig()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if config.Server.Port != 8080 {
			t.Errorf("Expected overridden Port 8080, got %d", config.Server.Port)
		}
		if config.Database.Host != "db.example.com" {
			t.Errorf("Expected overridden Host 'db.example.com', got '%s'", config.Database.Host)
		}
	})

	t.Run("FileOverride", func(t *testing.T) {
		clearEnv()

		// Create a temporary TOML file
		tmpDir := t.TempDir()
		tmpFilePath := filepath.Join(tmpDir, "config.toml")
		content := `
[server]
port = 9000

[database]
host = "file-host"
`
		if err := os.WriteFile(tmpFilePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write tmp config file: %v", err)
		}

		os.Setenv("APP_CONFIG_PATH", tmpFilePath)
		defer clearEnv()

		config, err := LoadConfig()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if config.Server.Port != 9000 {
			t.Errorf("Expected file-overridden Port 9000, got %d", config.Server.Port)
		}
		if config.Database.Host != "file-host" {
			t.Errorf("Expected file-overridden Host 'file-host', got '%s'", config.Database.Host)
		}
	})

	t.Run("EnvWinsOverFile", func(t *testing.T) {
		clearEnv()

		// Setup file
		tmpDir := t.TempDir()
		tmpFilePath := filepath.Join(tmpDir, "config.toml")
		content := `[server]
port = 9000`
		os.WriteFile(tmpFilePath, []byte(content), 0644)

		// Setup env (should win)
		os.Setenv("APP_CONFIG_PATH", tmpFilePath)
		os.Setenv("PORT", "9999")
		defer clearEnv()

		config, err := LoadConfig()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if config.Server.Port != 9999 {
			t.Errorf("Expected Env to win (9999), but got %d (File might have won?)", config.Server.Port)
		}
	})
}
