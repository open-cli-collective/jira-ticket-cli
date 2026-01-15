package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	configDirName  = "jira-ticket-cli"
	configFileName = "config.json"
	configFileMode = 0600
	configDirMode  = 0700
)

// Config holds the CLI configuration
type Config struct {
	Domain   string `json:"domain"`
	Email    string `json:"email"`
	APIToken string `json:"api_token"`
}

// configPath returns the path to the config file
func configPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}
	return filepath.Join(configDir, configDirName, configFileName), nil
}

// Load loads the configuration from file
func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// Save saves the configuration to file
func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, configDirMode); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, configFileMode); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Clear removes the configuration file
func Clear() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to remove config file: %w", err)
	}

	return nil
}

// GetDomain returns the domain from config or environment
func GetDomain() string {
	if v := os.Getenv("JIRA_DOMAIN"); v != "" {
		return v
	}
	cfg, err := Load()
	if err != nil {
		return ""
	}
	return cfg.Domain
}

// GetEmail returns the email from config or environment
func GetEmail() string {
	if v := os.Getenv("JIRA_EMAIL"); v != "" {
		return v
	}
	cfg, err := Load()
	if err != nil {
		return ""
	}
	return cfg.Email
}

// GetAPIToken returns the API token from config or environment
func GetAPIToken() string {
	if v := os.Getenv("JIRA_API_TOKEN"); v != "" {
		return v
	}
	cfg, err := Load()
	if err != nil {
		return ""
	}
	return cfg.APIToken
}

// IsConfigured returns true if all required config values are set
func IsConfigured() bool {
	return GetDomain() != "" && GetEmail() != "" && GetAPIToken() != ""
}

// Path returns the path to the config file
func Path() string {
	path, _ := configPath()
	return path
}
