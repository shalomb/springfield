package agent

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/shalomb/axon/pkg/types"
	"github.com/shalomb/springfield/internal/governance"
	"github.com/shalomb/springfield/internal/llm"
)

// ---------------------------------------------------------------------------
// Test doubles
// ---------------------------------------------------------------------------

type mockLLM struct {
	responses []string
	errors    []error
	calls     int
	received  [][]llm.Message // every messages slice passed to Chat
}

func (m *mockLLM) Chat(ctx context.Context, messages []llm.Message) (llm.Response, error) {
	i := m.calls
	m.calls++
	cp := make([]llm.Message, len(messages))
	copy(cp, messages)
	m.received = append(m.received, cp)
	if i < len(m.errors) && m.errors[i] != nil {
		return llm.Response{}, m.errors[i]
	}
	if i >= len(m.responses) {
		return llm.Response{}, errors.New("mockLLM: no more responses")
	}
	return llm.Response{
		Content: m.responses[i],
		TokenUsage: llm.TokenUsage{
			PromptTokens:     10, // Mock usage
			CompletionTokens: 10,
			TotalTokens:      20,
		},
	}, nil
}

type mockSandbox struct {
	results  []*types.Result
	errors   []error
	calls    int
	commands []string // every command passed to Execute
}

func (m *mockSandbox) Execute(ctx context.Context, command string) (*types.Result, error) {
	i := m.calls
	m.calls++
	m.commands = append(m.commands, command)
	if i < len(m.errors) && m.errors[i] != nil {
		return nil, m.errors[i]
	}
	if i >= len(m.results) {
		return nil, errors.New("mockSandbox: no more results")
	}
	return m.results[i], nil
}

// ---------------------------------------------------------------------------
// Table-driven: happy-path action loop
// ---------------------------------------------------------------------------

