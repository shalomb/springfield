package agent

import (
	"context"
	"testing"
)

func TestAgent_Run_NoActionOrFinish(t *testing.T) {
	mLLM := &mockLLM{
		responses: []string{"Just talking...", "[[FINISH]]"},
	}
	mSB := &mockSandbox{}
	a := New(AgentProfile{Name: "agent", Role: "role"}, mLLM, mSB)
	a.Task = "talk"

	err := a.Run(context.Background())
	if err != nil {
		t.Fatalf("Run() unexpected error: %v", err)
	}

	if mLLM.calls != 2 {
		t.Errorf("expected 2 LLM calls, got %d", mLLM.calls)
	}
}
