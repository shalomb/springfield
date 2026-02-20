package agent

import (
	"context"
	"testing"

	"github.com/shalomb/axon/pkg/types"
)

func TestAgent_Run_RobustFinishDetection(t *testing.T) {
	tests := []struct {
		name         string
		llmResponses []string
		wantErr      bool
		wantCalls    int
	}{
		{
			name: "naive FINISH in sentence should not trigger completion",
			llmResponses: []string{
				"I will FINISH later. ACTION: ls",
				"[[FINISH]]",
			},
			wantErr:   false,
			wantCalls: 2,
		},
		{
			name: "[[FINISH]] on its own line should trigger completion",
			llmResponses: []string{
				"I am done.\n[[FINISH]]",
			},
			wantErr:   false,
			wantCalls: 1,
		},
		{
			name: "[[FINISH]] at the end of response should trigger completion",
			llmResponses: []string{
				"Task complete [[FINISH]]",
			},
			wantErr:   false,
			wantCalls: 1,
		},
		{
			name: "FINISH without brackets should not trigger completion anymore",
			llmResponses: []string{
				"FINISH",
				"[[FINISH]]",
			},
			wantErr:   false,
			wantCalls: 2,
		},
		{
			name: "[[FINISH]] followed by more text should not trigger completion",
			llmResponses: []string{
				"[[FINISH]]\nActually, I forgot something.\nACTION: ls",
				"[[FINISH]]",
			},
			wantErr:   false,
			wantCalls: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mLLM := &mockLLM{responses: tt.llmResponses}
			// We need a sandbox mock that returns something for the ACTION in the first test case
			mSB := &mockSandbox{
				results: []*types.Result{{Stdout: "file.txt", ExitCode: 0}},
			}
			a := New("agent", "role", mLLM, mSB)

			err := a.Run(context.Background(), "task")

			if (err != nil) != tt.wantErr {
				t.Fatalf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			if mLLM.calls != tt.wantCalls {
				t.Errorf("LLM calls = %d, want %d", mLLM.calls, tt.wantCalls)
			}
		})
	}
}
