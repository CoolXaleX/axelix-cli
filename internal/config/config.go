package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds connection settings for the Axelix SBS app.
type Config struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".axelix", "config.json"), nil
}

// Load reads the config from ~/.axelix/config.json.
func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
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

// Resolve merges flag values, env vars, and the config file.
// Priority: flag > env var > file.
func Resolve(flagURL, flagUser, flagPass string) *Config {
	file, _ := Load()
	if file == nil {
		file = &Config{}
	}

	url := flagURL
	if url == "" {
		url = os.Getenv("AXELIX_URL")
	}
	if url == "" {
		url = file.URL
	}

	user := flagUser
	if user == "" {
		user = os.Getenv("AXELIX_USER")
	}
	if user == "" {
		user = file.Username
	}

	pass := flagPass
	if pass == "" {
		pass = os.Getenv("AXELIX_PASSWORD")
	}
	if pass == "" {
		pass = file.Password
	}

	return &Config{URL: url, Username: user, Password: pass}
}
