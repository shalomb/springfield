package agent

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/shalomb/springfield/internal/llm"
)

// mockLLMClient is a simple test double for testing runners.
type mockLLMClient struct {
	responses []string
	calls     int
	received  [][]llm.Message
}

func (m *mockLLMClient) Chat(ctx context.Context, messages []llm.Message) (llm.Response, error) {
	if m.calls >= len(m.responses) {
		m.calls++
		return llm.Response{}, nil
	}
	cp := make([]llm.Message, len(messages))
	copy(cp, messages)
	m.received = append(m.received, cp)
	response := m.responses[m.calls]
	m.calls++
	return llm.Response{
		Content: response,
		TokenUsage: llm.TokenUsage{
			PromptTokens:     10,
			CompletionTokens: 10,
			TotalTokens:      20,
		},
	}, nil
}

// TestBaseRunnerHasRunMethod verifies BaseRunner implements Run() method.
func TestBaseRunnerHasRunMethod(t *testing.T) {
	// Create a temporary prompt file
	tmpDir := t.TempDir()
	promptPath := filepath.Join(tmpDir, "test_prompt.md")
	promptContent := "You are a test runner"
	err := os.WriteFile(promptPath, []byte(promptContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test prompt file: %v", err)
	}

	mock := &mockLLMClient{responses: []string{"Test response"}}
	br := &BaseRunner{
		Agent:      "test",
		LLMClient:  mock,
		PromptPath: promptPath,
	}

	// Should not panic or error
	err = br.Run(context.Background())
	if err != nil {
		t.Errorf("Run() returned unexpected error: %v", err)
	}
}

// TestBaseRunnerLoadsPrompt verifies BaseRunner loads and uses prompts.
func TestBaseRunnerLoadsPrompt(t *testing.T) {
	// Create a temporary prompt file
	tmpDir := t.TempDir()
	promptPath := filepath.Join(tmpDir, "test_prompt.md")
	promptContent := "You are a test agent"
	err := os.WriteFile(promptPath, []byte(promptContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test prompt file: %v", err)
	}

	mock := &mockLLMClient{responses: []string{"Test response"}}
	br := &BaseRunner{
		Agent:      "test",
		LLMClient:  mock,
		PromptPath: promptPath,
	}

	err = br.Run(context.Background())
	if err != nil {
		t.Errorf("Run() returned unexpected error: %v", err)
	}

	// Verify that the prompt was included in the LLM call
	if len(mock.received) == 0 {
		t.Fatal("Expected LLM Chat to be called with messages")
	}

	messages := mock.received[0]
	if len(messages) == 0 {
		t.Fatal("Expected at least one message in LLM call")
	}

	// First message should be system message with the prompt
	if messages[0].Role != "system" {
		t.Errorf("Expected first message to be system role, got %s", messages[0].Role)
	}

	if messages[0].Content != promptContent {
		t.Errorf("Expected system message to contain prompt: %q, got %q", promptContent, messages[0].Content)
	}
}

// TestBaseRunnerInitialization verifies BaseRunner can be properly initialized.
func TestBaseRunnerInitialization(t *testing.T) {
	mock := &mockLLMClient{responses: []string{"Test response"}}
	br := &BaseRunner{
		Agent:     "ralph",
		LLMClient: mock,
	}

	if br.Agent != "ralph" {
		t.Errorf("Expected agent name 'ralph', got %s", br.Agent)
	}

	if br.LLMClient != mock {
		t.Errorf("Expected LLM client to be set, got nil")
	}
}

// TestRunnerInterfaceIsImplemented verifies BaseRunner satisfies Runner interface.
func TestRunnerInterfaceIsImplemented(t *testing.T) {
	var _ Runner = (*BaseRunner)(nil)
}
