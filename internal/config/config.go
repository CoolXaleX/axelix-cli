package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds all named services.
type Config struct {
	Services map[string]string `json:"services"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".axelix", "config.json"), nil
}

// Load reads the config from ~/.axelix/config.json.
// Returns an empty config (not an error) if the file does not exist yet.
func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Services: map[string]string{}}, nil

		}
		return nil, err
	}
	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	if cfg.Services == nil {
		cfg.Services = map[string]string{}
	}
	return cfg, nil
}

// Save writes the config to ~/.axelix/config.json.
func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

