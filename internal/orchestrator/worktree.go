package orchestrator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	if _, err := os.Stat(worktreePath); err == nil {
		return worktreePath, nil
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
