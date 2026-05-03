package setup

import (
	"fmt"
	"os/exec"
	"strings"
)

// GitRoot returns the root of the git repository containing dir, using
// `git rev-parse --show-toplevel`. This correctly handles worktrees and
// submodules, unlike a manual .git directory walk.
func GitRoot(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository (or git not available): %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
