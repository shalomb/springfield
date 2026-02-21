package sandbox

import (
	"context"
	"testing"
)

func TestAxonSandbox_IsolationCheck(t *testing.T) {
	skipIfNoPodman(t)
	sb, err := NewAxonSandbox("", nil)
	if err != nil {
		t.Fatalf("NewAxonSandbox failed: %v", err)
	}

	// Try to list files in the current host directory (which shouldn't be visible in sandbox)
	result, err := sb.Execute(context.Background(), "ls AGENTS.md")
	if err != nil {
		t.Logf("Command failed with error: %v", err)
		return
	}
	t.Logf("Command stdout: %q, stderr: %q, exitCode: %d", result.Stdout, result.Stderr, result.ExitCode)
	if result.ExitCode == 0 {
		t.Errorf("Security Leak! Sandbox can see AGENTS.md from host: %s", result.Stdout)
	}
}
