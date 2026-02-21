package agent

import (
	"os"
	"path/filepath"
	"testing"
)

func setupPromptFiles(t *testing.T, tmpDir string) {
	agents := []string{"ralph", "lisa", "bart", "lovejoy", "marge"}
	for _, a := range agents {
		path := filepath.Join(tmpDir, ".github", "agents")
		_ = os.MkdirAll(path, 0755)
		_ = os.WriteFile(filepath.Join(path, "prompt_"+a+".md"), []byte("You are "+a), 0644)
	}
}

// TestNewRunnerCreatesAgent verifies the factory creates Agent instances.
func TestNewRunnerCreatesAgent(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()
	_ = os.Chdir(tmpDir)
	setupPromptFiles(t, tmpDir)

	mock := &mockLLM{responses: []string{"response"}}
	agents := []string{"ralph", "lisa", "bart", "lovejoy", "marge"}

	for _, name := range agents {
		runner, err := NewRunner(name, "test task", mock)
		if err != nil {
			t.Fatalf("NewRunner(%s) returned error: %v", name, err)
		}

		a, ok := runner.(*Agent)
		if !ok {
			t.Errorf("Expected *Agent for %s, got %T", name, runner)
			continue
		}

		if a.Profile.Name != name {
			t.Errorf("Expected agent name %s, got %s", name, a.Profile.Name)
		}
	}
}

// TestNewRunnerRejectsUnknownAgent verifies factory rejects unknown agents.
func TestNewRunnerRejectsUnknownAgent(t *testing.T) {
	mock := &mockLLM{responses: []string{"response"}}
	_, err := NewRunner("unknown-agent", "test task", mock)
	if err == nil {
		t.Errorf("Expected error for unknown agent, got nil")
	}
}

// TestNewRunnerTaskIsSet verifies the task is passed to runners.
func TestNewRunnerTaskIsSet(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()
	_ = os.Chdir(tmpDir)
	setupPromptFiles(t, tmpDir)

	mock := &mockLLM{responses: []string{"response"}}
	taskText := "important task"
	runner, err := NewRunner("ralph", taskText, mock)
	if err != nil {
		t.Fatalf("NewRunner() returned error: %v", err)
	}

	a := runner.(*Agent)
	if a.Task != taskText {
		t.Errorf("Expected task %q, got %q", taskText, a.Task)
	}
}

// TestNewRunnerWithBudgetSetsBudget verifies budget is properly set.
func TestNewRunnerWithBudgetSetsBudget(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()
	_ = os.Chdir(tmpDir)
	setupPromptFiles(t, tmpDir)

	mock := &mockLLM{responses: []string{"response"}}
	budget := 1000
	runner, err := NewRunnerWithBudget("ralph", "task", mock, nil, budget)
	if err != nil {
		t.Fatalf("NewRunnerWithBudget() returned error: %v", err)
	}

	a := runner.(*Agent)
	if a.Budget != budget {
		t.Errorf("Expected budget %d, got %d", budget, a.Budget)
	}
}
