package agent

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RalphRunner implements the multi-iteration loop for Ralph.
// It continues executing until TODO.md is empty and there are no uncommitted changes.
type RalphRunner struct {
	*BaseRunner
	maxLoops int // For testing: limit the number of loops to prevent infinite runs
}

// Run executes Ralph's multi-iteration loop.
// The loop continues while TODO.md exists OR there are uncommitted changes.
// It exits when TODO.md is removed AND there are no uncommitted changes.
func (rr *RalphRunner) Run(ctx context.Context) error {
	// Set default maxLoops if not set (for production use)
	if rr.maxLoops == 0 {
		rr.maxLoops = 100 // Arbitrary large number for production
	}

	loopCount := 0

	for loopCount < rr.maxLoops {
		loopCount++

		// Check if we should continue looping
		todoExists, err := RalphTODOExists()
		if err != nil {
			return fmt.Errorf("error checking TODO.md existence: %w", err)
		}

		hasChanges, err := hasUncommittedChanges()
		if err != nil {
			return fmt.Errorf("error checking git status: %w", err)
		}

		// Exit condition: no TODO.md and no uncommitted changes
		if !todoExists && !hasChanges {
			fmt.Println("✅ No TODO.md found and no uncommitted changes. Work complete!")
			return nil
		}

		// Log status
		if hasChanges {
			fmt.Println("📝 Uncommitted changes detected. Engaging Ralph to finalize...")
		} else {
			fmt.Println("📋 Tasks found in TODO.md. Engaging Ralph...")
		}

		// Append TODO.md content to the task if it exists
		originalTask := rr.Task
		if todoExists {
			todoContent, err := readTODOFile()
			if err == nil && todoContent != "" {
				rr.Task = originalTask + "\n\nCurrent TODO.md:\n" + todoContent
			}
		}

		// Execute one iteration
		err = rr.BaseRunner.Run(ctx)
		if err != nil {
			return fmt.Errorf("error in Ralph iteration %d: %w", loopCount, err)
		}

		// Restore original task
		rr.Task = originalTask

		fmt.Println("\n********")
	}

	return fmt.Errorf("exceeded maximum loop iterations (%d)", rr.maxLoops)
}

// RalphTODOExists checks whether TODO.md exists in the current directory.
func RalphTODOExists() (bool, error) {
	_, err := os.Stat("TODO.md")
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// hasUncommittedChanges checks whether there are uncommitted changes using git.
func hasUncommittedChanges() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain", "--untracked-files=no")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("git status command failed: %w", err)
	}

	// If output is empty, there are no uncommitted changes
	return strings.TrimSpace(string(output)) != "", nil
}

// readTODOFile reads the content of TODO.md from disk.
func readTODOFile() (string, error) {
	content, err := os.ReadFile("TODO.md")
	if err != nil {
		return "", err
	}
	return string(content), nil
}
