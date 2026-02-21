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
