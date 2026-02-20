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

func TestGetAgentConfig(t *testing.T) {
	tomlContent := `
[agent]
model = "default-model"
temperature = 0.5
max_iterations = 10

[agents.lisa]
model = "claude-opus-4-1"
temperature = 0.3
max_iterations = 15

[agents.ralph]
model = "gpt-4o-mini"
fallback_model = "gemini-2.0-flash"
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

	tests := []struct {
		agent     string
		wantModel string
		wantTemp  float64
	}{
		{"lisa", "claude-opus-4-1", 0.3},
		{"ralph", "gpt-4o-mini", 0.5},
		{"bart", "default-model", 0.5},
		{"LISA", "claude-opus-4-1", 0.3}, // Test case-insensitivity
	}

	for _, tt := range tests {
		t.Run(tt.agent, func(t *testing.T) {
			agentCfg := cfg.GetAgentConfig(tt.agent)
			if agentCfg.Model != tt.wantModel {
				t.Errorf("GetAgentConfig(%s).Model = %s, want %s", tt.agent, agentCfg.Model, tt.wantModel)
			}
			if agentCfg.Temperature != tt.wantTemp {
				t.Errorf("GetAgentConfig(%s).Temperature = %f, want %f", tt.agent, agentCfg.Temperature, tt.wantTemp)
			}
		})
	}

	// Test fallback model for Ralph
	ralphCfg := cfg.GetAgentConfig("ralph")
	if ralphCfg.FallbackModel != "gemini-2.0-flash" {
		t.Errorf("Expected Ralph fallback_model to be gemini-2.0-flash, got %s", ralphCfg.FallbackModel)
	}
}
