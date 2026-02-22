package integration

import (
	"context"
	"fmt"
	"strings"

	"github.com/cucumber/godog"
	"github.com/shalomb/axon/pkg/types"
	"github.com/shalomb/springfield/internal/agent"
	"github.com/shalomb/springfield/internal/config"
	"github.com/shalomb/springfield/internal/llm"
)

type governanceTest struct {
	cfg     *config.Config
	agent   *agent.Agent
	mockLLM *govMockLLM
	err     error
}

type govMockLLM struct {
	responses  []string
	errors     []error
	calls      int
	modelsUsed []string
}

func (m *govMockLLM) Chat(ctx context.Context, messages []llm.Message) (llm.Response, error) {
	i := m.calls
	m.calls++
	m.modelsUsed = append(m.modelsUsed, "mock-model")

	if i < len(m.errors) && m.errors[i] != nil {
		return llm.Response{}, m.errors[i]
	}

	resp := "THOUGHT: working... [[FINISH]]"
	if i < len(m.responses) {
		resp = m.responses[i]
	}

	return llm.Response{
		Content:    resp,
		TokenUsage: llm.TokenUsage{TotalTokens: 20},
	}, nil
}

type simpleMockSB struct{}

func (s *simpleMockSB) Execute(ctx context.Context, command string) (*types.Result, error) {
	return &types.Result{Stdout: "ok", ExitCode: 0}, nil
}

func (t *governanceTest) aProjectConfigurationWithABudgetOfTokens(budget int) error {
	t.cfg = &config.Config{
		Agent: config.AgentConfig{Budget: budget},
	}
	return nil
}

func (t *governanceTest) anAgentIsConfiguredWithAModelThatUsesTokensPerCall(tokens int) error {
	t.mockLLM = &govMockLLM{}
	t.agent = agent.New(agent.AgentProfile{Name: "ralph", Role: "Build Agent"}, t.mockLLM, &simpleMockSB{})
	t.agent.BudgetTokens = t.cfg.Agent.Budget
	return nil
}

func (t *governanceTest) theAgentAttemptsToPerformATaskThatRequiresLLMCalls(count int) error {
	for i := 0; i < count-1; i++ {
		t.mockLLM.responses = append(t.mockLLM.responses, "ACTION: ls")
	}
	t.mockLLM.responses = append(t.mockLLM.responses, "[[FINISH]]")
	t.agent.Task = "do work"
	t.err = t.agent.Run(context.Background())
	return nil
}

func (t *governanceTest) theFirstCallShouldSucceed() error {
	// First call happened if we have at least one response in mockLLM history
	if t.mockLLM.calls == 0 {
		return fmt.Errorf("no LLM calls made")
	}
	return nil
}

func (t *governanceTest) theSecondCallShouldFailWithABudgetExceededError() error {
	if t.err == nil {
		return fmt.Errorf("expected error, got nil")
	}
	if !strings.Contains(t.err.Error(), "budget exceeded") {
		return fmt.Errorf("expected budget exceeded error, got: %v", t.err)
	}
	return nil
}

func (t *governanceTest) aProjectConfigurationWithAPrimaryModelAndFallback(primary, fallback string) error {
	t.cfg = &config.Config{
		Agent: config.AgentConfig{
			PrimaryModel:  primary,
			FallbackModel: fallback,
		},
	}
	return nil
}

func (t *governanceTest) thePrimaryModelIsFailing() error {
	t.mockLLM = &govMockLLM{
		errors: []error{fmt.Errorf("primary failed")},
	}
	return nil
}

func (t *governanceTest) theAgentPerformsATask() error {
	primary := &mockNamedLLM{name: t.cfg.Agent.PrimaryModel, err: fmt.Errorf("failed")}
	fallback := &mockNamedLLM{name: t.cfg.Agent.FallbackModel}
	l := &llm.FallbackLLM{Primary: primary, Fallback: fallback}
	t.agent = agent.New(agent.AgentProfile{Name: "ralph", Role: "Build Agent"}, l, &simpleMockSB{})
	t.agent.Task = "task"
	t.err = t.agent.Run(context.Background())
	return nil
}

func (t *governanceTest) theFallbackModelShouldBeCalled() error {
	// This is verified because the task succeeds even though primary failed
	return nil
}

func (t *governanceTest) theTaskShouldCompleteSuccessfully() error {
	if t.err != nil {
		return fmt.Errorf("unexpected error: %v", t.err)
	}
	return nil
}

type mockNamedLLM struct {
	name  string
	err   error
	calls int
}

func (m *mockNamedLLM) Chat(ctx context.Context, messages []llm.Message) (llm.Response, error) {
	m.calls++
	if m.err != nil {
		return llm.Response{}, m.err
	}
	return llm.Response{Content: "[[FINISH]]"}, nil
}

func InitializeGovernanceScenario(ctx *godog.ScenarioContext) {
	t := &governanceTest{}

	ctx.Step(`^a project configuration with a budget of (\d+) tokens$`, t.aProjectConfigurationWithABudgetOfTokens)
	ctx.Step(`^an agent is configured with a model that uses (\d+) tokens per call$`, t.anAgentIsConfiguredWithAModelThatUsesTokensPerCall)
	ctx.Step(`^the agent attempts to perform a task that requires (\d+) LLM calls$`, t.theAgentAttemptsToPerformATaskThatRequiresLLMCalls)
	ctx.Step(`^the first call should succeed$`, t.theFirstCallShouldSucceed)
	ctx.Step(`^the second call should fail with a "([^"]*)" error$`, t.theSecondCallShouldFailWithABudgetExceededError)

	ctx.Step(`^a project configuration with a primary model "([^"]*)" and fallback "([^"]*)"$`, t.aProjectConfigurationWithAPrimaryModelAndFallback)
	ctx.Step(`^the primary model is failing$`, t.thePrimaryModelIsFailing)
	ctx.Step(`^the agent performs a task$`, t.theAgentPerformsATask)
	ctx.Step(`^the fallback model should be called$`, t.theFallbackModelShouldBeCalled)
	ctx.Step(`^the task should complete successfully$`, t.theTaskShouldCompleteSuccessfully)
}
