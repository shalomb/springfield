package agent

import (
	"context"
	"errors"
	"github.com/shalomb/axon/pkg/types"
	"testing"
)

func TestAgent_Run_MaxIterations(t *testing.T) {
	mLLM := &mockLLM{
		responses: []string{"ACTION: ls", "ACTION: ls", "ACTION: ls"},
	}
	mSB := &mockSandbox{
		results: []*types.Result{
			{Stdout: "ok", ExitCode: 0},
			{Stdout: "ok", ExitCode: 0},
			{Stdout: "ok", ExitCode: 0},
		},
	}
	a := New(AgentProfile{Name: "agent", Role: "role"}, mLLM, mSB)
	a.Task = "task"
	a.MaxIterations = 2

	err := a.Run(context.Background())
	if err == nil {
		t.Fatal("expected error on max iterations, got nil")
	}
	if err.Error() != "max iterations reached" {
		t.Errorf("unexpected error message: %v", err)
	}
	if mLLM.calls != 2 {
		t.Errorf("expected 2 LLM calls, got %d", mLLM.calls)
	}
}

func TestAgent_Run_EmptyAction(t *testing.T) {
	mLLM := &mockLLM{
		responses: []string{"ACTION: ", "[[FINISH]]"},
	}
	mSB := &mockSandbox{}
	a := New(AgentProfile{Name: "agent", Role: "role"}, mLLM, mSB)
	a.Task = "task"

	err := a.Run(context.Background())
	if err != nil {
		t.Fatalf("Run() unexpected error: %v", err)
	}

	if mSB.calls != 0 {
		t.Error("expected sandbox not to be called for empty action")
	}
}

func TestAgent_Run_SandboxMaxRetriesReached(t *testing.T) {
	mLLM := &mockLLM{responses: []string{"ACTION: ls"}}
	mSB := &mockSandbox{errors: []error{
		errors.New("e1"), errors.New("e2"), errors.New("e3"), errors.New("e4"),
	}}
	a := New(AgentProfile{Name: "agent", Role: "role"}, mLLM, mSB)
	a.Task = "task"
	a.MaxRetries = 2

	err := a.Run(context.Background())
	if err == nil {
		t.Fatal("expected error after max retries, got nil")
	}
	if mSB.calls != 3 {
		t.Errorf("Sandbox calls = %d, want 3", mSB.calls)
	}
}
