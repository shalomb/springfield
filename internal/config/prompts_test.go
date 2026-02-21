package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadPrompt verifies that a prompt can be loaded from a markdown file.
func TestLoadPrompt(t *testing.T) {
	// Create a temporary directory for testing.
	tmpDir := t.TempDir()

	// Write a test prompt file.
	testPromptFile := filepath.Join(tmpDir, "prompt_test.md")
	testPromptContent := "This is a test prompt"
	err := os.WriteFile(testPromptFile, []byte(testPromptContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test prompt file: %v", err)
	}

	// Test LoadPrompt function.
	content, err := LoadPrompt(testPromptFile)
	if err != nil {
		t.Fatalf("LoadPrompt failed: %v", err)
	}

	if content != testPromptContent {
		t.Errorf("Expected prompt content '%s', got '%s'", testPromptContent, content)
	}

	// Test with YAML front matter
	testPromptWithFrontMatter := `---
name: test
role: tester
---
Actual content`
	err = os.WriteFile(testPromptFile, []byte(testPromptWithFrontMatter), 0644)
	if err != nil {
		t.Fatalf("Failed to write test prompt file: %v", err)
	}

	content, err = LoadPrompt(testPromptFile)
	if err != nil {
		t.Fatalf("LoadPrompt failed: %v", err)
	}

	expectedContent := "Actual content"
	if content != expectedContent {
		t.Errorf("Expected prompt content '%s', got '%s'", expectedContent, content)
	}
}

// TestLoadPrompt_FileNotFound verifies error handling for missing files.
func TestLoadPrompt_FileNotFound(t *testing.T) {
	_, err := LoadPrompt("/nonexistent/path/prompt_nonexistent.md")
	if err == nil {
		t.Errorf("Expected error for missing file, got nil")
	}
}

// TestLoadPromptNames verifies that all required agent prompts can be located.
func TestLoadPromptNames(t *testing.T) {
	prompts := []string{"ralph", "lisa", "bart", "lovejoy"}

	for _, agent := range prompts {
		t.Run(agent, func(t *testing.T) {
			// Just verify the function doesn't panic.
			_ = GetPromptPath(agent)
		})
	}
}

// TestAllPromptsExist verifies that all agent prompt files exist.
func TestAllPromptsExist(t *testing.T) {
	agents := []string{"ralph", "lisa", "bart", "lovejoy"}

	for _, agent := range agents {
		t.Run(agent, func(t *testing.T) {
			path := GetPromptPath(agent)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("Prompt file does not exist: %s", path)
			}
		})
	}
}
