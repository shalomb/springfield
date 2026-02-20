package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/shalomb/springfield/internal/config"
	"github.com/shalomb/springfield/internal/llm"
)

// BartRunner implements quality verification and adversarial testing.
// It reviews code for issues, runs tests, and provides feedback.
type BartRunner struct {
	*BaseRunner
}

// Run executes Bart's quality review cycle.
func (br *BartRunner) Run(ctx context.Context) error {
	// Determine the prompt path
	promptPath := br.PromptPath
	if promptPath == "" {
		promptPath = config.GetPromptPath(br.Agent)
		br.PromptPath = promptPath
	}

	// Load the prompt from the file
	prompt, err := config.LoadPrompt(promptPath)
	if err != nil {
		return fmt.Errorf("failed to load prompt for agent %s: %w", br.Agent, err)
	}

	// Build the initial message with the system prompt
	messages := []llm.Message{
		{
			Role:    "system",
			Content: prompt,
		},
	}

	// Build the user message with feedback context
	userMessage := br.aggregateFeedbackContext()

	if userMessage != "" {
		messages = append(messages, llm.Message{
			Role:    "user",
			Content: userMessage,
		})
	} else if br.Task != "" {
		// If no feedback context, use the task description
		messages = append(messages, llm.Message{
			Role:    "user",
			Content: br.Task,
		})
	}

	// Call the LLM
	response, err := br.LLMClient.Chat(ctx, messages)
	if err != nil {
		return fmt.Errorf("LLM call failed: %w", err)
	}

	// Track token usage
	br.TotalTokensUsed += response.TokenUsage.TotalTokens
	if br.Budget > 0 && br.TotalTokensUsed > br.Budget {
		return fmt.Errorf("budget exceeded: %d tokens used of %d allocated", br.TotalTokensUsed, br.Budget)
	}

	return nil
}

// aggregateFeedbackContext builds the user message with feedback context.
func (br *BartRunner) aggregateFeedbackContext() string {
	var parts []string

	// Load FEEDBACK.md if it exists
	if feedback, err := readContextFile("FEEDBACK.md"); err == nil && feedback != "" {
		parts = append(parts, fmt.Sprintf("## Current Feedback\n\n%s", feedback))
	}

	// Include the task if provided
	if br.Task != "" {
		parts = append(parts, br.Task)
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, "\n\n")
}
