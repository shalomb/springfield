package llm

import (
	"context"
	"testing"
)

func TestPiLLM_Chat(t *testing.T) {
	mockExec := func(ctx context.Context, name string, args ...string) ([]byte, error) {
		if name != "pi" {
			t.Errorf("expected command name 'pi', got %q", name)
		}

		foundSystem := false
		for i, arg := range args {
			if arg == "--system-prompt" {
				if i+1 < len(args) && args[i+1] == "you are a bot" {
					foundSystem = true
				}
			}
		}
		if !foundSystem {
			t.Error("system prompt not found in args")
		}

		foundUser := false
		for _, arg := range args {
			if arg == "hello" {
				foundUser = true
			}
		}
		if !foundUser {
			t.Error("user message not found in args")
		}

		return []byte("response from pi"), nil
	}

	p := &PiLLM{executor: mockExec}
	messages := []Message{
		{Role: "system", Content: "you are a bot"},
		{Role: "user", Content: "hello"},
	}

	resp, err := p.Chat(context.Background(), messages)
	if err != nil {
		t.Fatalf("Chat failed: %v", err)
	}

	if resp.Content != "response from pi" {
		t.Errorf("expected 'response from pi', got %q", resp.Content)
	}
}

func TestPiLLM_Chat_Error(t *testing.T) {
	mockExec := func(ctx context.Context, name string, args ...string) ([]byte, error) {
		return nil, context.DeadlineExceeded
	}

	p := &PiLLM{executor: mockExec}
	_, err := p.Chat(context.Background(), []Message{{Role: "user", Content: "hi"}})
	if err == nil {
		t.Error("expected error, got nil")
	}
}
