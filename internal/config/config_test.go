package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tomlContent := `
[agent]
model = "gpt-4"
temperature = 0.5
max_iterations = 10
`
	err := os.WriteFile(".springfield.toml", []byte(tomlContent), 0644)
	if err != nil {
		t.Fatalf("failed to create temp config: %v", err)
	}
	defer os.Remove(".springfield.toml")

	cfg, err := LoadConfig(".")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Agent.Model != "gpt-4" {
		t.Errorf("expected model gpt-4, got %s", cfg.Agent.Model)
	}
	if cfg.Agent.Temperature != 0.5 {
		t.Errorf("expected temperature 0.5, got %f", cfg.Agent.Temperature)
	}
	if cfg.Agent.MaxIterations != 10 {
		t.Errorf("expected max_iterations 10, got %d", cfg.Agent.MaxIterations)
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	// No file
	cfg, err := LoadConfig("non-existent")
	if err != nil {
		t.Fatalf("LoadConfig should not fail on missing file, got %v", err)
	}

	if cfg.Agent.Model == "" {
		t.Error("expected default model to be set")
	}
}
