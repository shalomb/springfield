package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
	"github.com/shalomb/axon/pkg/types"
	"github.com/shalomb/springfield/internal/sandbox"
)

// axonAvailable is set by probing the sandbox library in init().
var (
	axonAvailable bool
)

func init() {
	// Probe for library availability by attempting to create a sandbox.
	// We don't need a real config here; we just want to see if the library
	// can initialize (e.g. check if podman is available if that's a requirement).
	_, err := sandbox.NewAxonSandbox("", nil)
	axonAvailable = (err == nil)
}

type sandboxingTest struct {
	sb     sandbox.Sandbox
	result *types.Result
	err    error
}

// skipIfNoAxon skips the scenario when the axon binary is not present.
func (t *sandboxingTest) skipIfNoAxon(ctx context.Context) (context.Context, error) {
	if !axonAvailable {
		return ctx, godog.ErrPending
	}
	return ctx, nil
}

func (t *sandboxingTest) aSandboxEnvironmentIsConfigured() error {
	var err error
	t.sb, err = sandbox.NewAxonSandbox("", nil)
	return err
}

func (t *sandboxingTest) iRunTheCommandWithinTheSandbox(command string) error {
	t.result, t.err = t.sb.Execute(context.Background(), command)
	return nil
}

func (t *sandboxingTest) theCommandShouldSucceed() error {
	if t.err != nil {
		return t.err
	}
	if t.result.ExitCode != 0 {
		return fmt.Errorf("expected exit code 0, got %d. Stdout: %s Stderr: %s",
			t.result.ExitCode, t.result.Stdout, t.result.Stderr)
	}
	return nil
}

func (t *sandboxingTest) theOutputShouldContain(expected string) error {
	combined := t.result.Stdout + t.result.Stderr
	if !strings.Contains(combined, expected) {
		return fmt.Errorf("expected output to contain %q, but got stdout=%q stderr=%q",
			expected, t.result.Stdout, t.result.Stderr)
	}
	return nil
}

func (t *sandboxingTest) theExitCodeShouldBe(expected int) error {
	if t.result.ExitCode != expected {
		return fmt.Errorf("expected exit code %d, got %d", expected, t.result.ExitCode)
	}
	return nil
}

func (t *sandboxingTest) theStderrShouldContain(expected string) error {
	if !strings.Contains(t.result.Stderr, expected) {
		return fmt.Errorf("expected stderr to contain %q, but got %q", expected, t.result.Stderr)
	}
	return nil
}

func (t *sandboxingTest) theHostHasADirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("directory %s does not exist on host", dir)
	}
	return nil
}

func (t *sandboxingTest) iAttemptToRunWithinTheSandbox(command string) error {
	t.result, t.err = t.sb.Execute(context.Background(), command)
	return nil
}

func (t *sandboxingTest) theOperationShouldFail() error {
	if t.err == nil && t.result != nil && t.result.ExitCode == 0 {
		return fmt.Errorf("expected operation to fail, but it succeeded with output: %s", t.result.Stdout)
	}
	return nil
}

func (t *sandboxingTest) aSandboxEnvironmentIsConfiguredWithAWorkspaceVolume() error {
	var err error
	t.sb, err = sandbox.NewAxonSandbox("", nil)
	return err
}

func (t *sandboxingTest) aFileExistsInTheWorkspace(filename string) error {
	return os.WriteFile(filename, []byte("Initial content"), 0644)
}

func (t *sandboxingTest) theFileInTheHostWorkspaceShouldContain(filename, content string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	if !strings.Contains(string(data), content) {
		return fmt.Errorf("expected file %s to contain %q, but got %q", filename, content, string(data))
	}
	return nil
}

func (t *sandboxingTest) theChangesShouldPersistAfterTheSandboxExecution() error {
	// Persistence is partially verified by theFileInTheHostWorkspaceShouldContain.
	// Here we additionally verify that the file is not empty.
	info, err := os.Stat("test.txt")
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		return fmt.Errorf("test.txt is empty after sandbox execution")
	}
	return nil
}

func (t *sandboxingTest) theFileShouldNotExist(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("expected file %s to not exist, but it does", filename)
	}
	return nil
}

func (t *sandboxingTest) theExecutionMetadataShouldBePopulatedDuration() error {
	if t.result.Execution.DurationMs <= 0 {
		return fmt.Errorf("expected duration_ms > 0, got %d", t.result.Execution.DurationMs)
	}
	return nil
}

func (t *sandboxingTest) theExecutionContextShouldIdentifyProjectTypeAs(expected string) error {
	if t.result.Context.ProjectType != expected {
		return fmt.Errorf("expected project_type %q, got %q", expected, t.result.Context.ProjectType)
	}
	return nil
}

