package orchestrator

type EpicStatus string

const (
	StatusPlanned     EpicStatus = "planned"
	StatusReady       EpicStatus = "ready"
	StatusInProgress  EpicStatus = "in_progress"
	StatusImplemented EpicStatus = "implemented"
	StatusVerified    EpicStatus = "verified"
	StatusDone        EpicStatus = "done"
	StatusDeferred    EpicStatus = "deferred"
	StatusBlocked     EpicStatus = "blocked"
)

// Next returns the next status for an Epic based on the current status and a signal.
func (s EpicStatus) Next(signal string) EpicStatus {
	switch s {
	case StatusPlanned:
		if signal == "lisa_ready" {
			return StatusReady
		}
	case StatusReady:
		if signal == "tick" {
			return StatusInProgress
		}
	case StatusInProgress:
		if signal == "ralph_done" {
			return StatusImplemented
		}
	case StatusImplemented:
		switch signal {
		case "bart_ok":
			return StatusVerified
		case "bart_fail_impl":
			return StatusInProgress
		case "bart_fail_viability", "bart_fail_adr":
			return StatusBlocked
		}
	case StatusVerified:
		if signal == "lovejoy_merge" {
			return StatusDone
		}
	case StatusBlocked:
		if signal == "lisa_redecide" {
			return StatusReady
		}
	case StatusDeferred:
		// Needs external intervention to un-defer
	}

	return s
}
