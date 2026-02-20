package agent

import (
	"context"
	"github.com/shalomb/axon/pkg/types"
	"testing"
)

func TestAgent_Run_PipesAllowed(t *testing.T) {
	mLLM := &mockLLM{
		responses: []string{"ACTION: ls | grep go", "[[FINISH]]"},
	}
	mSB := &mockSandbox{
		results: []*types.Result{{Stdout: "main.go\n", ExitCode: 0}},
	}
	a := New("agent", "role", mLLM, mSB)

	err := a.Run(context.Background(), "list go files")
	if err != nil {
		t.Fatalf("Run() unexpected error: %v", err)
	}

	if mSB.calls != 1 {
		t.Error("expected sandbox to be called for action with pipe")
	}
}
