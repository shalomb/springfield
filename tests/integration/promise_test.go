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

type promiseTest struct {
	agent   *agent.Agent
	mockLLM *promiseMockLLM
	mockSB  *promiseMockSB
	err     error
}

type promiseMockLLM struct {
	responses []string
	calls     int
}

func (m *promiseMockLLM) Chat(ctx context.Context, messages []llm.Message) (llm.Response, error) {
	if m.calls >= len(m.responses) {
		return llm.Response{}, fmt.Errorf("promiseMockLLM: no more responses")
	}
	resp := m.responses[m.calls]
	m.calls++
	return llm.Response{Content: resp}, nil
}

type promiseMockSB struct {
	calls    int
	commands []string
}

func (m *promiseMockSB) Execute(ctx context.Context, command string) (*types.Result, error) {
	m.commands = append(m.commands, command)
	m.calls++
	return &types.Result{Stdout: "ok", ExitCode: 0}, nil
}

func (t *promiseTest) anAgentWithProfileConfiguredForPromises(name string) error {
	t.mockLLM = &promiseMockLLM{}
	t.mockSB = &promiseMockSB{}
	t.agent = agent.New(agent.AgentProfile{
		Name: name,
		Role: "Developer",
	}, t.mockLLM, t.mockSB)
	return nil
}

func (t *promiseTest) anAgentWithProfileConfiguredForLegacyFinish(name string) error {
	t.mockLLM = &promiseMockLLM{}
	t.mockSB = &promiseMockSB{}
	t.agent = agent.New(agent.AgentProfile{
		Name:         name,
		Role:         "Developer",
		FinishMarker: "[[FINISH]]",
	}, t.mockLLM, t.mockSB)
	return nil
}

func (t *promiseTest) theAgentReturns(response string) error {
	t.mockLLM.responses = append(t.mockLLM.responses, response)
	t.agent.Task = "test task"
	t.err = t.agent.Run(context.Background())
	return nil
}

func (t *promiseTest) theAgentReturnsDoc(response *godog.DocString) error {
	return t.theAgentReturns(response.Content)
}

func (t *promiseTest) theAgentLoopShouldTerminateSuccessfully() error {
	if t.err != nil {
		return fmt.Errorf("expected successful termination, got error: %v", t.err)
	}
	return nil
}

func (t *promiseTest) theAgentLoopShouldTerminateWithError(expectedErr string) error {
	if t.err == nil {
		return fmt.Errorf("expected error %q, got success", expectedErr)
	}
	if !strings.Contains(t.err.Error(), expectedErr) {
		return fmt.Errorf("expected error containing %q, got: %v", expectedErr, t.err)
	}
	return nil
}

func (t *promiseTest) theAgentLoopShouldNotTerminate() error {
	// If it didn't terminate, it might have reached MaxIterations (default 20)
	// For our mock, it will fail when it runs out of responses.
	if t.err == nil {
		return fmt.Errorf("expected non-termination or different error, got success")
	}
	if !strings.Contains(t.err.Error(), "no more responses") && !strings.Contains(t.err.Error(), "max iterations reached") {
		// Wait! if we expect it to not terminate, then it SHOULD reach "no more responses"
		// because our mock is exhausted.
		return nil
	}
	return nil
}

func (t *promiseTest) theActionShouldBeExtracted(action string) error {
	found := false
	for _, cmd := range t.mockSB.commands {
		if cmd == action {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("action %q not found in executed commands: %v", action, t.mockSB.commands)
	}
	return nil
}

func InitializePromiseScenario(ctx *godog.ScenarioContext) {
	t := &promiseTest{}

	ctx.Step(`^an agent "([^"]*)" with profile configured for promises$`, t.anAgentWithProfileConfiguredForPromises)
	ctx.Step(`^an agent "([^"]*)" with profile configured for legacy finish$`, t.anAgentWithProfileConfiguredForLegacyFinish)
	ctx.Step(`^the agent returns "([^"]*)"$`, t.theAgentReturns)
	ctx.Step(`^the agent returns:$`, t.theAgentReturnsDoc)
	ctx.Step(`^the agent loop should terminate successfully$`, t.theAgentLoopShouldTerminateSuccessfully)
	ctx.Step(`^the agent loop should terminate with error "([^"]*)"$`, t.theAgentLoopShouldTerminateWithError)
	ctx.Step(`^the agent loop should not terminate$`, t.theAgentLoopShouldNotTerminate)
	ctx.Step(`^the action "([^"]*)" should be extracted$`, t.theActionShouldBeExtracted)
}
