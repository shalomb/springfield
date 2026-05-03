package main

import (
	"bytes"
	"testing"
)

func TestOrchestrateCommandExists(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"orchestrate", "--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("orchestrate command failed: %v", err)
	}

	if buf.String() == "" {
		t.Errorf("orchestrate command output should not be empty")
	}
	t.Logf("Output: %s", buf.String())
}

func TestOrchestrateDaemonFlag(t *testing.T) {
	daemon = false // reset global flag
	rootCmd.SetArgs([]string{"orchestrate", "--daemon"})
	// We don't execute it because it would loop forever
	// but we can check if the flag is registered.
	f := orchestrateCmd.Flag("daemon")
	if f == nil {
		t.Errorf("orchestrate command missing --daemon flag")
	}
}
