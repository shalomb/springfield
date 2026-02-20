package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds the Springfield configuration.
type Config struct {
	Agent AgentConfig `toml:"agent"`
}

// AgentConfig holds agent-specific settings.
type AgentConfig struct {
	Model         string  `toml:"model"` // Keep for backward compatibility or as Primary
	PrimaryModel  string  `toml:"primary_model"`
	FallbackModel string  `toml:"fallback_model"`
	Temperature   float64 `toml:"temperature"`
	MaxIterations int     `toml:"max_iterations"`
	Budget        int     `toml:"budget"`
}

// LoadConfig loads the configuration from a .springfield.toml file in the given directory.
func LoadConfig(dir string) (*Config, error) {
	cfg := &Config{
		Agent: AgentConfig{
			Model:         "gemini-2.0-flash", // Default model
			Temperature:   0.7,
			MaxIterations: 20,
		},
	}

	path := filepath.Join(dir, ".springfield.toml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return cfg, nil
}
