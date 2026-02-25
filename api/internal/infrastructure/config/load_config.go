package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

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
	config := defaultConfig

	if configFilePath == "" {
		return &config, nil
	}

	// If the path is relative, resolve it from the project root
	if !filepath.IsAbs(configFilePath) {
		root, err := findProjectRoot()
		if err == nil {
			configFilePath = filepath.Join(root, configFilePath)
		}
	}

	info, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist at: %s", configFilePath)
	}
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, fmt.Errorf("config path is a directory, not a file: %s", configFilePath)
	}

	if _, err := toml.DecodeFile(configFilePath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("could not find project root (go.mod)")
		}
		dir = parent
	}
}
