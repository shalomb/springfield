package orchestrator

import (
	"fmt"
)

// EpicStatus represents the various states an Epic can be in during its lifecycle.
// Transitions are driven by signals from agents or the orchestrator.
type EpicStatus string

const (
	// StatusPlanned: Lisa has identified the Epic but not yet refined it.
	StatusPlanned EpicStatus = "planned"
	// StatusReady: Lisa has deposited the TODO-{id}.md and the Epic is ready for Ralph.
	StatusReady EpicStatus = "ready"
	// StatusInProgress: Ralph is currently working on the Epic.
	StatusInProgress EpicStatus = "in_progress"
	// StatusImplemented: Ralph has completed implementation and handed off to Bart.
	StatusImplemented EpicStatus = "implemented"
	// StatusVerified: Bart has verified the implementation and handed off to Lovejoy.
	StatusVerified EpicStatus = "verified"
	// StatusDone: Lovejoy has merged the changes and the Epic is complete.
	StatusDone EpicStatus = "done"
	// StatusDeferred: The Epic has been put on hold.
	StatusDeferred EpicStatus = "deferred"
	// StatusBlocked: An issue was found that requires Lisa's intervention.
	StatusBlocked EpicStatus = "blocked"
)

func (s EpicStatus) Transition(signal string) (EpicStatus, error) {
	switch s {
	case StatusPlanned:
		if signal == "lisa_ready" {
			return StatusReady, nil
		}
	case StatusReady:
		if signal == "tick" {
			return StatusInProgress, nil
		}
	case StatusInProgress:
		if signal == "ralph_done" {
			return StatusImplemented, nil
		}
	case StatusImplemented:
		switch signal {
		case "bart_ok":
			return StatusVerified, nil
		case "bart_fail_implementation":
			return StatusInProgress, nil
		case "bart_fail_viability":
			return StatusBlocked, nil
		case "bart_fail_adr":
			return StatusBlocked, nil
		}
	case StatusVerified:
		if signal == "lovejoy_merge" {
			return StatusDone, nil
		}
	case StatusBlocked:
		if signal == "lisa_redecide" {
			return StatusReady, nil
		}
	}

	return "", fmt.Errorf("invalid transition from %s with signal %s", s, signal)
}