func (t *sandboxingTest) theExecutionContextShouldIdentifyBuildToolAs(expected string) error {
	if t.result.Context.BuildTool != expected {
		return fmt.Errorf("expected build_tool %q, got %q", expected, t.result.Context.BuildTool)
	}
	return nil
}

func (t *sandboxingTest) theToolsListShouldContain(expected string) error {
	found := false
	for _, tool := range t.result.Tools.List {
		if tool == expected {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("expected tools list to contain %q, but it didn't: %v", expected, t.result.Tools.List)
	}
	return nil
}

func (t *sandboxingTest) aFileExistsOnTheHost(filename string) error {
	return os.WriteFile(filename, []byte("Host only content"), 0644)
}

func (t *sandboxingTest) iRunAMemoryIntensiveProcess(ctx context.Context, limitMb int) error {
	// Try to allocate 1GB in Perl, which is definitely more than the 512m limit.
	command := "perl -e '$x = \"A\" x 1024_000_000'"
	t.result, t.err = t.sb.Execute(ctx, command)
	return nil
}

func (t *sandboxingTest) theProcessShouldBeKilledOrRestricted() error {
	// If it was killed by memory limit, exit code might be non-zero (often 137 for OOM)
	if t.err == nil && t.result != nil && t.result.ExitCode == 0 {
		return fmt.Errorf("expected process to be killed or restricted, but it succeeded")
	}
	return nil
}

func InitializeSandboxingScenario(ctx *godog.ScenarioContext) {
	t := &sandboxingTest{}

	// Cleanup host-only file after scenario
	ctx.After(func(bddCtx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		home, _ := os.UserHomeDir()
		if home != "" {
			_ = os.Remove(filepath.Join(home, "springfield_host_only.txt"))
		}
		return bddCtx, nil
	})

	// Gate all @requires_axon scenarios on binary availability.
	ctx.Before(func(bddCtx context.Context, sc *godog.Scenario) (context.Context, error) {
		for _, tag := range sc.Tags {
			if tag.Name == "@requires_axon" {
				return t.skipIfNoAxon(bddCtx)
			}
		}
		return bddCtx, nil
	})

	ctx.Step(`^a sandbox environment is configured$`, t.aSandboxEnvironmentIsConfigured)
	ctx.Step(`^a sandbox environment is configured with a workspace volume$`, t.aSandboxEnvironmentIsConfiguredWithAWorkspaceVolume)
	ctx.Step(`^a file "([^"]*)" exists on the host$`, t.aFileExistsOnTheHost)
	ctx.Step(`^I run the command "([^"]*)" within the sandbox$`, t.iRunTheCommandWithinTheSandbox)
	ctx.Step(`^I attempt to run "([^"]*)" within the sandbox$`, t.iAttemptToRunWithinTheSandbox)
	ctx.Step(`^the command should succeed$`, t.theCommandShouldSucceed)
	ctx.Step(`^the exit code should be (\d+)$`, t.theExitCodeShouldBe)
	ctx.Step(`^the output should contain "([^"]*)"$`, t.theOutputShouldContain)
	ctx.Step(`^the stderr should contain "([^"]*)"$`, t.theStderrShouldContain)
	ctx.Step(`^the host has a directory "([^"]*)"$`, t.theHostHasADirectory)
	ctx.Step(`^the operation should fail$`, t.theOperationShouldFail)
	ctx.Step(`^a file "([^"]*)" exists in the workspace$`, t.aFileExistsInTheWorkspace)
	ctx.Step(`^the file "([^"]*)" in the host workspace should contain "([^"]*)"$`, t.theFileInTheHostWorkspaceShouldContain)
	ctx.Step(`^the changes should persist after the sandbox execution$`, t.theChangesShouldPersistAfterTheSandboxExecution)
	ctx.Step(`^the file "([^"]*)" should not exist$`, t.theFileShouldNotExist)
	ctx.Step(`^the execution metadata should be populated \(duration > 0\)$`, t.theExecutionMetadataShouldBePopulatedDuration)
	ctx.Step(`^the execution context should identify project_type as "([^"]*)"$`, t.theExecutionContextShouldIdentifyProjectTypeAs)
	ctx.Step(`^the execution context should identify build_tool as "([^"]*)"$`, t.theExecutionContextShouldIdentifyBuildToolAs)
	ctx.Step(`^the tools list should contain "([^"]*)"$`, t.theToolsListShouldContain)
	ctx.Step(`^I run a memory-intensive process in the sandbox that exceeds (\d+)MB$`, t.iRunAMemoryIntensiveProcess)
	ctx.Step(`^the process should be killed or restricted$`, t.theProcessShouldBeKilledOrRestricted)
}
