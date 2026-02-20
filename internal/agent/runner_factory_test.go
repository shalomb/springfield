package agent

import (
	"context"
	"testing"
)

// TestNewRunnerCreatesRalphRunner verifies the factory creates RalphRunner.
func TestNewRunnerCreatesRalphRunner(t *testing.T) {
	mock := &mockLLMClient{responses: []string{"response"}}
	runner, err := NewRunner("ralph", "test task", mock)
	if err != nil {
		t.Fatalf("NewRunner() returned error: %v", err)
	}

	if _, ok := runner.(*RalphRunner); !ok {
		t.Errorf("Expected *RalphRunner, got %T", runner)
	}
}

// TestNewRunnerCreatesLisaRunner verifies the factory creates LisaRunner.
func TestNewRunnerCreatesLisaRunner(t *testing.T) {
	mock := &mockLLMClient{responses: []string{"response"}}
	runner, err := NewRunner("lisa", "test task", mock)
	if err != nil {
		t.Fatalf("NewRunner() returned error: %v", err)
	}

	if _, ok := runner.(*LisaRunner); !ok {
		t.Errorf("Expected *LisaRunner, got %T", runner)
	}
}

// TestNewRunnerCreatesBartRunner verifies the factory creates BartRunner.
func TestNewRunnerCreatesBartRunner(t *testing.T) {
	mock := &mockLLMClient{responses: []string{"response"}}
	runner, err := NewRunner("bart", "test task", mock)
	if err != nil {
		t.Fatalf("NewRunner() returned error: %v", err)
	}

	if _, ok := runner.(*BartRunner); !ok {
		t.Errorf("Expected *BartRunner, got %T", runner)
	}
}

// TestNewRunnerCreatesLovejoyRunner verifies the factory creates LovejoyRunner.
func TestNewRunnerCreatesLovejoyRunner(t *testing.T) {
	mock := &mockLLMClient{responses: []string{"response"}}
	runner, err := NewRunner("lovejoy", "test task", mock)
	if err != nil {
		t.Fatalf("NewRunner() returned error: %v", err)
	}

	if _, ok := runner.(*LovejoyRunner); !ok {
		t.Errorf("Expected *LovejoyRunner, got %T", runner)
	}
}

// TestNewRunnerCreatesBaseRunnerForMarge verifies Marge uses BaseRunner.
func TestNewRunnerCreatesBaseRunnerForMarge(t *testing.T) {
	mock := &mockLLMClient{responses: []string{"response"}}
	runner, err := NewRunner("marge", "test task", mock)
	if err != nil {
		t.Fatalf("NewRunner() returned error: %v", err)
	}

	if _, ok := runner.(*BaseRunner); !ok {
		t.Errorf("Expected *BaseRunner, got %T", runner)
	}

	// Verify it's not a specialized runner
	if _, ok := runner.(*RalphRunner); ok {
		t.Errorf("Expected BaseRunner, not RalphRunner")
	}
}

// TestNewRunnerRejectsUnknownAgent verifies factory rejects unknown agents.
func TestNewRunnerRejectsUnknownAgent(t *testing.T) {
	mock := &mockLLMClient{responses: []string{"response"}}
	_, err := NewRunner("unknown-agent", "test task", mock)
	if err == nil {
		t.Errorf("Expected error for unknown agent, got nil")
	}
}

// TestNewRunnerTaskIsSet verifies the task is passed to runners.
func TestNewRunnerTaskIsSet(t *testing.T) {
	mock := &mockLLMClient{responses: []string{"response"}}
	taskText := "important task"
	runner, err := NewRunner("ralph", taskText, mock)
	if err != nil {
		t.Fatalf("NewRunner() returned error: %v", err)
	}

	br := runner.(*RalphRunner).BaseRunner
	if br.Task != taskText {
		t.Errorf("Expected task %q, got %q", taskText, br.Task)
	}
}

// TestNewRunnerLLMClientIsSet verifies the LLM client is set.
func TestNewRunnerLLMClientIsSet(t *testing.T) {
	mock := &mockLLMClient{responses: []string{"response"}}
	runner, err := NewRunner("ralph", "task", mock)
	if err != nil {
		t.Fatalf("NewRunner() returned error: %v", err)
	}

	br := runner.(*RalphRunner).BaseRunner
	if br.LLMClient != mock {
		t.Errorf("Expected LLM client to be set")
	}
}

// TestNewRunnerImplementsRunner verifies all created runners satisfy Runner.
func TestNewRunnerImplementsRunner(t *testing.T) {
	tests := []string{"ralph", "lisa", "bart", "lovejoy", "marge"}
	mock := &mockLLMClient{responses: []string{"response"}}

	for _, agentName := range tests {
		t.Run(agentName, func(t *testing.T) {
			runner, err := NewRunner(agentName, "task", mock)
			if err != nil {
				t.Fatalf("NewRunner() returned error: %v", err)
			}

			// Verify it can be called as a Runner
			ctx := context.Background()
			_ = runner.Run(ctx)
			// We don't check error here because mocks may fail; we just
			// verify the method exists and is callable.
		})
	}
}

// TestNewRunnerWithBudgetSetsBudget verifies budget is properly set.
func TestNewRunnerWithBudgetSetsBudget(t *testing.T) {
	mock := &mockLLMClient{responses: []string{"response"}}
	budget := 1000
	runner, err := NewRunnerWithBudget("ralph", "task", mock, budget)
	if err != nil {
		t.Fatalf("NewRunnerWithBudget() returned error: %v", err)
	}

	br := runner.(*RalphRunner).BaseRunner
	if br.Budget != budget {
		t.Errorf("Expected budget %d, got %d", budget, br.Budget)
	}
}

// TestNewRunnerWithBudgetZeroBudgetAllowed verifies zero budget is allowed.
func TestNewRunnerWithBudgetZeroBudgetAllowed(t *testing.T) {
	mock := &mockLLMClient{responses: []string{"response"}}
	runner, err := NewRunnerWithBudget("lisa", "task", mock, 0)
	if err != nil {
		t.Fatalf("NewRunnerWithBudget(0) returned error: %v", err)
	}

	br := runner.(*LisaRunner).BaseRunner
	if br.Budget != 0 {
		t.Errorf("Expected budget 0, got %d", br.Budget)
	}
}