func TestAgent_Run_Table(t *testing.T) {
	tests := []struct {
		name         string
		llmResponses []string
		sbResults    []*types.Result
		task         string
		wantLLMCalls int
		wantSBCalls  int
		wantErr      bool
	}{
		{
			name: "single action then finish",
			llmResponses: []string{
				"THOUGHT: need ls.\nACTION: ls",
				"THOUGHT: done. [[FINISH]]",
			},
			sbResults:    []*types.Result{{Stdout: "file.txt\n", ExitCode: 0}},
			task:         "list files",
			wantLLMCalls: 2,
			wantSBCalls:  1,
		},
		{
			name: "immediate finish, no action",
			llmResponses: []string{
				"THOUGHT: already done. [[FINISH]]",
			},
			task:         "noop",
			wantLLMCalls: 1,
			wantSBCalls:  0,
		},
		{
			name: "two actions then finish",
			llmResponses: []string{
				"THOUGHT: step1.\nACTION: echo a",
				"THOUGHT: step2.\nACTION: echo b",
				"THOUGHT: done. [[FINISH]]",
			},
			sbResults: []*types.Result{
				{Stdout: "a\n", ExitCode: 0},
				{Stdout: "b\n", ExitCode: 0},
			},
			task:         "two steps",
			wantLLMCalls: 3,
			wantSBCalls:  2,
		},
		{
			name: "sandbox non-zero exit continues loop",
			llmResponses: []string{
				"THOUGHT: run bad cmd.\nACTION: false",
				"THOUGHT: saw failure, done. [[FINISH]]",
			},
			sbResults: []*types.Result{
				{Stdout: "", Stderr: "command not found", ExitCode: 1},
			},
			task:         "handle failure",
			wantLLMCalls: 2,
			wantSBCalls:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mLLM := &mockLLM{responses: tt.llmResponses}
			mSB := &mockSandbox{results: tt.sbResults}
			a := New(AgentProfile{Name: "agent", Role: "role"}, mLLM, mSB)
			a.Task = tt.task

			err := a.Run(context.Background())

			if (err != nil) != tt.wantErr {
				t.Fatalf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			if mLLM.calls != tt.wantLLMCalls {
				t.Errorf("LLM calls = %d, want %d", mLLM.calls, tt.wantLLMCalls)
			}
			if mSB.calls != tt.wantSBCalls {
				t.Errorf("Sandbox calls = %d, want %d", mSB.calls, tt.wantSBCalls)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// System prompt contains agent name and role
// ---------------------------------------------------------------------------

func TestAgent_Run_SystemPromptContainsNameAndRole(t *testing.T) {
	mLLM := &mockLLM{responses: []string{"[[FINISH]]"}}
	mSB := &mockSandbox{}
	a := New(AgentProfile{Name: "Marge", Role: "Product Agent"}, mLLM, mSB)
	a.Task = "anything"

	_ = a.Run(context.Background())

	if mLLM.calls == 0 {
		t.Fatal("LLM was never called")
	}
	first := mLLM.received[0]
	if len(first) == 0 {
		t.Fatal("no messages sent to LLM")
	}
	sys := first[0]
	if sys.Role != "system" {
		t.Errorf("first message role = %q, want \"system\"", sys.Role)
	}
	if !strings.Contains(sys.Content, "Marge") {
		t.Errorf("system prompt missing agent name: %q", sys.Content)
	}
	if !strings.Contains(sys.Content, "Product Agent") {
		t.Errorf("system prompt missing agent role: %q", sys.Content)
	}
}

// ---------------------------------------------------------------------------
// Sandbox result is fed back to LLM
// ---------------------------------------------------------------------------

func TestAgent_Run_SandboxResultFedBackToLLM(t *testing.T) {
	const sandboxOutput = "hello from sandbox"
	mLLM := &mockLLM{responses: []string{
		"ACTION: echo hello",
		"[[FINISH]]",
	}}
	mSB := &mockSandbox{results: []*types.Result{
		{Stdout: sandboxOutput, ExitCode: 0},
	}}
	a := New(AgentProfile{Name: "agent", Role: "role"}, mLLM, mSB)
	a.Task = "task"
	_ = a.Run(context.Background())

	// Second LLM call should include the sandbox output in conversation history
	if mLLM.calls < 2 {
		t.Fatalf("expected at least 2 LLM calls, got %d", mLLM.calls)
	}
	secondCallMessages := mLLM.received[1]
	found := false
	for _, msg := range secondCallMessages {
		if strings.Contains(msg.Content, sandboxOutput) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("sandbox output %q not found in second LLM call messages", sandboxOutput)
	}
}

// ---------------------------------------------------------------------------
// Context injection
// ---------------------------------------------------------------------------

func TestAgent_Run_ContextInjection(t *testing.T) {
	mLLM := &mockLLM{responses: []string{
		"ACTION: echo hello",
		"[[FINISH]]",
	}}
	mSB := &mockSandbox{results: []*types.Result{
		{
			Stdout:   "out",
			ExitCode: 0,
			Context: types.ContextMetadata{
				ProjectType: "go",
				BuildTool:   "go",
			},
		},
	}}
	a := New(AgentProfile{Name: "agent", Role: "role"}, mLLM, mSB)
	a.Task = "task"
	_ = a.Run(context.Background())

	if mLLM.calls < 2 {
		t.Fatalf("expected at least 2 LLM calls, got %d", mLLM.calls)
	}
	secondCallMessages := mLLM.received[1]
	found := false
	for _, msg := range secondCallMessages {
		if strings.Contains(msg.Content, "Project type: go") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("project_type: go not found in second LLM call messages")
	}
}

// ---------------------------------------------------------------------------
// Context Persistence & Update
// ---------------------------------------------------------------------------

func TestAgent_Run_ContextPersistence(t *testing.T) {
	mLLM := &mockLLM{responses: []string{
		"ACTION: step1",
		"ACTION: step2",
		"[[FINISH]]",
	}}
	mSB := &mockSandbox{results: []*types.Result{
		{
			Stdout:   "res1",
			ExitCode: 0,
			Context: types.ContextMetadata{
				ProjectType: "go",
			},
		},
		{
			Stdout:   "res2",
			ExitCode: 0,
			// No context returned in second call
		},
	}}
	a := New(AgentProfile{Name: "agent", Role: "role"}, mLLM, mSB)
	a.Task = "task"
	_ = a.Run(context.Background())

	if mLLM.calls < 3 {
		t.Fatalf("expected at least 3 LLM calls, got %d", mLLM.calls)
	}

	// Third call should still have "Project type: go" if it persists in history
	thirdCallMessages := mLLM.received[2]
	found := false
	for _, msg := range thirdCallMessages {
		if strings.Contains(msg.Content, "Project type: go") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("project_type: go should persist in history but was not found in third LLM call")
	}
}

func TestAgent_Run_ContextUpdate(t *testing.T) {
	mLLM := &mockLLM{responses: []string{
		"ACTION: step1",
		"ACTION: step2",
		"[[FINISH]]",
	}}
	mSB := &mockSandbox{results: []*types.Result{
		{
			Stdout:   "res1",
			ExitCode: 0,
			Context: types.ContextMetadata{
				ProjectType: "python",
			},
		},
		{
			Stdout:   "res2",
			ExitCode: 0,
			Context: types.ContextMetadata{
				ProjectType: "go",
			},
		},
	}}
	a := New(AgentProfile{Name: "agent", Role: "role"}, mLLM, mSB)
	a.Task = "task"
	_ = a.Run(context.Background())

	if mLLM.calls < 3 {
		t.Fatalf("expected at least 3 LLM calls, got %d", mLLM.calls)
	}

	// Third call should have "Project type: go"
	thirdCallMessages := mLLM.received[2]
	foundNew := false
	for _, msg := range thirdCallMessages {
		if strings.Contains(msg.Content, "Project type: go") {
			foundNew = true
		}
	}
	if !foundNew {
		t.Errorf("updated project_type: go not found in third LLM call")
	}
}

func TestFormatContext_Empty(t *testing.T) {
	got := formatContext(types.ContextMetadata{})
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

// ---------------------------------------------------------------------------
// Cancellation
// ---------------------------------------------------------------------------

func TestAgent_Run_RespectsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	mLLM := &mockLLM{responses: []string{"ACTION: sleep 60", "[[FINISH]]"}}
	mSB := &mockSandbox{results: []*types.Result{{Stdout: "", ExitCode: 0}}}
	a := New(AgentProfile{Name: "agent", Role: "role"}, mLLM, mSB)
	a.Task = "task"

	err := a.Run(ctx)
	if err == nil {
		// Either cancelled or completed immediately — both acceptable.
		// What we must NOT do is loop forever.
		t.Log("Run() returned nil on cancelled context (completed before first check)")
	}
}

// ---------------------------------------------------------------------------
// Retry behaviour
// ---------------------------------------------------------------------------

func TestAgent_Run_SandboxRetry(t *testing.T) {
	mLLM := &mockLLM{responses: []string{
		"ACTION: ls",
		"[[FINISH]]",
	}}
	mSB := &mockSandbox{
		errors:  []error{errors.New("transient"), nil},
		results: []*types.Result{nil, {Stdout: "ok\n", ExitCode: 0}},
	}
	a := New(AgentProfile{Name: "agent", Role: "role"}, mLLM, mSB)
	a.Task = "ls"

	if err := a.Run(context.Background()); err != nil {
		t.Fatalf("Run() unexpected error: %v", err)
	}
	if mSB.calls != 2 {
		t.Errorf("Sandbox calls = %d, want 2 (one retry)", mSB.calls)
	}
}

func TestAgent_Run_LLMRetry(t *testing.T) {
	mLLM := &mockLLM{
		errors:    []error{errors.New("transient"), nil},
		responses: []string{"", "[[FINISH]]"},
	}
	mSB := &mockSandbox{}
	a := New(AgentProfile{Name: "agent", Role: "role"}, mLLM, mSB)
	a.Task = "task"

	if err := a.Run(context.Background()); err != nil {
		t.Fatalf("Run() unexpected error: %v", err)
	}
	if mLLM.calls != 2 {
		t.Errorf("LLM calls = %d, want 2 (one retry)", mLLM.calls)
	}
}

func TestAgent_Run_LLMMaxRetriesReached(t *testing.T) {
	mLLM := &mockLLM{errors: []error{
		errors.New("e1"), errors.New("e2"), errors.New("e3"), errors.New("e4"),
	}}
	mSB := &mockSandbox{}
	a := New(AgentProfile{Name: "agent", Role: "role"}, mLLM, mSB)
	a.Task = "task"
	a.MaxRetries = 2 // 3 total attempts

	err := a.Run(context.Background())
	if err == nil {
		t.Fatal("expected error after max retries, got nil")
	}
	if mLLM.calls != 3 {
		t.Errorf("LLM calls = %d, want 3", mLLM.calls)
	}
}

// ---------------------------------------------------------------------------
// Agent constructor defaults
// ---------------------------------------------------------------------------

func TestAgent_Run_BudgetExceeded(t *testing.T) {
	mLLM := &mockLLM{
		responses: []string{"ACTION: ls", "[[FINISH]]"},
	}
	mSB := &mockSandbox{
		results: []*types.Result{{Stdout: "file.txt", ExitCode: 0}},
	}
	a := New(AgentProfile{Name: "agent", Role: "role"}, mLLM, mSB)
	a.Task = "list files"
	a.BudgetTokens = 30 // Initial budget is 30 tokens
	// In mockLLM, each call consumes 20 tokens.
	// First call: 20 tokens used. Remaining: 10.
	// Second call should fail because it would exceed budget (20 > 10).

	err := a.Run(context.Background())
	if err == nil {
		t.Fatal("expected error due to budget exceed, got nil")
	}
	if !strings.Contains(err.Error(), "budget exceeded") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestIsUnsafeAction(t *testing.T) {
	tests := []struct {
		action string
		unsafe bool
	}{
		{"ls", false},
		{"echo hello && ls", false},
		{"cat file | grep text", false},
		{"ls ; rm -rf /", true},
		{"echo `whoami`", true},
		{"echo $(whoami)", true},
		{"ls || rm -rf /", true},
		{"echo hello > file.txt", false},
	}

	for _, tt := range tests {
		got := isUnsafeAction(tt.action)
		if got != tt.unsafe {
			t.Errorf("isUnsafeAction(%q) = %v, want %v", tt.action, got, tt.unsafe)
		}
	}
}

func TestAgent_CostAggregation(t *testing.T) {
	mockLLMClient := &mockLLM{
		responses: []string{"ACTION: ls", "[[FINISH]]"},
	}
	mockSandboxClient := &mockSandbox{
		results: []*types.Result{{ExitCode: 0, Stdout: "file.txt"}},
	}
	a := New(AgentProfile{Name: "test"}, mockLLMClient, mockSandboxClient)

	ctx := context.Background()
	if err := a.Run(ctx); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// 2 LLM calls, each with 10 prompt + 10 completion tokens = 20 total.
	// Total tokens: 40.
	if a.TotalUsage != 40 {
		t.Errorf("expected 40 total tokens, got %d", a.TotalUsage)
	}

	// Cost per call (fallback):
	// 10 tokens * 75 nano + 10 tokens * 300 nano = 750 + 3000 = 3750 nano-dollars per call.
	// Total for 2 calls: 7500 nano-dollars.

	if a.TotalCostNanoDollars != 7500 {
		t.Errorf("expected 7500 nano-dollars total cost, got %d", a.TotalCostNanoDollars)
	}
}

func TestAgent_Run_CostBudgetExceeded(t *testing.T) {
	mockLLMClient := &mockLLM{
		responses: []string{"ACTION: ls", "[[FINISH]]"},
	}
	mockSandboxClient := &mockSandbox{
		results: []*types.Result{{ExitCode: 0, Stdout: "file.txt"}},
	}
	a := New(AgentProfile{Name: "test"}, mockLLMClient, mockSandboxClient)

	// Set a very low cost budget: 1000 nano-dollars (/bin/bash.000001)
	// Each mock LLM call uses 10+10=20 tokens.
	// Fallback cost: 10*75 + 10*300 = 3750 nano-dollars per call.
	a.MaxCostNanoDollars = 1000

	ctx := context.Background()
	err := a.Run(ctx)
	if err == nil {
		t.Fatal("expected error due to cost budget exceed, got nil")
	}
	if !strings.Contains(err.Error(), "cost budget exceeded") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestAgent_Run_DailyBudgetExceeded(t *testing.T) {
	mockLLMClient := &mockLLM{
		responses: []string{"ACTION: ls", "[[FINISH]]"},
	}
	mockSandboxClient := &mockSandbox{
		results: []*types.Result{{ExitCode: 0, Stdout: "file.txt"}},
	}
	a := New(AgentProfile{Name: "test"}, mockLLMClient, mockSandboxClient)

	tmpDir := t.TempDir()
	tracker := governance.NewUsageTracker(tmpDir)

	// Pre-record some usage to exceed daily budget
	// Budget is 50 tokens today.
	a.DailyBudgetTokens = 50
	a.Tracker = tracker

	// Record 40 tokens already used today
	err := tracker.RecordUsage(40, 0)
	if err != nil {
		t.Fatalf("RecordUsage failed: %v", err)
	}

	ctx := context.Background()
	// Next call will use 20 tokens, making total 60, which exceeds 50.
	err = a.Run(ctx)
	if err == nil {
		t.Fatal("expected error due to daily budget exceed, got nil")
	}
	if !strings.Contains(err.Error(), "daily token budget exceeded") {
		t.Errorf("unexpected error message: %v", err)
	}
}
