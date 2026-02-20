package testutils

import (
	"context"
	"fmt"
	"os"

	"github.com/shalomb/springfield/internal/llm"
)

type MockLLM struct{}

func (m *MockLLM) Chat(ctx context.Context, messages []llm.Message) (llm.Response, error) {
	if os.Getenv("MOCK_LLM_ERROR") == "true" {
		return llm.Response{}, fmt.Errorf("mock llm error")
	}
	// Very simple mock response to allow the loop to finish in tests
	return llm.Response{
		Content: "THOUGHT: I am a mock agent. [[FINISH]]",
		TokenUsage: llm.TokenUsage{
			PromptTokens:     10,
			CompletionTokens: 10,
			TotalTokens:      20,
		},
	}, nil
}
