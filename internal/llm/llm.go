package llm

import "context"

// TokenUsage represents the token counts for an LLM response.
type TokenUsage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// Response represents a full response from an LLM.
type Response struct {
	Content    string
	TokenUsage TokenUsage
}

// LLMClient defines the interface for interacting with a Large Language Model.
type LLMClient interface {
	Chat(ctx context.Context, messages []Message) (Response, error)
}

// FallbackLLM wraps a primary and fallback LLM client.
type FallbackLLM struct {
	Primary  LLMClient
	Fallback LLMClient
}

func (f *FallbackLLM) Chat(ctx context.Context, messages []Message) (Response, error) {
	resp, err := f.Primary.Chat(ctx, messages)
	if err == nil {
		return resp, nil
	}
	if f.Fallback == nil {
		return resp, err
	}
	return f.Fallback.Chat(ctx, messages)
}

// Message represents a single message in a chat conversation.
type Message struct {
	Role    string
	Content string
}
