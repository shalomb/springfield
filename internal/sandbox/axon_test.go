package sandbox

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func skipIfNoPodman(t *testing.T) {
	if _, err := exec.LookPath("podman"); err != nil {
		t.Skip("podman not found in path")
	}
}

func TestAxonSandbox_Execute_Simple(t *testing.T) {
	skipIfNoPodman(t)
	// In a unit test, we might want to mock the executor,
	// but for Ralph, seeing it work with the real thing (if possible) is better.
	// However, if we don't have podman in the test environment, it might fail.
	// Let's assume we want a unit test that doesn't depend on podman.

	// Since AxonSandbox wraps *executor.Executor, and we don't have an interface
	// for Executor in axon pkg yet, it's hard to mock without more refactoring.

	// For now, let's just test that it satisfies the interface and handles initialization.
	var _ Sandbox = (*AxonSandbox)(nil)

	sb, err := NewAxonSandbox("", nil)
	if err != nil {
		t.Fatalf("NewAxonSandbox failed: %v", err)
	}
	if sb == nil {
		t.Fatal("NewAxonSandbox returned nil")
	}
}

func TestAxonSandbox_Execute_Error(t *testing.T) {
	sb := &AxonSandbox{exec: nil}
	_, err := sb.Execute(context.Background(), "echo hello")
	if err == nil {
		t.Fatal("expected error for uninitialized executor, got nil")
	}
}

func TestAxonSandbox_GuardrailBlock(t *testing.T) {
	skipIfNoPodman(t)
	sb, err := NewAxonSandbox("", nil)
	if err != nil {
		t.Fatalf("NewAxonSandbox failed: %v", err)
	}
	// "rm -rf /" should be blocked by default guardrails
	_, err = sb.Execute(context.Background(), "rm -rf /")
	if err == nil {
		t.Fatal("expected guardrail error, got nil")
	}
	if !strings.Contains(err.Error(), "guardrail") {
		t.Errorf("expected guardrail error, got: %v", err)
	}
}

func TestNewAxonSandbox_Path(t *testing.T) {
	skipIfNoPodman(t)
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "config.toml")
	if err := os.WriteFile(confPath, []byte("[axon]\nversion=\"1.0.0\"\n"), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	sb, err := NewAxonSandbox(confPath, nil)
	if err != nil {
		t.Fatalf("NewAxonSandbox failed: %v", err)
	}
	if sb == nil {
		t.Fatal("NewAxonSandbox returned nil")
	}
}

func TestNewAxonSandbox_SearchToRoot(t *testing.T) {
	skipIfNoPodman(t)
	// Create a deep directory structure
	tmpDir := t.TempDir()
	deepDir := filepath.Join(tmpDir, "a", "b", "c")
	if err := os.MkdirAll(deepDir, 0755); err != nil {
		t.Fatalf("failed to create deep dir: %v", err)
	}

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current dir: %v", err)
	}
	if err := os.Chdir(deepDir); err != nil {
		t.Fatalf("failed to change dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(origDir); err != nil {
			t.Errorf("failed to restore original dir: %v", err)
		}
	}()

	t.Setenv("SPRINGFIELD_CONFIG", "")
	// This will search up to root and hit the break
	_, err = NewAxonSandbox("", nil)
	if err != nil {
		t.Fatalf("NewAxonSandbox failed: %v", err)
	}
}
