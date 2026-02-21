package orchestrator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWorktreeManager_RobustVerification(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	tempDir, err := os.MkdirTemp("", "worktree-robust-test")
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
	epicID := "td-999"
	worktreePath := filepath.Join(tempDir, "worktrees", "epic-"+epicID)

	// Case 1: Pre-existing non-worktree directory
	err = os.MkdirAll(worktreePath, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(worktreePath, "stale.txt"), []byte("I am stale"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// EnsureWorktree should handle this stale directory
	path, err := wm.EnsureWorktree(epicID)
	if err != nil {
		t.Fatalf("EnsureWorktree failed with stale directory: %v", err)
	}

	if path != worktreePath {
		t.Errorf("expected path %s, got %s", worktreePath, path)
	}

	// Verify it's now a real worktree (has .git file)
	if _, err := os.Stat(filepath.Join(worktreePath, ".git")); os.IsNotExist(err) {
		t.Errorf("worktree .git file missing after EnsureWorktree on stale directory")
	}

	// Verify stale file is gone
	if _, err := os.Stat(filepath.Join(worktreePath, "stale.txt")); err == nil {
		t.Errorf("stale file still exists in worktree")
	}

	// Case 2: Pre-existing .git directory (e.g. manually initialized)
	// We want to make sure it's treated as stale if not in `git worktree list`
	err = os.RemoveAll(worktreePath)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(filepath.Join(worktreePath, ".git"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	_, err = wm.EnsureWorktree(epicID)
	if err != nil {
		t.Fatalf("EnsureWorktree failed with manually initialized .git: %v", err)
	}

	// Verify it's now a real worktree (git worktree add should have worked)
	isWorktree, err := wm.isManagedWorktree(worktreePath)
	if err != nil {
		t.Fatal(err)
	}
	if !isWorktree {
		t.Errorf("expected it to be a managed worktree now")
	}
}
