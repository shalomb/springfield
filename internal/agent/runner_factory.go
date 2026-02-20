package agent

import (
	"fmt"

	"github.com/shalomb/springfield/internal/llm"
)

// NewRunner creates a specialized runner based on the agent name.
// It returns a Runner interface implementation tailored to the agent's role.
func NewRunner(agentName string, task string, llmClient llm.LLMClient) (Runner, error) {
	return NewRunnerWithBudget(agentName, task, llmClient, 0)
}

// NewRunnerWithBudget creates a specialized runner with a specified budget.
func NewRunnerWithBudget(agentName string, task string, llmClient llm.LLMClient, budget int) (Runner, error) {
	baseRunner := &BaseRunner{
		Agent:     agentName,
		Task:      task,
		LLMClient: llmClient,
		Budget:    budget,
	}

	switch agentName {
	case "ralph":
		return &RalphRunner{
			BaseRunner: baseRunner,
		}, nil

	case "lisa":
		return &LisaRunner{
			BaseRunner: baseRunner,
		}, nil

	case "bart":
		return &BartRunner{
			BaseRunner: baseRunner,
		}, nil

	case "lovejoy":
		return &LovejoyRunner{
			BaseRunner: baseRunner,
		}, nil

	case "marge", "":
		// For marge or default, use BaseRunner (simple single-call agent)
		return baseRunner, nil

	default:
		return nil, fmt.Errorf("unknown agent: %s", agentName)
	}
}
