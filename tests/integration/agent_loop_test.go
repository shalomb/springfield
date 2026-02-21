package integration

import (
	"context"
	"fmt"
	"strings"

	"github.com/cucumber/godog"
	"github.com/shalomb/axon/pkg/types"
	"github.com/shalomb/springfield/internal/agent"
	"github.com/shalomb/springfield/internal/llm"
)

type agentLoopTest struct {
	agent   *agent.Agent
	mockLLM *bddMockLLM
	mockSB  *bddMockSB
	err     error
}

type bddMockLLM struct {
	responses []string
	calls     int
	received  [][]llm.Message
}

func (m *bddMockLLM) Chat(ctx context.Context, messages []llm.Message) (llm.Response, error) {
	if m.calls >= len(m.responses) {
		return llm.Response{}, fmt.Errorf("bddMockLLM: no more responses")
	}
	cp := make([]llm.Message, len(messages))
	copy(cp, messages)
	m.received = append(m.received, cp)
	resp := m.responses[m.calls]
	m.calls++
	return llm.Response{Content: resp}, nil
}

type bddMockSB struct {
	calls   int
	outputs []string
}

func (m *bddMockSB) Execute(ctx context.Context, command string) (*types.Result, error) {
	out := "file1.txt\n"
	if m.calls < len(m.outputs) {
		out = m.outputs[m.calls]
	}
	m.calls++
	return &types.Result{Stdout: out, ExitCode: 0}, nil
}

// ---------------------------------------------------------------------------
// Step definitions
// ---------------------------------------------------------------------------

func (t *agentLoopTest) anAgentNamedWithRole(name, role string) error {
	t.mockLLM = &bddMockLLM{}
	t.mockSB = &bddMockSB{outputs: []string{"hello from sandbox\n"}}
	t.agent = agent.New(agent.AgentProfile{Name: name, Role: role}, t.mockLLM, t.mockSB)
	return nil
}

func (t *agentLoopTest) theLLMIsConfiguredWithResponses(table *godog.Table) error {
	for _, row := range table.Rows[1:] { // skip header
		var resp string
		if row.Cells[2].Value == "true" {
			resp = fmt.Sprintf("THOUGHT: %s.\n%s", row.Cells[0].Value, agent.FinishMarker)
		} else {
			resp = fmt.Sprintf("THOUGHT: %s.\nACTION: %s", row.Cells[0].Value, row.Cells[1].Value)
		}
		t.mockLLM.responses = append(t.mockLLM.responses, resp)
	}
	return nil
}

func (t *agentLoopTest) theAgentRunsTheTask(task string) error {
	t.agent.Task = task
	t.err = t.agent.Run(context.Background())
	return nil
}

func (t *agentLoopTest) theAgentShouldHaveCalledTheLLMTimes(count int) error {
	if t.mockLLM.calls != count {
		return fmt.Errorf("expected %d LLM calls, got %d", count, t.mockLLM.calls)
	}
	return nil
}

func (t *agentLoopTest) theAgentShouldHaveExecutedActionInTheSandbox(count int) error {
	if t.mockSB.calls != count {
		return fmt.Errorf("expected %d sandbox calls, got %d", count, t.mockSB.calls)
	}
	return nil
}

func (t *agentLoopTest) theTaskShouldBeSuccessful() error {
	if t.err != nil {
		return fmt.Errorf("task failed with error: %v", t.err)
	}
	return nil
}

func (t *agentLoopTest) theSecondLLMCallShouldIncludeTheSandboxOutput() error {
	if len(t.mockLLM.received) < 2 {
		return fmt.Errorf("expected at least 2 LLM calls, got %d", len(t.mockLLM.received))
	}
	msgs := t.mockLLM.received[1]
	for _, m := range msgs {
		if strings.Contains(m.Content, "hello from sandbox") {
			return nil
		}
	}
	return fmt.Errorf("sandbox output not found in second LLM call messages")
}

func (t *agentLoopTest) theSystemPromptShouldContain(text string) error {
	if len(t.mockLLM.received) == 0 {
		return fmt.Errorf("LLM was never called")
	}
	first := t.mockLLM.received[0]
	if len(first) == 0 {
		return fmt.Errorf("no messages sent to LLM")
	}
	if !strings.Contains(first[0].Content, text) {
		return fmt.Errorf("system prompt %q does not contain %q", first[0].Content, text)
	}
	return nil
}

func InitializeAgentLoopScenario(ctx *godog.ScenarioContext) {
	t := &agentLoopTest{}

	ctx.Step(`^an agent named "([^"]*)" with role "([^"]*)"$`, t.anAgentNamedWithRole)
	ctx.Step(`^the LLM is configured with responses:$`, t.theLLMIsConfiguredWithResponses)
	ctx.Step(`^the agent runs the task "([^"]*)"$`, t.theAgentRunsTheTask)
	ctx.Step(`^the agent should have called the LLM (\d+) times$`, t.theAgentShouldHaveCalledTheLLMTimes)
	ctx.Step(`^the agent should have executed (\d+) actions? in the sandbox$`, t.theAgentShouldHaveExecutedActionInTheSandbox)
	ctx.Step(`^the task should be successful$`, t.theTaskShouldBeSuccessful)
	ctx.Step(`^the second LLM call should include the sandbox output$`, t.theSecondLLMCallShouldIncludeTheSandboxOutput)
	ctx.Step(`^the system prompt should contain "([^"]*)"$`, t.theSystemPromptShouldContain)
}
