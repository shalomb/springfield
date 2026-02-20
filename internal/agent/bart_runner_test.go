package agent

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBartRunnerChecksFeedback verifies Bart checks FEEDBACK.md for issues.
func TestBartRunnerChecksFeedback(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create prompt
	promptPath := filepath.Join(tmpDir, "test_prompt.md")
	os.WriteFile(promptPath, []byte("You are Bart"), 0644)

	// Create FEEDBACK.md with issues
	os.WriteFile(filepath.Join(tmpDir, "FEEDBACK.md"),
		[]byte("Status: BLOCKED\nCritical issue found"), 0644)

	mock := &mockLLMClient{responses: []string{"Feedback review complete"}}
	br := &BartRunner{
		BaseRunner: &BaseRunner{
			Agent:      "bart",
			LLMClient:  mock,
			PromptPath: promptPath,
		},
	}

	err := br.Run(context.Background())
	if err != nil {
		t.Errorf("Run() returned unexpected error: %v", err)
	}

	// Verify FEEDBACK.md was checked
	if mock.calls == 0 {
		t.Fatal("Expected LLM Chat to be called")
	}

	messages := mock.received[0]
	userMessage := messages[len(messages)-1]
	if !strings.Contains(userMessage.Content, "BLOCKED") {
		t.Errorf("Expected user message to include feedback content")
	}
}

// TestBartRunnerIncludesFeedbackContent verifies Bart includes FEEDBACK.md in context.
func TestBartRunnerIncludesFeedbackContent(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create prompt
	promptPath := filepath.Join(tmpDir, "test_prompt.md")
	os.WriteFile(promptPath, []byte("You are Bart"), 0644)

	// Create FEEDBACK.md
	feedbackContent := "Test failure in module X"
	os.WriteFile(filepath.Join(tmpDir, "FEEDBACK.md"), []byte(feedbackContent), 0644)

	mock := &mockLLMClient{responses: []string{"Review complete"}}
	br := &BartRunner{
		BaseRunner: &BaseRunner{
			Agent:      "bart",
			LLMClient:  mock,
			PromptPath: promptPath,
		},
	}

	err := br.Run(context.Background())
	if err != nil {
		t.Errorf("Run() returned unexpected error: %v", err)
	}

	messages := mock.received[0]
	userMessage := messages[len(messages)-1]
	if !strings.Contains(userMessage.Content, feedbackContent) {
		t.Errorf("Expected user message to contain FEEDBACK content")
	}
}

// TestBartRunnerWorksWithoutFeedback verifies Bart handles missing FEEDBACK.md.
func TestBartRunnerWorksWithoutFeedback(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create prompt (no FEEDBACK.md)
	promptPath := filepath.Join(tmpDir, "test_prompt.md")
	os.WriteFile(promptPath, []byte("You are Bart"), 0644)

	mock := &mockLLMClient{responses: []string{"Review complete"}}
	br := &BartRunner{
		BaseRunner: &BaseRunner{
			Agent:      "bart",
			LLMClient:  mock,
			PromptPath: promptPath,
			Task:       "Review the code",
		},
	}

	err := br.Run(context.Background())
	if err != nil {
		t.Errorf("Run() returned unexpected error: %v", err)
	}

	if mock.calls == 0 {
		t.Fatal("Expected LLM Chat to be called")
	}
}
