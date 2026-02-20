package orchestrator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWorktreeManager(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "worktree-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize a git repo in tempDir
	runCmd(t, tempDir, "git", "init")
	runCmd(t, tempDir, "git", "config", "user.email", "test@example.com")
	runCmd(t, tempDir, "git", "config", "user.name", "Test User")
	runCmd(t, tempDir, "git", "commit", "--allow-empty", "-m", "initial commit")

	wm := &WorktreeManager{BaseDir: tempDir}
	epicID := "td-123456"

	worktreePath, err := wm.EnsureWorktree(epicID)
	if err != nil {
		t.Fatalf("EnsureWorktree failed: %v", err)
	}

	expectedPath := filepath.Join(tempDir, "worktrees", "epic-td-123456")
	if worktreePath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, worktreePath)
	}

	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		t.Errorf("worktree directory %s was not created", worktreePath)
	}

	// Verify git branch exists
	runCmd(t, worktreePath, "git", "rev-parse", "--abbrev-ref", "HEAD")

	// Test Handoff Deposit
	handoffFile := "TODO-td-123456.md"
	err = os.WriteFile(filepath.Join(tempDir, handoffFile), []byte("test handoff"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = wm.DepositHandoff(epicID)
	if err != nil {
		t.Fatalf("DepositHandoff failed: %v", err)
	}

	// Verify handoff was deposited in worktree as TODO.md
	depositedPath := filepath.Join(worktreePath, "TODO.md")
	if _, err := os.Stat(depositedPath); os.IsNotExist(err) {
		t.Errorf("handoff file was not deposited as TODO.md at %s", depositedPath)
	}
}

func TestWorktreeManager_ExistingBranch(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "worktree-existing-branch-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	runCmd(t, tempDir, "git", "init")
	runCmd(t, tempDir, "git", "config", "user.email", "test@example.com")
	runCmd(t, tempDir, "git", "config", "user.name", "Test User")
	runCmd(t, tempDir, "git", "commit", "--allow-empty", "-m", "initial commit")

	wm := &WorktreeManager{BaseDir: tempDir}
	epicID := "td-789"
	branchName := "feat/epic-" + epicID

	// Manually create the branch
	runCmd(t, tempDir, "git", "branch", branchName)

	// Now try to EnsureWorktree. It should succeed even if branch exists.
	_, err = wm.EnsureWorktree(epicID)
	if err != nil {
		t.Fatalf("EnsureWorktree failed with existing branch: %v", err)
	}
}
