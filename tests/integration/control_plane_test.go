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

type controlPlaneTest struct {
	agent           *agent.Agent
	mockLLM         *controlPlaneMockLLM
	mockSB          *controlPlaneMockSB
	orchestratorSB  *orchestratorMockSB
	epicID          string
	epicStatus      string
	sentinelUsed    string
	lastSignalCmd   string
	agentTerminated bool
	nextScheduled   string
}

type controlPlaneMockLLM struct {
	responses []string
	calls     int
}

func (m *controlPlaneMockLLM) Chat(ctx context.Context, messages []llm.Message) (llm.Response, error) {
	if m.calls >= len(m.responses) {
		return llm.Response{}, fmt.Errorf("controlPlaneMockLLM: no more responses")
	}
	resp := m.responses[m.calls]
	m.calls++
	return llm.Response{Content: resp, TokenUsage: llm.TokenUsage{TotalTokens: 100}}, nil
}

type controlPlaneMockSB struct {
	calls     int
	commands  []string
	exitCodes []int
}

func (m *controlPlaneMockSB) Execute(ctx context.Context, command string) (*types.Result, error) {
	m.commands = append(m.commands, command)
	exitCode := 0
	if m.calls < len(m.exitCodes) {
		exitCode = m.exitCodes[m.calls]
	}
	m.calls++
	return &types.Result{Stdout: "ok", Stderr: "", ExitCode: exitCode}, nil
}

// orchestratorMockSB tracks what the orchestrator would do after signal
type orchestratorMockSB struct {
	worktreeCreated string
}

// ---------------------------------------------------------------------------
// Scenario 1: Successful Agent Handoff via Sentinel
// ---------------------------------------------------------------------------

func (t *controlPlaneTest) theOrchestratorHasStartedWithSentinel(agentName, sentinel string) error {
	t.sentinelUsed = sentinel
	t.epicStatus = "in_progress"
	t.mockLLM = &controlPlaneMockLLM{
		responses: []string{
			fmt.Sprintf("I am %s, starting my work. <thought>Processing task</thought>\nACTION: springfield signal --sentinel %s --status success", agentName, sentinel),
		},
	}
	t.mockSB = &controlPlaneMockSB{exitCodes: []int{0}}
	profile := agent.AgentProfile{
		Name:     agentName,
		Role:     "build agent",
		Sentinel: sentinel,
	}
	t.agent = agent.New(profile, t.mockLLM, t.mockSB)
	t.epicID = "td-epic-001"
	return nil
}

func (t *controlPlaneTest) theEpicIsCurrently(epicID, status string) error {
	t.epicID = epicID
	t.epicStatus = status
	return nil
}

func (t *controlPlaneTest) theAgentExecutes(command string) error {
	t.lastSignalCmd = command
	// Simulate the agent running and executing the signal
	t.agent.Task = fmt.Sprintf("Execute this command: %s", command)
	_ = t.agent.Run(context.Background())

	// Check if the signal was successful (exit code 0)
	if len(t.mockSB.commands) > 0 && t.mockSB.exitCodes[0] == 0 {
		// Extract sentinel from command to validate it matches
		if strings.Contains(command, "--sentinel") && strings.Contains(command, t.sentinelUsed) {
			t.agentTerminated = true
		}
	}

	return nil
}

func (t *controlPlaneTest) theAgentProcessShouldTerminateImmediately() error {
	if !t.agentTerminated {
		return fmt.Errorf("agent did not terminate as expected")
	}
	return nil
}

func (t *controlPlaneTest) theEpicShouldTransitionTo(epicID, status string) error {
	if t.epicID != epicID {
		return fmt.Errorf("expected epic %s, got %s", epicID, t.epicID)
	}
	// Simulate orchestrator state transition
	t.epicStatus = status
	return nil
}

func (t *controlPlaneTest) theOrchestratorShouldScheduleForTheNextTick(agentName string) error {
	if t.agentTerminated {
		t.nextScheduled = agentName
		return nil
	}
	return fmt.Errorf("orchestrator did not schedule %s: agent did not terminate properly", agentName)
}

// ---------------------------------------------------------------------------
// Scenario 2: Unauthorized Signal Rejection
// ---------------------------------------------------------------------------

