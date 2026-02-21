package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// LoadPrompt reads a prompt from a markdown file and returns its content.
// It supports YAML front matter and strips it if present.
func LoadPrompt(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file %s: %w", path, err)
	}

	s := string(content)
	if strings.HasPrefix(s, "---") {
		// Find the end of the front matter
		parts := strings.SplitN(s, "---", 3)
		if len(parts) == 3 {
			// parts[0] is empty, parts[1] is front matter, parts[2] is markdown
			return strings.TrimSpace(parts[2]), nil
		}
	}

	return s, nil
}

// GetPromptPath returns the path to a prompt markdown file for the given agent.
// It prioritizes global (~/.pi/) then local (.pi/) directories.
func GetPromptPath(agent string) string {
	home, _ := os.UserHomeDir()
	cwd, _ := os.Getwd()

	// Define search directories in order of priority
	searchDirs := []string{}

	// 1. Global User Directory (~/.pi)
	if home != "" {
		searchDirs = append(searchDirs,
			filepath.Join(home, ".pi", "prompts"),
			filepath.Join(home, ".pi", "agents"),
		)
	}

	// 2. Project Local Directory (.pi)
	// Try to find project root by looking for .git
	projectRoot := cwd
	tempDir := cwd
	for {
		if _, err := os.Stat(filepath.Join(tempDir, ".git")); err == nil {
			projectRoot = tempDir
			break
		}
		parent := filepath.Dir(tempDir)
		if parent == tempDir {
			break
		}
		tempDir = parent
	}
	searchDirs = append(searchDirs,
		filepath.Join(projectRoot, ".pi", "prompts"),
		filepath.Join(projectRoot, ".pi", "agents"),
	)

	// Try specific patterns
	patterns := []string{
		agent + ".prompt.md",
		agent + ".prompt",
		"prompt_" + agent + ".md",
		agent + ".md",
	}

	for _, dir := range searchDirs {
		for _, pattern := range patterns {
			path := filepath.Join(dir, pattern)
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}

	// Fallback to local .pi/agents
	return filepath.Join(projectRoot, ".pi", "agents", "prompt_"+agent+".md")
}

// Config holds the Springfield configuration.
type Config struct {
	Agent   AgentConfig            `toml:"agent"`
	Agents  map[string]AgentConfig `toml:"agents"`
	Sandbox SandboxConfig          `toml:"sandbox"`
	Env     map[string]string      `toml:"env"`
}

// AgentConfig holds agent-specific settings.
type AgentConfig struct {
	// Model specification can be:
	// - "claude-opus-4-1" (uses default provider)
	// - "anthropic/claude-opus-4-1" (explicit provider)
	// - "openai/gpt-4o" (explicit provider)
	// - "google-gemini-cli/gemini-2.0-flash" (explicit provider)
	Model         string `toml:"model"`          // Default model or primary model
	PrimaryModel  string `toml:"primary_model"`  // Override for primary model (can include provider)
	FallbackModel string `toml:"fallback_model"` // Fallback model (can include provider)
	MaxIterations int    `toml:"max_iterations"`
	Budget        int    `toml:"budget"`
}

// SandboxConfig holds sandbox/Axon-specific settings.
type SandboxConfig struct {
	Image        string `toml:"image"`
	ImageBuilder string `toml:"image_builder"`
}

// LoadConfig loads the configuration from a .springfield.toml or config.toml file in the given directory.
// It searches in the specified directory, then in XDG config paths.
func LoadConfig(dir string) (*Config, error) {
	cfg := &Config{
		Agent: AgentConfig{
			Model:         "gemini-2.0-flash", // Default model
			MaxIterations: 20,
		},
		Agents: make(map[string]AgentConfig),
		Sandbox: SandboxConfig{
			Image:        "docker.io/library/debian:trixie-slim",
			ImageBuilder: "podman",
		},
		Env: map[string]string{
			// Git interaction: disable interactive prompts and pagers
			"GIT_EDITOR":          "true",
			"GIT_PAGER":           "cat",
			"GIT_ASKPASS":         "false",
			"GIT_TERMINAL_PROMPT": "0",

			// Automation & CI/CD signals
			"CI":               "true",
			"TF_IN_AUTOMATION": "true",

			// Disable color output (for log parsing and clarity)
			"NO_COLOR": "1",
			"TERM":     "dumb",
			"PAGER":    "cat",

			// Prevent interactive prompts from other tools
			"PROMPT_COMMAND": "",
		},
	}

	// Define search paths in order of priority
	home, _ := os.UserHomeDir()
	configDirs := []string{dir}
	if home != "" {
		configDirs = append(configDirs, filepath.Join(home, ".config", "springfield"))
	}

	// Configuration file names to look for
	configNames := []string{".springfield.toml", "config.toml", "config.yaml", "config.yml"}

	var path string
	for _, d := range configDirs {
		for _, name := range configNames {
			p := filepath.Join(d, name)
			if _, err := os.Stat(p); err == nil {
				path = p
				break
			}
		}
		if path != "" {
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
	if agentConfig.MaxIterations == 0 {
		agentConfig.MaxIterations = c.Agent.MaxIterations
	}
	if agentConfig.Budget == 0 {
		agentConfig.Budget = c.Agent.Budget
	}
	return agentConfig
}
