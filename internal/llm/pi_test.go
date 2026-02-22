package llm

import (
	"context"
	"fmt"
	"strings"
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

		foundMode := false
		for i, arg := range args {
			if arg == "--mode" {
				if i+1 < len(args) && args[i+1] == "json" {
					foundMode = true
				}
			}
		}
		if !foundMode {
			t.Error("--mode json not found in args")
		}

		jsonOutput := `
{"type":"message_start","message":{"role":"assistant","content":[]}}
{"type":"message_update","assistantMessageEvent":{"type":"text_delta","delta":"response from pi"}}
{"type":"turn_end","message":{"role":"assistant","content":[{"type":"text","text":"response from pi"}],"usage":{"input":10,"output":20,"totalTokens":30,"cost":{"total":0.00000375}}}}
`
		return []byte(jsonOutput), nil
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

	if resp.TokenUsage.PromptTokens != 10 {
		t.Errorf("expected PromptTokens 10, got %d", resp.TokenUsage.PromptTokens)
	}
	if resp.TokenUsage.CompletionTokens != 20 {
		t.Errorf("expected CompletionTokens 20, got %d", resp.TokenUsage.CompletionTokens)
	}
	if resp.TokenUsage.TotalTokens != 30 {
		t.Errorf("expected TotalTokens 30, got %d", resp.TokenUsage.TotalTokens)
	}
	if resp.TokenUsage.CostNanoDollars != 3750 {
		t.Errorf("expected CostNanoDollars 3750, got %d", resp.TokenUsage.CostNanoDollars)
	}
}

func TestPiLLM_Chat_JSONError(t *testing.T) {
	mockExec := func(ctx context.Context, name string, args ...string) ([]byte, error) {
		jsonOutput := `
{"type":"agent_start"}
{"type":"turn_start"}
{"type":"turn_end","message":{"role":"assistant","content":[],"errorMessage":"Some API error"},"toolResults":[]}
{"type":"agent_end"}
`
		return []byte(jsonOutput), nil
	}

	p := &PiLLM{executor: mockExec}
	_, err := p.Chat(context.Background(), []Message{{Role: "user", Content: "hi"}})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "pi CLI error: Some API error") {
		t.Errorf("expected error message to contain 'pi CLI error: Some API error', got %v", err)
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

func TestIsQuotaExceeded(t *testing.T) {
	testCases := []struct {
		stderr string
		expect bool
		name   string
	}{
		{"Cloud Code Assist API error (429): You have exhausted your capacity", true, "gemini quota"},
		{"429 Too Many Requests", true, "http 429"},
		{"exhausted your capacity on this model", true, "exhausted capacity"},
		{"rate limit exceeded", true, "rate limit"},
		{"quota_exceeded", true, "quota keyword"},
		{"billing_exception", true, "billing error"},
		{"401 Unauthorized", true, "unauthorized"},
		{"403 Forbidden", true, "forbidden"},
		{"random error message", false, "no quota marker"},
		{"connection timeout", false, "timeout"},
		{"", false, "empty string"},
		{`Error: 429 {"type":"error","error":{"type":"rate_limit_error","message":"This request would exceed your account's rate limit. Please try again later."},"request_id":"req_011CYL3UE3hxCyxpV9ELnHVR"}`, true, "anthropic rate limit error"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isQuotaExceeded(tc.stderr)
			if result != tc.expect {
				t.Errorf("isQuotaExceeded(%q) = %v, expected %v", tc.stderr, result, tc.expect)
			}
		})
	}
}

func TestExtractAnthropicErrorMessage(t *testing.T) {
	testCases := []struct {
		stderr   string
		expected string
		name     string
	}{
		{
			name:     "anthropic rate limit error",
			stderr:   `Error: 429 {"type":"error","error":{"type":"rate_limit_error","message":"This request would exceed your account's rate limit. Please try again later."},"request_id":"req_011CYL3UE3hxCyxpV9ELnHVR"}`,
			expected: "Anthropic API error (rate_limit_error): This request would exceed your account's rate limit. Please try again later.",
		},
		{
			name:     "anthropic invalid request error",
			stderr:   `Error: 400 {"type":"error","error":{"type":"invalid_request_error","message":"Invalid request body"},"request_id":"req_123"}`,
			expected: "Anthropic API error (invalid_request_error): Invalid request body",
		},
		{
			name:     "non-json error",
			stderr:   "Error: 429 plain text error",
			expected: "",
		},
		{
			name:     "empty string",
			stderr:   "",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := extractAnthropicErrorMessage(tc.stderr)
			if result != tc.expected {
				t.Errorf("extractAnthropicErrorMessage got %q, expected %q", result, tc.expected)
			}
		})
	}
}

func TestIsQuotaExceededError(t *testing.T) {
	quotaErr := &QuotaExceededError{Message: "quota exceeded"}
	if !IsQuotaExceededError(quotaErr) {
		t.Error("IsQuotaExceededError failed to detect QuotaExceededError")
	}

	genericErr := fmt.Errorf("some error")
	if IsQuotaExceededError(genericErr) {
		t.Error("IsQuotaExceededError incorrectly detected generic error as quota error")
	}
}
