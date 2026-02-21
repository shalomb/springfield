package agent

import (
	"fmt"
	"strings"

	"github.com/shalomb/springfield/internal/config"
	"github.com/shalomb/springfield/internal/llm"
	"github.com/shalomb/springfield/internal/sandbox"
)

// NewRunner creates a specialized runner based on the agent name.
func NewRunner(agentName string, task string, llmClient llm.LLMClient) (Runner, error) {
	return NewRunnerWithBudget(agentName, task, llmClient, nil, 0)
}

// NewRunnerWithBudget creates a specialized runner with a specified budget and optional sandbox.
func NewRunnerWithBudget(agentName string, task string, llmClient llm.LLMClient, sb sandbox.Sandbox, budget int) (Runner, error) {
	normalizedAgent := strings.ToLower(agentName)

	profile, err := GetAgentProfile(normalizedAgent)
	if err != nil {
		return nil, err
	}

	// We'll transition to using the unified Agent as the primary Runner
	// For now, we keep the specialized ones for compatibility if needed,
	// but the goal is to migrate all to the Agent struct.

	a := New(profile, llmClient, sb)
	a.Task = task
	a.Budget = budget

	return a, nil
}

// GetAgentProfile returns the profile for a given agent name.
func GetAgentProfile(agentName string) (AgentProfile, error) {
	roles := map[string]string{
		"marge":   "Product Agent",
		"lisa":    "Planning Agent",
		"ralph":   "Build Agent",
		"bart":    "Quality Agent",
		"lovejoy": "Release Agent",
	}

	role, ok := roles[agentName]
	if !ok {
		return AgentProfile{}, fmt.Errorf("unknown agent: %s", agentName)
	}

	promptPath := config.GetPromptPath(agentName)
	prompt, err := config.LoadPrompt(promptPath)
	if err != nil {
		return AgentProfile{}, fmt.Errorf("failed to load prompt for %s: %w", agentName, err)
	}

	profile := AgentProfile{
		Name:         agentName,
		Role:         role,
		SystemPrompt: prompt,
	}

	// Specialized profile settings
	// NOTE: DO NOT load context files here - they're too large and exceed ARG_MAX
	// instead, instruct agents to use the `read` tool to access PLAN.md, FEEDBACK.md, etc.
	switch agentName {
	case "lisa":
		profile.ContextFiles = []string{} // Use `read` tool instead
		profile.OutputTarget = "PLAN.md"
	case "ralph":
		profile.ContextFiles = []string{} // Use `read` tool instead
		// Ralph handles his own persistence via git/filesystem actions
	case "bart":
		profile.ContextFiles = []string{} // Use `read` tool instead
		profile.OutputTarget = "FEEDBACK.md"
	case "lovejoy":
		profile.ContextFiles = []string{} // Use `read` tool instead
	}

	return profile, nil
}
