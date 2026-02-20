package integration

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/cucumber/godog"
	"github.com/google/shlex"
)

type cliTest struct {
	task   string
	output string
	err    error
}

func (t *cliTest) aTask(task string) error {
	t.task = task
	return nil
}

func (t *cliTest) iRun(command string) error {
	// command will be "springfield --agent ralph --task 'say hello'"
	// We should run the binary from bin/springfield
	// We'll assume 'springfield' in the command refers to our binary.

	args, err := shlex.Split(command)
	if err != nil {
		return err
	}

	if len(args) > 0 && args[0] == "springfield" {
		args[0] = "../../bin/springfield"
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = append(cmd.Environ(), "USE_MOCK_LLM=true")
	out, err := cmd.CombinedOutput()
	t.output = string(out)
	t.err = err
	return nil
}

func (t *cliTest) theAgentShouldReceiveTheTask(agent, task string) error {
	if !strings.Contains(t.output, fmt.Sprintf("Agent: %s", agent)) {
		return fmt.Errorf("agent %s not found in output", agent)
	}
	if !strings.Contains(t.output, fmt.Sprintf("Task: %s", task)) {
		return fmt.Errorf("task %s not found in output", task)
	}
	return nil
}

func (t *cliTest) theCliOutputShouldContain(expected string) error {
	if !strings.Contains(t.output, expected) {
		return fmt.Errorf("expected output to contain %q, but got:\n%s", expected, t.output)
	}
	return nil
}

func InitializeCliScenario(ctx *godog.ScenarioContext) {
	t := &cliTest{}

	ctx.Step(`^a task "([^"]*)"$`, t.aTask)
	ctx.Step(`^I run "([^"]*)"$`, t.iRun)
	ctx.Step(`^the agent "([^"]*)" should receive the task "([^"]*)"$`, t.theAgentShouldReceiveTheTask)
	ctx.Step(`^the CLI output should contain "([^"]*)"$`, t.theCliOutputShouldContain)
}
