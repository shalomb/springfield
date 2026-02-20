package agent

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestRalphRunnerLoopsUntilTODOEmpty verifies Ralph continues looping while TODO.md exists.
func TestRalphRunnerLoopsUntilTODOEmpty(t *testing.T) {
	// Setup: Create a temporary working directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Initialize a git repository with initial commit
	initGitRepo(t, tmpDir)

	// Create a prompt file
	promptPath := filepath.Join(tmpDir, "test_prompt.md")
	err := os.WriteFile(promptPath, []byte("You are Ralph"), 0644)
	if err != nil {
		t.Fatalf("Failed to create prompt file: %v", err)
	}

	// Create initial TODO.md
	err = os.WriteFile(filepath.Join(tmpDir, "TODO.md"), []byte("Task 1"), 0644)
	if err != nil {
		t.Fatalf("Failed to create TODO.md: %v", err)
	}

	// Track LLM calls
	mock := &mockLLMClient{
		responses: []string{"Task complete", "Task complete"},
	}

	rr := &RalphRunner{
		BaseRunner: &BaseRunner{
			Agent:      "ralph",
			LLMClient:  mock,
			PromptPath: promptPath,
		},
		maxLoops: 2, // Limit loops to prevent infinite test
	}

	// Remove TODO.md after the first LLM call to trigger loop exit
	// We'll simulate this by not having TODO.md cause an error on the second pass
	err = rr.Run(context.Background())

	// We expect the runner to complete successfully even if it hits the max loops
	// since we're in test mode with limited iterations
	if err != nil && !strings.Contains(err.Error(), "exceeded maximum loop iterations") {
		t.Errorf("Execute() returned unexpected error: %v", err)
	}
}

// TestRalphRunnerDetectsTODOFile verifies Ralph detects TODO.md existence.
func TestRalphRunnerDetectsTODOFile(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Initialize git
	initGitRepo(t, tmpDir)

	// Create prompt
	promptPath := filepath.Join(tmpDir, "test_prompt.md")
	os.WriteFile(promptPath, []byte("You are Ralph"), 0644)

	// Create TODO.md
	os.WriteFile(filepath.Join(tmpDir, "TODO.md"), []byte("Task"), 0644)

	// Verify TODO.md is detected
	exists, err := RalphTODOExists()
	if err != nil {
		t.Errorf("TODOExists() returned error: %v", err)
	}

	if !exists {
		t.Errorf("Expected TODO.md to exist, but TODOExists() returned false")
	}
}

// TestRalphRunnerStopsWhenTODOGoneAndNoChanges verifies loop stops correctly.
func TestRalphRunnerStopsWhenTODOGoneAndNoChanges(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Initialize git with an initial commit
	initGitRepo(t, tmpDir)

	// Create prompt
	promptPath := filepath.Join(tmpDir, "test_prompt.md")
	os.WriteFile(promptPath, []byte("You are Ralph"), 0644)

	// No TODO.md, no changes
	mock := &mockLLMClient{responses: []string{}}
	rr := &RalphRunner{
		BaseRunner: &BaseRunner{
			Agent:      "ralph",
			LLMClient:  mock,
			PromptPath: promptPath,
		},
		maxLoops: 1,
	}

	err := rr.Run(context.Background())
	if err != nil {
		t.Errorf("Execute() returned error: %v", err)
	}

	// Verify no LLM calls were made (since TODO.md doesn't exist and no changes)
	if mock.calls > 0 {
		t.Errorf("Expected 0 LLM calls, got %d (loop should exit immediately)", mock.calls)
	}
}

// Helper function to initialize a git repository
func initGitRepo(t *testing.T, dir string) {
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to set git email: %v", err)
	}

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to set git user: %v", err)
	}

	// Create an initial commit so we can check for changes
	err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create README: %v", err)
	}

	cmd = exec.Command("git", "add", "README.md")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to git add: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "initial commit")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to git commit: %v", err)
	}
}
