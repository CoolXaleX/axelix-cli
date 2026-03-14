package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds all named services and the name of the currently active one.
type Config struct {
	Services map[string]string `json:"services"`
	Current  string            `json:"current"`
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

// Resolve returns the effective URL for the given session.
// Priority: --url flag > AXELIX_URL env > --service flag > current service in config file.
func Resolve(flagURL, flagService string) (string, error) {
	if flagURL != "" {
		return flagURL, nil
	}
	if envURL := os.Getenv("AXELIX_URL"); envURL != "" {
		return envURL, nil
	}
	cfg, err := Load()
	if err != nil {
		return "", err
	}
	serviceName := flagService
	if serviceName == "" {
		serviceName = cfg.Current
	}
	if serviceName == "" {
		return "", fmt.Errorf("no service selected — run 'axelix config use <name>' or use --url / --service")
	}
	url, ok := cfg.Services[serviceName]
	if !ok {
		return "", fmt.Errorf("service %q not found in config — run 'axelix config list' to see available services", serviceName)
	}
	return url, nil
}
