package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the CLI configuration.
type Config struct {
	Server    string `json:"server"`
	Token     string `json:"token"`
	ProjectID int    `json:"project_id"`
}

const (
	configFileName = ".yapi.json"

	defaultServer = "http://127.0.0.1:3000"
)

// LoadConfig reads configuration from files and environment variables.
// Priority: env > ./.yapi.json > ~/.yapi.json
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	// Load from home directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		homePath := filepath.Join(homeDir, configFileName)
		loadFromFile(homePath, cfg)
	}

	// Load from current directory (overrides home)
	loadFromFile(configFileName, cfg)

	// Environment variables override file config
	if v := os.Getenv("YAPI_SERVER"); v != "" {
		cfg.Server = v
	}
	if v := os.Getenv("YAPI_TOKEN"); v != "" {
		cfg.Token = v
	}
	if v := os.Getenv("YAPI_PROJECT_ID"); v != "" {
		var id int
		if _, err := fmt.Sscanf(v, "%d", &id); err == nil {
			cfg.ProjectID = id
		}
	}

	if cfg.Server == "" {
		cfg.Server = defaultServer
	}

	return cfg, nil
}

// SaveConfig writes configuration to a file.
func SaveConfig(cfg *Config, local bool) error {
	var path string
	if local {
		path = configFileName
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cannot determine home directory: %w", err)
		}
		path = filepath.Join(homeDir, configFileName)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write config file %s: %w", path, err)
	}
	return nil
}

func loadFromFile(path string, cfg *Config) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var fileCfg Config
	if err := json.Unmarshal(data, &fileCfg); err != nil {
		return
	}
	if fileCfg.Server != "" {
		cfg.Server = fileCfg.Server
	}
	if fileCfg.Token != "" {
		cfg.Token = fileCfg.Token
	}
	if fileCfg.ProjectID != 0 {
		cfg.ProjectID = fileCfg.ProjectID
	}
}
