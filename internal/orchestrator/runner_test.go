package orchestrator

import (
	"os"
	"strings"
	"testing"
)

func TestCommandAgentRunner_Run(t *testing.T) {
	// We'll use a script that just prints its environment and arguments to a file
	// so we can verify them.
	tempDir := t.TempDir()
	logFile := tempDir + "/log.txt"
	scriptPath := tempDir + "/fake_springfield"

	script := `#!/bin/sh
echo "ARGS: $@" > ` + logFile + `
env >> ` + logFile + `
`
	_ = os.WriteFile(scriptPath, []byte(script), 0755)

	runner := &CommandAgentRunner{BinaryPath: scriptPath}
	err := runner.Run("ralph", "epic-123", tempDir)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	logContent, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read log: %v", err)
	}

	content := string(logContent)
	if !strings.Contains(content, "--agent ralph") {
		t.Errorf("missing --agent ralph in args: %s", content)
	}
	if !strings.Contains(content, "--epic epic-123") {
		t.Errorf("missing --epic epic-123 in args: %s", content)
	}
	if !strings.Contains(content, "--sentinel ") {
		t.Errorf("missing --sentinel in args: %s", content)
	}
	if !strings.Contains(content, "SPRINGFIELD_SENTINEL=") {
		t.Errorf("missing SPRINGFIELD_SENTINEL in env: %s", content)
	}
	if !strings.Contains(content, "SPRINGFIELD_EPIC=epic-123") {
		t.Errorf("missing SPRINGFIELD_EPIC in env: %s", content)
	}
	if !strings.Contains(content, "SPRINGFIELD_AGENT=ralph") {
		t.Errorf("missing SPRINGFIELD_AGENT in env: %s", content)
	}
}
