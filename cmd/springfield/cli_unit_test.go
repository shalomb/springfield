package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRootCmd_Help(t *testing.T) {
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"--help"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if !bytes.Contains(b.Bytes(), []byte("Springfield is an AI agent orchestration tool")) {
		t.Errorf("unexpected help output: %s", b.String())
	}
}

func TestRootCmd_NoArgs(t *testing.T) {
	// Should show help if no agent/task
	rootCmd.SetArgs([]string{})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestRootCmd_RunMock(t *testing.T) {
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "config.toml")
	_ = os.WriteFile(confPath, []byte("[axon]\nversion=\"1.0.0\"\n"), 0644)

	t.Setenv("USE_MOCK_LLM", "true")
	t.Setenv("SPRINGFIELD_CONFIG", confPath)

	// Reset global flags because cobra doesn't reset them between Execute calls
	// Use lisa instead of ralph to avoid multi-iteration loop in tests
	agentName = "lisa"
	task = "test task"
	configPath = ""

	err := rootCmd.RunE(rootCmd, []string{})
	if err != nil {
		t.Fatalf("RunE failed: %v", err)
	}
}

func TestRootCmd_Roles(t *testing.T) {
	t.Setenv("USE_MOCK_LLM", "true")
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "config.toml")
	_ = os.WriteFile(confPath, []byte("[axon]\nversion=\"1.0.0\"\n"), 0644)
	t.Setenv("SPRINGFIELD_CONFIG", confPath)

	// Test case-insensitive agent name matching
	// These agents use BaseRunner (single call) and don't require special setup
	validRoles := []string{"marge", "lisa", "bart", "lovejoy"}
	for _, name := range validRoles {
		agentName = name
		task = "test"
		err := rootCmd.RunE(rootCmd, []string{})
		if err != nil {
			t.Errorf("RunE failed for valid role %s: %v", name, err)
		}
	}

	// Test case normalization: MARGE (uppercase) should work like marge
	agentName = "MARGE"
	task = "test"
	err := rootCmd.RunE(rootCmd, []string{})
	if err != nil {
		t.Errorf("RunE failed for case-normalized MARGE: %v", err)
	}

	// Test invalid agent (should fail)
	agentName = "other"
	task = "test"
	err = rootCmd.RunE(rootCmd, []string{})
	if err == nil {
		t.Error("expected error for invalid agent 'other', got nil")
	}
}

func TestRootCmd_MissingAgent(t *testing.T) {
	agentName = ""
	task = "test"
	_ = rootCmd.RunE(rootCmd, []string{})
}

func TestRootCmd_RunError(t *testing.T) {
	t.Setenv("USE_MOCK_LLM", "true")
	t.Setenv("MOCK_LLM_ERROR", "true")

	agentName = "ralph"
	task = "test"

	err := rootCmd.RunE(rootCmd, []string{})
	if err == nil {
		t.Error("expected error from mock llm, got nil")
	}
}

func TestRunMain(t *testing.T) {
	agentName = ""
	task = ""
	rootCmd.SetArgs([]string{"--help"})
	if err := runMain(); err != nil {
		t.Fatalf("runMain failed: %v", err)
	}
}
