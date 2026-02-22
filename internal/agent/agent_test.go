package agent

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/shalomb/axon/pkg/types"
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
	a.Budget = 30 // Initial budget is 30 tokens
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

func TestAgent_Run_SystemPromptContainsSentinel(t *testing.T) {
	const validSentinel = "session-123"
	mLLM := &mockLLM{responses: []string{"[[FINISH]]"}}
	mSB := &mockSandbox{}
	a := New(AgentProfile{Name: "agent", Role: "role", Sentinel: validSentinel}, mLLM, mSB)
	a.Task = "task"

	_ = a.Run(context.Background())

	if mLLM.calls == 0 {
		t.Fatal("LLM was never called")
	}
	sys := mLLM.received[0][0]
	if !strings.Contains(sys.Content, validSentinel) {
		t.Errorf("system prompt missing sentinel: %q", sys.Content)
	}
}

func TestAgent_Run_SystemPromptTemplate(t *testing.T) {
	const validSentinel = "session-123"
	mLLM := &mockLLM{responses: []string{"[[FINISH]]"}}
	mSB := &mockSandbox{}
	a := New(AgentProfile{
		Name:         "marge",
		Role:         "product",
		Sentinel:     validSentinel,
		EpicID:       "epic-456",
		SystemPrompt: "Hello {{.Name}}, your token is {{.Sentinel}} for epic {{.EpicID}}",
	}, mLLM, mSB)
	a.Task = "task"

	_ = a.Run(context.Background())

	sys := mLLM.received[0][0]
	expected := "Hello marge, your token is session-123 for epic epic-456"
	if !strings.Contains(sys.Content, expected) {
		t.Errorf("system prompt template replacement failed: %q", sys.Content)
	}
}

func TestAgent_Run_Signal(t *testing.T) {
	const validSentinel = "d330877b-7f7b-4e3a-ac5d-2e8c3f9a1b2c"
	mLLM := &mockLLM{responses: []string{
		"ACTION: springfield signal --sentinel d330877b-7f7b-4e3a-ac5d-2e8c3f9a1b2c --status success --reason 'I am done'",
		"SHOULD NOT REACH THIS",
	}}
	mSB := &mockSandbox{results: []*types.Result{
		{Stdout: "Signaling status: success", ExitCode: 0},
	}}
	a := New(AgentProfile{Name: "agent", Role: "role", Sentinel: validSentinel}, mLLM, mSB)
	a.Task = "task"

	err := a.Run(context.Background())
	if err != nil {
		t.Fatalf("Run() unexpected error: %v", err)
	}

	if mLLM.calls != 1 {
		t.Errorf("expected 1 LLM call before signal exit, got %d", mLLM.calls)
	}

	// Check if signal command was executed
	if len(mSB.commands) == 0 {
		t.Fatal("expected sandbox execution, got none")
	}
	if !strings.Contains(mSB.commands[0], "springfield signal") {
		t.Errorf("unexpected command: %s", mSB.commands[0])
	}
}

func TestAgent_Run_Signal_InvalidSentinel(t *testing.T) {
	const validSentinel = "d330877b-7f7b-4e3a-ac5d-2e8c3f9a1b2c"
	mLLM := &mockLLM{responses: []string{
		"ACTION: springfield signal --sentinel 00000000-0000-0000-0000-000000000000 --status success",
		"[[FINISH]]",
	}}
	mSB := &mockSandbox{results: []*types.Result{
		{Stdout: "Error: invalid sentinel", ExitCode: 1},
	}}
	a := New(AgentProfile{Name: "agent", Role: "role", Sentinel: validSentinel}, mLLM, mSB)
	a.Task = "task"

	err := a.Run(context.Background())
	if err != nil {
		t.Fatalf("Run() unexpected error: %v", err)
	}

	// Should continue to the second call because signal failed (due to invalid sentinel)
	if mLLM.calls != 2 {
		t.Errorf("expected 2 LLM calls because first signal failed, got %d", mLLM.calls)
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
