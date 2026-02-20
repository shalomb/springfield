package agent

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLovejoyRunnerVerifiesToDOEmpty verifies Lovejoy checks if TODO.md is empty.
func TestLovejoyRunnerVerifiesToDOEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Initialize git
	initGitRepo(t, tmpDir)

	// Create prompt
	promptPath := filepath.Join(tmpDir, "test_prompt.md")
	os.WriteFile(promptPath, []byte("You are Lovejoy"), 0644)

	// No TODO.md - ready for release
	mock := &mockLLMClient{responses: []string{"Release ceremony complete"}}
	lr := &LovejoyRunner{
		BaseRunner: &BaseRunner{
			Agent:      "lovejoy",
			LLMClient:  mock,
			PromptPath: promptPath,
		},
	}

	err := lr.Run(context.Background())
	if err != nil {
		t.Errorf("Run() returned unexpected error: %v", err)
	}

	if mock.calls == 0 {
		t.Fatal("Expected LLM Chat to be called")
	}
}

// TestLovejoyRunnerChecksReleaseReadiness verifies release readiness context.
func TestLovejoyRunnerChecksReleaseReadiness(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Initialize git
	initGitRepo(t, tmpDir)

	// Create prompt
	promptPath := filepath.Join(tmpDir, "test_prompt.md")
	os.WriteFile(promptPath, []byte("You are Lovejoy"), 0644)

	mock := &mockLLMClient{responses: []string{"Ready to release"}}
	lr := &LovejoyRunner{
		BaseRunner: &BaseRunner{
			Agent:      "lovejoy",
			LLMClient:  mock,
			PromptPath: promptPath,
		},
	}

	err := lr.Run(context.Background())
	if err != nil {
		t.Errorf("Run() returned unexpected error: %v", err)
	}

	messages := mock.received[0]
	userMessage := messages[len(messages)-1]

	// Should include readiness information
	if userMessage.Content == "" {
		t.Errorf("Expected non-empty user message with readiness context")
	}
}

// TestLovejoyRunnerIncludesCHANGELOGContext verifies CHANGELOG.md is included if present.
func TestLovejoyRunnerIncludesCHANGELOGContext(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Initialize git
	initGitRepo(t, tmpDir)

	// Create prompt
	promptPath := filepath.Join(tmpDir, "test_prompt.md")
	os.WriteFile(promptPath, []byte("You are Lovejoy"), 0644)

	// Create CHANGELOG.md
	os.WriteFile(filepath.Join(tmpDir, "CHANGELOG.md"),
		[]byte("## v1.0.0\n- Initial release"), 0644)

	mock := &mockLLMClient{responses: []string{"Release complete"}}
	lr := &LovejoyRunner{
		BaseRunner: &BaseRunner{
			Agent:      "lovejoy",
			LLMClient:  mock,
			PromptPath: promptPath,
		},
	}

	err := lr.Run(context.Background())
	if err != nil {
		t.Errorf("Run() returned unexpected error: %v", err)
	}

	messages := mock.received[0]
	userMessage := messages[len(messages)-1]
	if !strings.Contains(userMessage.Content, "v1.0.0") {
		t.Errorf("Expected user message to include CHANGELOG.md content")
	}
}
