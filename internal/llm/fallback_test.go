package llm

import (
	"context"
	"errors"
	"testing"
)

type mockSimpleLLM struct {
	response string
	err      error
	calls    int
}

func (m *mockSimpleLLM) Chat(ctx context.Context, messages []Message) (Response, error) {
	m.calls++
	if m.err != nil {
		return Response{}, m.err
	}
	return Response{Content: m.response}, nil
}

func TestFallbackLLM_Success(t *testing.T) {
	primary := &mockSimpleLLM{response: "primary success"}
	fallback := &mockSimpleLLM{response: "fallback success"}
	f := &FallbackLLM{Primary: primary, Fallback: fallback}

	resp, err := f.Chat(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Content != "primary success" {
		t.Errorf("expected 'primary success', got %q", resp.Content)
	}
	if primary.calls != 1 {
		t.Errorf("expected 1 call to primary, got %d", primary.calls)
	}
	if fallback.calls != 0 {
		t.Errorf("expected 0 calls to fallback, got %d", fallback.calls)
	}
}

func TestFallbackLLM_Fallback(t *testing.T) {
	primary := &mockSimpleLLM{err: errors.New("primary failed")}
	fallback := &mockSimpleLLM{response: "fallback success"}
	f := &FallbackLLM{Primary: primary, Fallback: fallback}

	resp, err := f.Chat(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Content != "fallback success" {
		t.Errorf("expected 'fallback success', got %q", resp.Content)
	}
	if primary.calls != 1 {
		t.Errorf("expected 1 call to primary, got %d", primary.calls)
	}
	if fallback.calls != 1 {
		t.Errorf("expected 1 call to fallback, got %d", fallback.calls)
	}
}
