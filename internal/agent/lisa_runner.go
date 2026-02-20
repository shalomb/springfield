package agent

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/shalomb/springfield/internal/config"
	"github.com/shalomb/springfield/internal/llm"
)

// LisaRunner implements context aggregation for Lisa.
// It loads FEEDBACK.md and PLAN.md, injects user instructions, and provides
// intelligent pre-processing for task planning.
type LisaRunner struct {
	*BaseRunner

	// ExtraInstruction is user-provided guidance for Lisa.
	ExtraInstruction string
}

// Run executes Lisa's planning cycle with context aggregation.
// It overrides BaseRunner.Run to include additional context.
func (lr *LisaRunner) Run(ctx context.Context) error {
	// Determine the prompt path
	promptPath := lr.PromptPath
	if promptPath == "" {
		promptPath = config.GetPromptPath(lr.Agent)
		lr.PromptPath = promptPath
	}

	// Load the prompt from the file
	prompt, err := config.LoadPrompt(lr.PromptPath)
	if err != nil {
		return fmt.Errorf("failed to load prompt for agent %s: %w", lr.Agent, err)
	}

	// Build the initial message with the system prompt
	messages := []llm.Message{
		{
			Role:    "system",
			Content: prompt,
		},
	}

	// Build the user message with aggregated context
	userMessage := lr.aggregateContext()

	if userMessage != "" {
		messages = append(messages, llm.Message{
			Role:    "user",
			Content: userMessage,
		})
	}

	// Call the LLM using the BaseRunner's LLMClient
	response, err := lr.LLMClient.Chat(ctx, messages)
	if err != nil {
		return fmt.Errorf("LLM call failed: %w", err)
	}

	// Track token usage
	lr.TotalTokensUsed += response.TokenUsage.TotalTokens
	if lr.Budget > 0 && lr.TotalTokensUsed > lr.Budget {
		return fmt.Errorf("budget exceeded: %d tokens used of %d allocated", lr.TotalTokensUsed, lr.Budget)
	}

	return nil
}

// aggregateContext builds the user message by combining task, context files, and instructions.
func (lr *LisaRunner) aggregateContext() string {
	var parts []string

	// Include the original task
	if lr.Task != "" {
		parts = append(parts, lr.Task)
	}

	// Load FEEDBACK.md if it exists
	if feedback, err := readContextFile("FEEDBACK.md"); err == nil && feedback != "" {
		parts = append(parts, fmt.Sprintf("## Current Feedback\n\n%s", feedback))
	}

	// Load PLAN.md if it exists
	if plan, err := readContextFile("PLAN.md"); err == nil && plan != "" {
		parts = append(parts, fmt.Sprintf("## Current Plan\n\n%s", plan))
	}

	// Include extra instructions if provided
	if lr.ExtraInstruction != "" {
		parts = append(parts, fmt.Sprintf("## User Instruction\n\n%s", lr.ExtraInstruction))
	}

	return strings.Join(parts, "\n\n")
}

// readContextFile reads a context file if it exists.
func readContextFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
