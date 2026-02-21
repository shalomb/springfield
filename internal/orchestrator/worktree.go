package orchestrator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// WorktreeManager manages git worktrees for Epics.
type WorktreeManager struct {
	BaseDir string
}

// EnsureWorktree ensures that a git worktree exists for the given Epic ID.
// It returns the path to the worktree.
func (m *WorktreeManager) EnsureWorktree(epicID string) (string, error) {
	worktreePath := filepath.Join(m.BaseDir, "worktrees", "epic-"+epicID)
	branchName := "feat/epic-" + epicID

	if info, err := os.Stat(worktreePath); err == nil {
		if !info.IsDir() {
			return "", fmt.Errorf("expected directory but found file at %s", worktreePath)
		}

		// Verify if it's a managed worktree
		isWorktree, err := m.isManagedWorktree(worktreePath)
		if err != nil {
			return "", fmt.Errorf("failed to verify worktree: %w", err)
		}

		// Secondary check: .git file presence
		_, gitErr := os.Stat(filepath.Join(worktreePath, ".git"))
		if isWorktree && gitErr == nil {
			return worktreePath, nil
		}

		// Not a worktree, remove it
		if err := os.RemoveAll(worktreePath); err != nil {
			return "", fmt.Errorf("failed to remove stale directory %s: %w", worktreePath, err)
		}
	}

	// Check if branch exists
	branchExists := false
	checkBranchCmd := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+branchName)
	checkBranchCmd.Dir = m.BaseDir
	if err := checkBranchCmd.Run(); err == nil {
		branchExists = true
	}

	// Create worktree
	var cmd *exec.Cmd
	if branchExists {
		// Use existing branch
		cmd = exec.Command("git", "worktree", "add", worktreePath, branchName)
	} else {
		// Create new branch
		cmd = exec.Command("git", "worktree", "add", worktreePath, "-b", branchName)
	}
	cmd.Dir = m.BaseDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("git worktree add failed: %w (output: %s)", err, string(output))
	}

	return worktreePath, nil
}

// DepositHandoff copies the TODO-{id}.md handoff document into the worktree as TODO.md.
func (m *WorktreeManager) DepositHandoff(epicID string) error {
	handoffFile := "TODO-" + epicID + ".md"
	sourcePath := filepath.Join(m.BaseDir, handoffFile)
	worktreePath := filepath.Join(m.BaseDir, "worktrees", "epic-"+epicID)
	destPath := filepath.Join(worktreePath, "TODO.md")

	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("handoff file not found: %s", sourcePath)
	}

	input, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	return os.WriteFile(destPath, input, 0644)
}

func (m *WorktreeManager) isManagedWorktree(path string) (bool, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}

	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = m.BaseDir
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "worktree ") {
			wtPath := strings.TrimPrefix(line, "worktree ")
			wtAbsPath, err := filepath.Abs(wtPath)
			if err != nil {
				continue
			}
			if wtAbsPath == absPath {
				return true, nil
			}
		}
	}

	return false, nil
}
