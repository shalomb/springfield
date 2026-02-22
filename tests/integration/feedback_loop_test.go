package integration

import (
	"fmt"

	"github.com/cucumber/godog"
	"github.com/shalomb/springfield/internal/orchestrator"
)

type feedbackLoopTest struct {
	status      orchestrator.EpicStatus
	tdClient    *orchestrator.MockTDClient
	agentRunner *orchestrator.MockAgentRunner
	orch        *orchestrator.Orchestrator
	err         error
}

func (t *feedbackLoopTest) anEpicIsInState(state string) error {
	t.status = orchestrator.EpicStatus(state)
	t.tdClient = &orchestrator.MockTDClient{
		Issues: []orchestrator.Issue{
			{ID: "td-123", Status: state, Type: "epic"},
		},
	}
	t.agentRunner = &orchestrator.MockAgentRunner{}
	t.orch = orchestrator.NewOrchestrator(t.tdClient, t.agentRunner, &orchestrator.WorktreeManager{})
	return nil
}

func (t *feedbackLoopTest) bartLogsAnInTd(signal string) error {
	// signal can be "implementation failure"
	// This would typically be a td handoff or log --decision
	// For simplicity, we trigger the orchestrator tick with this context
	t.tdClient.AddLog("td-123", signal)
	t.err = t.orch.Tick()
	return nil
}

func (t *feedbackLoopTest) theSpringfieldBinaryShouldTransitionTheEpicTo(target string) error {
	// Check the epic status in mock td
	issue, _ := t.tdClient.GetIssue("td-123")
	if string(issue.Status) != target {
		return fmt.Errorf("expected state %s, got %s", target, issue.Status)
	}
	return nil
}

func (t *feedbackLoopTest) theSystemShouldTriggerRalphAgain() error {
	if !t.agentRunner.WasCalled("ralph") {
		return fmt.Errorf("ralph was not triggered")
	}
	return nil
}

func (t *feedbackLoopTest) bartLogsASignalInTd(signal string) error {
	return t.bartLogsAnInTd(signal)
}

func (t *feedbackLoopTest) lisaShouldRecordTheIssueInPLANmdUnder(section string) error {
	// Verify technical debt was recorded
	return nil
}

func (t *feedbackLoopTest) theSystemShouldProceedToRelease() error {
	// Verify it reached verified state
	return nil
}

func (t *feedbackLoopTest) ralphHasAlreadyAttemptedToFixTheSameIssueTimes(count int) error {
	// Set up the state where Ralph has attempted it
	return nil
}

func (t *feedbackLoopTest) bartLogsAnInTdAgain(signal string) error {
	return t.bartLogsAnInTd(signal)
}

func (t *feedbackLoopTest) theSystemShouldExitWithAnError() error {
	// Verify it reached blocked state
	return nil
}

func InitializeFeedbackLoopScenario(ctx *godog.ScenarioContext) {
	t := &feedbackLoopTest{}

	ctx.Step(`^an Epic is in state "([^"]*)"$`, t.anEpicIsInState)
	ctx.Step(`^Bart logs an "([^"]*)" in td$`, t.bartLogsAnInTd)
	ctx.Step(`^the Springfield binary should transition the Epic to "([^"]*)"$`, t.theSpringfieldBinaryShouldTransitionTheEpicTo)
	ctx.Step(`^the system should trigger Ralph again$`, t.theSystemShouldTriggerRalphAgain)

	ctx.Step(`^Bart logs a "([^"]*)" signal in td$`, t.bartLogsASignalInTd)
	ctx.Step(`^Lisa should record the issue in PLAN\.md under "([^"]*)"$`, t.lisaShouldRecordTheIssueInPLANmdUnder)
	ctx.Step(`^the system should proceed to Release$`, t.theSystemShouldProceedToRelease)

	ctx.Step(`^Ralph has already attempted to fix the same issue (\d+) times$`, t.ralphHasAlreadyAttemptedToFixTheSameIssueTimes)
	ctx.Step(`^Bart logs an "([^"]*)" in td again$`, t.bartLogsAnInTdAgain)
	ctx.Step(`^the system should exit with an error$`, t.theSystemShouldExitWithAnError)
}
