package config

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
)

var defaultConfig = Config{
	Database: &DatabaseConfig{
		Password: os.Getenv("POSTGRES_PASS"),
	},
	Auth: &AuthConfig{
		JwtSecret:     os.Getenv("JWT_SECRET"),
		RefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
	},
}

// LoadConfig load a global configuration based off an input path
func LoadConfig(configFilePath string) (*Config, error) {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, errors.New("config file does not exist")
	}

	config := defaultConfig

	if _, err := toml.DecodeFile(configFilePath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
