package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/shalomb/springfield/internal/config"
	"github.com/shalomb/springfield/internal/llm"
)

// LovejoyRunner implements release ceremony logic.
// It verifies readiness, captures learnings, and coordinates the merge.
type LovejoyRunner struct {
	*BaseRunner
}

// Run executes Lovejoy's release ceremony.
func (lr *LovejoyRunner) Run(ctx context.Context) error {
	// Determine the prompt path
	promptPath := lr.PromptPath
	if promptPath == "" {
		promptPath = config.GetPromptPath(lr.Agent)
		lr.PromptPath = promptPath
	}

	// Load the prompt from the file
	prompt, err := config.LoadPrompt(promptPath)
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

	// Build the user message with release context
	userMessage := lr.aggregateReleaseContext()

	if userMessage != "" {
		messages = append(messages, llm.Message{
			Role:    "user",
			Content: userMessage,
		})
	} else if lr.Task != "" {
		// If no release context, use the task description
		messages = append(messages, llm.Message{
			Role:    "user",
			Content: lr.Task,
		})
	}

	// Call the LLM
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

// aggregateReleaseContext builds the user message with release readiness information.
func (lr *LovejoyRunner) aggregateReleaseContext() string {
	var parts []string

	// Check if TODO.md exists (blocking issue)
	if todoExists, _ := RalphTODOExists(); todoExists {
		parts = append(parts, "⚠️  WARNING: TODO.md still exists. Release may be blocked.")
	} else {
		parts = append(parts, "✅ TODO.md is empty. Work is complete.")
	}

	// Load CHANGELOG.md if it exists
	if changelog, err := readContextFile("CHANGELOG.md"); err == nil && changelog != "" {
		parts = append(parts, fmt.Sprintf("## Changelog\n\n%s", changelog))
	}

	// Include the task if provided
	if lr.Task != "" {
		parts = append(parts, lr.Task)
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, "\n\n")
}
