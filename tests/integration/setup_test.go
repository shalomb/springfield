package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
	"github.com/google/shlex"
)

type setupTest struct {
	dir    string // working directory for the scenario
	output string // combined stdout+stderr of last command
	err    error  // exit error of last command
}

// --- Given steps ---

func (t *setupTest) aFreshGitRepository() error {
	dir, err := os.MkdirTemp("", "springfield-e2e-*")
	if err != nil {
		return err
	}
	t.dir = dir

	for _, args := range [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@example.com"},
		{"git", "config", "user.name", "Test"},
	} {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("%v failed: %w\n%s", args, err, string(out))
		}
	}
	return nil
}

func (t *setupTest) aDirectoryThatIsNotAGitRepository() error {
	dir, err := os.MkdirTemp("", "springfield-nongit-*")
	if err != nil {
		return err
	}
	t.dir = dir
	return nil
}

func (t *setupTest) initHasAlreadyBeenRun() error {
	return t.runCommand("springfield init")
}

func (t *setupTest) theTDDatabaseHasBeenInitialised() error {
	cmd := exec.Command("td", "-w", t.dir, "init")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("td init failed: %w\n%s", err, string(out))
	}
	return nil
}

// --- When steps ---

func (t *setupTest) iRunSetup(command string) error {
	return t.runCommand(command)
}

// --- Then steps ---

func (t *setupTest) theCommandShouldSucceed() error {
	if t.err != nil {
		return fmt.Errorf("expected success but got error: %v\noutput:\n%s", t.err, t.output)
	}
	return nil
}

func (t *setupTest) theCommandShouldFail() error {
	if t.err == nil {
		return fmt.Errorf("expected failure but command succeeded\noutput:\n%s", t.output)
	}
	return nil
}

func (t *setupTest) theOutputShouldContain(expected string) error {
	if !strings.Contains(t.output, expected) {
		return fmt.Errorf("expected output to contain %q\ngot:\n%s", expected, t.output)
	}
	return nil
}

func (t *setupTest) theFileShouldExistAt(relPath string) error {
	full := filepath.Join(t.dir, relPath)
	if _, err := os.Stat(full); err != nil {
		return fmt.Errorf("expected %s to exist: %w", full, err)
	}
	return nil
}

func (t *setupTest) fileShouldContain(relPath, expected string) error {
	full := filepath.Join(t.dir, relPath)
	data, err := os.ReadFile(full)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", full, err)
	}
	if !strings.Contains(string(data), expected) {
		return fmt.Errorf("%s does not contain %q\ncontents:\n%s", relPath, expected, string(data))
	}
	return nil
}

// --- helpers ---

func (t *setupTest) runCommand(command string) error {
	args, err := shlex.Split(command)
	if err != nil {
		return fmt.Errorf("cannot parse command %q: %w", command, err)
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = t.dir
	out, err := cmd.CombinedOutput()
	t.output = string(out)
	t.err = err
	return nil
}

func (t *setupTest) cleanup() {
	if t.dir != "" {
		_ = os.RemoveAll(t.dir)
	}
}

func InitializeSetupScenario(ctx *godog.ScenarioContext) {
	t := &setupTest{}

	ctx.After(func(bddCtx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		t.cleanup()
		t.dir = ""
		t.output = ""
		t.err = nil
		return bddCtx, nil
	})

	ctx.Step(`^a fresh git repository$`, t.aFreshGitRepository)
	ctx.Step(`^a directory that is not a git repository$`, t.aDirectoryThatIsNotAGitRepository)
	ctx.Step(`^"springfield init" has already been run$`, t.initHasAlreadyBeenRun)
	ctx.Step(`^the td database has been initialised$`, t.theTDDatabaseHasBeenInitialised)

	ctx.Step(`^I run setup command "([^"]*)"$`, t.iRunSetup)

	ctx.Step(`^the setup command should succeed$`, t.theCommandShouldSucceed)
	ctx.Step(`^the setup command should fail$`, t.theCommandShouldFail)
	ctx.Step(`^the setup output should contain "([^"]*)"$`, t.theOutputShouldContain)
	ctx.Step(`^the td database should exist at "([^"]*)"$`, t.theFileShouldExistAt)
	ctx.Step(`^setup file "([^"]*)" should exist$`, t.theFileShouldExistAt)
	ctx.Step(`^setup file "([^"]*)" should contain "([^"]*)"$`, t.fileShouldContain)
}
