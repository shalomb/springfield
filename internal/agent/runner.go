package agent

import (
	"context"
	"fmt"

	"github.com/shalomb/springfield/internal/config"
	"github.com/shalomb/springfield/internal/llm"
	log "github.com/sirupsen/logrus"
)

// Runner defines the interface for agent runners.
// A runner is responsible for executing an agent's task using an LLM.
type Runner interface {
	// Run executes the agent's task and returns any error.
	Run(ctx context.Context) error
}

// BaseRunner provides common functionality for all agent runners.
// It handles prompt loading and LLM interaction.
type BaseRunner struct {
	// Agent is the name of the agent (e.g., "ralph", "lisa", "bart", "lovejoy").
	Agent string

	// LLMClient is the language model client used for chat interaction.
	LLMClient llm.LLMClient

	// PromptPath is the path to the prompt markdown file.
	// If empty, the path is derived from the Agent name using config.GetPromptPath().
	PromptPath string

	// Task is the user-provided task description.
	Task string

	// Budget is the maximum number of tokens allowed for this run (0 = unlimited).
	Budget int

	// TotalTokensUsed tracks the total tokens consumed by this runner.
	TotalTokensUsed int
}

// SetBudget sets the budget for this runner.
func (br *BaseRunner) SetBudget(budget int) {
	br.Budget = budget
}

// Run executes the agent runner by loading the prompt and calling the LLM.
func (br *BaseRunner) Run(ctx context.Context) error {
	logger := log.WithField("agent", br.Agent)
	
	// Determine the prompt path
	promptPath := br.PromptPath
	if promptPath == "" {
		promptPath = config.GetPromptPath(br.Agent)
	}

	// Load the prompt from the file
	logger.Debugf("Loading prompt from: %s", promptPath)
	prompt, err := config.LoadPrompt(promptPath)
	if err != nil {
		logger.WithError(err).Errorf("Failed to load prompt")
		return fmt.Errorf("failed to load prompt for agent %s: %w", br.Agent, err)
	}
	logger.Debugf("Prompt loaded, %d bytes", len(prompt))

	// Build the initial message with the system prompt
	messages := []llm.Message{
		{
			Role:    "system",
			Content: prompt,
		},
	}

	// Add the task as the user message if provided
	if br.Task != "" {
		messages = append(messages, llm.Message{
			Role:    "user",
			Content: br.Task,
		})
	}

	// Call the LLM
	logger.Infof("🤖 Calling LLM with %d messages...", len(messages))
	response, err := br.LLMClient.Chat(ctx, messages)
	if err != nil {
		logger.WithError(err).Errorf("LLM call failed")
		return fmt.Errorf("LLM call failed: %w", err)
	}
	logger.Infof("✅ LLM responded with %d bytes", len(response.Content))

	// Track token usage
	br.TotalTokensUsed += response.TokenUsage.TotalTokens
	if br.Budget > 0 && br.TotalTokensUsed > br.Budget {
		logger.Warnf("Budget exceeded: %d tokens used of %d allocated", br.TotalTokensUsed, br.Budget)
		return fmt.Errorf("budget exceeded: %d tokens used of %d allocated", br.TotalTokensUsed, br.Budget)
	}

	return nil
}
