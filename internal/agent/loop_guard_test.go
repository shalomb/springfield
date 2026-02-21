package agent

import (
	"context"
	"testing"

	"github.com/shalomb/axon/pkg/types"
)

func TestAgent_Run_InfiniteLoopGuard(t *testing.T) {
	// Setup an LLM that always returns an action and never FINISH
	mLLM := &mockLLM{
		responses: make([]string, 100),
	}
	for i := range mLLM.responses {
		mLLM.responses[i] = "ACTION: echo infinite"
	}

	// Mock sandbox that always succeeds
	mSB := &mockSandbox{
		results: make([]*types.Result, 100),
	}
	for i := range mSB.results {
		mSB.results[i] = &types.Result{ExitCode: 0}
	}

	a := New(AgentProfile{Name: "agent", Role: "role"}, mLLM, mSB)
	a.Task = "loop forever"
	a.MaxIterations = 5 // Set a small limit for testing

	err := a.Run(context.Background())

	if err == nil {
		t.Fatal("expected error from infinite loop guard, got nil")
	}

	if err.Error() != "max iterations reached" {
		t.Errorf("unexpected error: %v", err)
	}

	if mLLM.calls != 5 {
		t.Errorf("expected 5 LLM calls, got %d", mLLM.calls)
	}
}
