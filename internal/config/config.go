package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config holds the Springfield configuration.
type Config struct {
	Agent   AgentConfig            `toml:"agent"`
	Agents  map[string]AgentConfig `toml:"agents"`
	Sandbox SandboxConfig          `toml:"sandbox"`
}

// AgentConfig holds agent-specific settings.
type AgentConfig struct {
	// Model specification can be:
	// - "claude-opus-4-1" (uses default provider)
	// - "anthropic/claude-opus-4-1" (explicit provider)
	// - "openai/gpt-4o" (explicit provider)
	// - "google-gemini-cli/gemini-2.0-flash" (explicit provider)
	Model         string  `toml:"model"`         // Default model or primary model
	PrimaryModel  string  `toml:"primary_model"` // Override for primary model (can include provider)
	FallbackModel string  `toml:"fallback_model"` // Fallback model (can include provider)
	Temperature   float64 `toml:"temperature"`
	MaxIterations int     `toml:"max_iterations"`
	Budget        int     `toml:"budget"`
}

// SandboxConfig holds sandbox/Axon-specific settings.
type SandboxConfig struct {
	Image       string `toml:"image"`
	ImageBuilder string `toml:"image_builder"`
}

// LoadConfig loads the configuration from a .springfield.toml or config.toml file in the given directory.
func LoadConfig(dir string) (*Config, error) {
	cfg := &Config{
		Agent: AgentConfig{
			Model:         "gemini-2.0-flash", // Default model
			Temperature:   0.7,
			MaxIterations: 20,
		},
		Agents: make(map[string]AgentConfig),
		Sandbox: SandboxConfig{
			Image:        "docker.io/library/debian:trixie-slim",
			ImageBuilder: "podman",
		},
	}

	// Try .springfield.toml first, then fall back to config.toml
	paths := []string{
		filepath.Join(dir, ".springfield.toml"),
		filepath.Join(dir, "config.toml"),
	}

	var path string
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			path = p
			break
		}
	}

	if path == "" {
		return cfg, nil // No config file, use defaults
	}

	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file %s: %w", path, err)
	}

	return cfg, nil
}

// GetAgentConfig returns the configuration for a specific agent.
// Falls back to the default Agent config if no agent-specific config exists.
func (c *Config) GetAgentConfig(agentName string) AgentConfig {
	agentName = strings.ToLower(agentName)

	// Check for agent-specific config
	if agentConfig, ok := c.Agents[agentName]; ok {
		return c.mergeWithDefaults(agentConfig)
	}

	// Fall back to default agent config
	return c.Agent
}

// mergeWithDefaults fills in any missing values from the default agent config.
func (c *Config) mergeWithDefaults(agentConfig AgentConfig) AgentConfig {
	if agentConfig.Model == "" && agentConfig.PrimaryModel == "" {
		agentConfig.Model = c.Agent.Model
	}
	if agentConfig.FallbackModel == "" {
		agentConfig.FallbackModel = c.Agent.FallbackModel
	}
	if agentConfig.Temperature == 0 {
		agentConfig.Temperature = c.Agent.Temperature
	}
	if agentConfig.MaxIterations == 0 {
		agentConfig.MaxIterations = c.Agent.MaxIterations
	}
	if agentConfig.Budget == 0 {
		agentConfig.Budget = c.Agent.Budget
	}
	return agentConfig
}
