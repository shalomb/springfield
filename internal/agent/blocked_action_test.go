package agent

import (
	"context"
	"testing"

	"github.com/shalomb/axon/pkg/types"
)

func TestAgent_Run_BlockedAction(t *testing.T) {
	mLLM := &mockLLM{
		responses: []string{"ACTION: rm -rf / ; echo hi", "[[FINISH]]"},
	}
	mSB := &mockSandbox{}
	a := New("agent", "role", mLLM, mSB)

	err := a.Run(context.Background(), "do something bad")
	if err != nil {
		t.Fatalf("Run() unexpected error: %v", err)
	}

	if mSB.calls != 0 {
		t.Error("expected sandbox not to be called for blocked action")
	}

	if mLLM.calls != 2 {
		t.Errorf("expected 2 LLM calls, got %d", mLLM.calls)
	}

	// Check that the blocked message was sent back to LLM
	lastMsg := mLLM.received[1][len(mLLM.received[1])-1]
	if lastMsg.Content != "Action blocked for security reasons." {
		t.Errorf("unexpected message sent to LLM: %q", lastMsg.Content)
	}
}

func TestAgent_Run_AllowedAction_Redirection(t *testing.T) {
	mLLM := &mockLLM{
		responses: []string{"ACTION: echo hello > out.txt", "[[FINISH]]"},
	}
	mSB := &mockSandbox{
		results: []*types.Result{{Stdout: "", ExitCode: 0}},
	}
	a := New("agent", "role", mLLM, mSB)

	err := a.Run(context.Background(), "write to file")
	if err != nil {
		t.Fatalf("Run() unexpected error: %v", err)
	}

	if mSB.calls != 1 {
		t.Errorf("expected sandbox to be called for redirection, got %d", mSB.calls)
	}
	if mSB.commands[0] != "echo hello > out.txt" {
		t.Errorf("unexpected command executed: %q", mSB.commands[0])
	}
}