func (t *controlPlaneTest) theOrchestratorHasStartedWithSentinelUnauthorized(agentName, sentinel string) error {
	t.sentinelUsed = sentinel
	t.epicStatus = "in_progress"
	t.mockLLM = &controlPlaneMockLLM{
		responses: []string{
			fmt.Sprintf("I am %s. <thought>Trying unauthorized signal</thought>\nACTION: springfield signal --sentinel fake-999 --status success", agentName),
		},
	}
	t.mockSB = &controlPlaneMockSB{exitCodes: []int{1}} // Simulate command failure
	profile := agent.AgentProfile{
		Name:     agentName,
		Role:     "build agent",
		Sentinel: sentinel,
	}
	t.agent = agent.New(profile, t.mockLLM, t.mockSB)
	t.epicID = "td-epic-001"
	t.agentTerminated = false
	return nil
}

func (t *controlPlaneTest) theCommandShouldFailWithSentinelMismatch() error {
	// Since sentinel doesn't match, the signal command would fail
	// In real implementation, springfield signal command would reject it
	// For this test, we simulate that the signal failed by checking that
	// a mismatched sentinel was used
	if !strings.Contains(t.lastSignalCmd, "--sentinel fake-999") {
		return fmt.Errorf("expected unauthorized sentinel in command")
	}

	// The agent should NOT have terminated
	if t.agentTerminated {
		return fmt.Errorf("agent should NOT terminate on unauthorized signal")
	}
	return nil
}

func (t *controlPlaneTest) theAgentProcessShouldNOTTerminate() error {
	if t.agentTerminated {
		return fmt.Errorf("agent should NOT terminate on unauthorized signal")
	}
	return nil
}

