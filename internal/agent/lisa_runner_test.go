package agent

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLisaRunnerLoadsContextFiles verifies Lisa loads FEEDBACK.md and PLAN.md.
func TestLisaRunnerLoadsContextFiles(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create prompt
	promptPath := filepath.Join(tmpDir, "test_prompt.md")
	os.WriteFile(promptPath, []byte("You are Lisa"), 0644)

	// Create FEEDBACK.md
	os.WriteFile(filepath.Join(tmpDir, "FEEDBACK.md"), []byte("Test feedback"), 0644)

	// Create PLAN.md
	os.WriteFile(filepath.Join(tmpDir, "PLAN.md"), []byte("Test plan"), 0644)

	mock := &mockLLMClient{responses: []string{"Analysis complete"}}
	lr := &LisaRunner{
		BaseRunner: &BaseRunner{
			Agent:      "lisa",
			LLMClient:  mock,
			PromptPath: promptPath,
		},
	}

	err := lr.Run(context.Background())
	if err != nil {
		t.Errorf("Run() returned unexpected error: %v", err)
	}

	// Verify LLM was called
	if mock.calls == 0 {
		t.Fatal("Expected LLM Chat to be called")
	}

	// Verify context files were included in the message
	messages := mock.received[0]
	if len(messages) < 2 {
		t.Fatal("Expected at least system and user messages")
	}

	userMessage := messages[len(messages)-1]
	if !strings.Contains(userMessage.Content, "Test feedback") {
		t.Errorf("Expected user message to contain FEEDBACK.md content")
	}
	if !strings.Contains(userMessage.Content, "Test plan") {
		t.Errorf("Expected user message to contain PLAN.md content")
	}
}

// TestLisaRunnerOptionalContextFiles verifies Lisa works without context files.
func TestLisaRunnerOptionalContextFiles(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create prompt (no FEEDBACK.md or PLAN.md)
	promptPath := filepath.Join(tmpDir, "test_prompt.md")
	os.WriteFile(promptPath, []byte("You are Lisa"), 0644)

	mock := &mockLLMClient{responses: []string{"Analysis complete"}}
	lr := &LisaRunner{
		BaseRunner: &BaseRunner{
			Agent:      "lisa",
			LLMClient:  mock,
			PromptPath: promptPath,
			Task:       "Analyze the project",
		},
	}

	err := lr.Run(context.Background())
	if err != nil {
		t.Errorf("Run() returned unexpected error: %v", err)
	}

	// Verify LLM was called
	if mock.calls == 0 {
		t.Fatal("Expected LLM Chat to be called")
	}
}

// TestLisaRunnerIncludesExtraInstruction verifies Lisa includes user instructions.
func TestLisaRunnerIncludesExtraInstruction(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create prompt
	promptPath := filepath.Join(tmpDir, "test_prompt.md")
	os.WriteFile(promptPath, []byte("You are Lisa"), 0644)

	mock := &mockLLMClient{responses: []string{"Analysis complete"}}
	lr := &LisaRunner{
		BaseRunner: &BaseRunner{
			Agent:      "lisa",
			LLMClient:  mock,
			PromptPath: promptPath,
			Task:       "Original task",
		},
		ExtraInstruction: "Please focus on module design",
	}

	err := lr.Run(context.Background())
	if err != nil {
		t.Errorf("Run() returned unexpected error: %v", err)
	}

	// Verify extra instruction was included
	messages := mock.received[0]
	userMessage := messages[len(messages)-1]
	if !strings.Contains(userMessage.Content, "Please focus on module design") {
		t.Errorf("Expected user message to contain extra instruction")
	}
}
