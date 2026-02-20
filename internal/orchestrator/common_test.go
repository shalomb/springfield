package orchestrator

import (
	"os/exec"
	"testing"
)

func runCmd(t *testing.T, dir string, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("cmd %s %v failed in %s: %v (output: %s)", name, args, dir, err, string(output))
	}
}