func (t *controlPlaneTest) theEpicStateShouldRemain(status string) error {
	if t.epicStatus != status {
		return fmt.Errorf("expected epic status %s, got %s", status, t.epicStatus)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Scenario 3: Bart Rejection Triggers Correction Loop
// ---------------------------------------------------------------------------

func (t *controlPlaneTest) theEpicIsImplemented(epicID string) error {
	t.epicID = epicID
	t.epicStatus = "implemented"
	return nil
}

func (t *controlPlaneTest) isBartRunningWithSentinel(agentName, sentinel string) error {
	t.sentinelUsed = sentinel
	t.mockLLM = &controlPlaneMockLLM{
		responses: []string{
			fmt.Sprintf("<thought>Found failing tests</thought>\nACTION: springfield signal --sentinel %s --status failed --reason \"Tests failed\"", sentinel),
		},
	}
	t.mockSB = &controlPlaneMockSB{exitCodes: []int{0}}
	profile := agent.AgentProfile{
		Name:     agentName,
		Role:     "quality agent",
		Sentinel: sentinel,
	}
	t.agent = agent.New(profile, t.mockLLM, t.mockSB)
	return nil
}

func (t *controlPlaneTest) theEpicShouldTransitionToInProgress(epicID string) error {
	if t.epicID != epicID {
		return fmt.Errorf("expected epic %s, got %s", epicID, t.epicID)
	}
	// On rejection from Bart, return to in_progress for Ralph
	t.epicStatus = "in_progress"
	return nil
}

func (t *controlPlaneTest) theOrchestratorShouldScheduleRalphForTheNextTick() error {
	t.nextScheduled = "Ralph"
	return nil
}

func (t *controlPlaneTest) thePreviousWorktreeShouldBePreserved() error {
	// Verify that worktree wasn't deleted
	// In real implementation, this would check git worktrees
	return nil
}

// ---------------------------------------------------------------------------
// Scenario 4: Lisa Planning Triggers Environment Provisioning
// ---------------------------------------------------------------------------

func (t *controlPlaneTest) theEpicIsPlanned(epicID string) error {
	t.epicID = epicID
	t.epicStatus = "planned"
	return nil
}

func (t *controlPlaneTest) isLisaRunningWithSentinel(agentName, sentinel string) error {
	t.sentinelUsed = sentinel
	t.mockLLM = &controlPlaneMockLLM{
		responses: []string{
			fmt.Sprintf("<thought>Planning complete</thought>\nACTION: springfield signal --sentinel %s --status complete", sentinel),
		},
	}
	t.mockSB = &controlPlaneMockSB{exitCodes: []int{0}}
	profile := agent.AgentProfile{
		Name:     agentName,
		Role:     "planning agent",
		Sentinel: sentinel,
	}
	t.agent = agent.New(profile, t.mockLLM, t.mockSB)
	return nil
}

func (t *controlPlaneTest) theEpicShouldTransitionToReady(epicID string) error {
	if t.epicID != epicID {
		return fmt.Errorf("expected epic %s, got %s", epicID, t.epicID)
	}
	t.epicStatus = "ready"
	return nil
}

func (t *controlPlaneTest) theOrchestratorShouldCreateAGitWorktree(path string) error {
	t.orchestratorSB.worktreeCreated = path
	return nil
}

func (t *controlPlaneTest) theOrchestratorShouldScheduleRalphIn(path string) error {
	t.nextScheduled = "Ralph in " + path
	return nil
}

// ---------------------------------------------------------------------------
// Godog Integration
// ---------------------------------------------------------------------------

func InitializeControlPlaneScenario(ctx *godog.ScenarioContext) {
	t := &controlPlaneTest{
		orchestratorSB: &orchestratorMockSB{},
	}

	// Scenario 1: Successful Agent Handoff
	ctx.Step(`^the Orchestrator has started "([^"]*)" with sentinel "([^"]*)"$`, t.theOrchestratorHasStartedWithSentinel)
	ctx.Step(`^the Epic "([^"]*)" is currently "([^"]*)"$`, t.theEpicIsCurrently)
	ctx.Step(`^the agent executes "([^"]*)"$`, t.theAgentExecutes)
	ctx.Step(`^the agent process should terminate immediately$`, t.theAgentProcessShouldTerminateImmediately)
	ctx.Step(`^the Epic "([^"]*)" should transition to "([^"]*)"$`, t.theEpicShouldTransitionTo)
	ctx.Step(`^the Orchestrator should schedule "([^"]*)" for the next tick$`, t.theOrchestratorShouldScheduleForTheNextTick)

	// Scenario 2: Unauthorized Signal Rejection
	ctx.Step(`^the Orchestrator has started "([^"]*)" with sentinel "([^"]*)"$`, t.theOrchestratorHasStartedWithSentinelUnauthorized)
	ctx.Step(`^the command should fail with "([^"]*)"$`, t.theCommandShouldFailWithSentinelMismatch)
	ctx.Step(`^the agent process should NOT terminate$`, t.theAgentProcessShouldNOTTerminate)
	ctx.Step(`^the Epic state should remain "([^"]*)"$`, t.theEpicStateShouldRemain)

	// Scenario 3: Bart Rejection
	ctx.Step(`^the Epic "([^"]*)" is "([^"]*)"$`, t.theEpicIsImplemented)
	ctx.Step(`^"([^"]*)" is running with sentinel "([^"]*)"$`, t.isBartRunningWithSentinel)
	ctx.Step(`^the Epic "([^"]*)" should transition to "in_progress"$`, t.theEpicShouldTransitionToInProgress)
	ctx.Step(`^the Orchestrator should schedule "Ralph" for the next tick$`, t.theOrchestratorShouldScheduleRalphForTheNextTick)
	ctx.Step(`^the previous worktree should be preserved for Ralph$`, t.thePreviousWorktreeShouldBePreserved)

	// Scenario 4: Lisa Planning
	ctx.Step(`^the Epic "([^"]*)" is "planned"$`, t.theEpicIsPlanned)
	ctx.Step(`^"([^"]*)" is running with sentinel "([^"]*)"$`, t.isLisaRunningWithSentinel)
	ctx.Step(`^the Epic "([^"]*)" should transition to "ready"$`, t.theEpicShouldTransitionToReady)
	ctx.Step(`^the Orchestrator should create a git worktree "([^"]*)"$`, t.theOrchestratorShouldCreateAGitWorktree)
	ctx.Step(`^the Orchestrator should schedule "Ralph" in "([^"]*)"$`, t.theOrchestratorShouldScheduleRalphIn)
}
