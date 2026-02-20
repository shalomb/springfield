package orchestrator

import (
	"testing"
)

func TestEpicStatus_Transitions(t *testing.T) {
	tests := []struct {
		name   string
		from   EpicStatus
		signal string
		want   EpicStatus
	}{
		{"Planned to Ready", StatusPlanned, "lisa_ready", StatusReady},
		{"Ready to InProgress", StatusReady, "tick", StatusInProgress},
		{"InProgress to Implemented", StatusInProgress, "ralph_done", StatusImplemented},
		{"Implemented to Verified (Bart OK)", StatusImplemented, "bart_ok", StatusVerified},
		{"Implemented to InProgress (Bart Fail)", StatusImplemented, "bart_fail_impl", StatusInProgress},
		{"Implemented to Blocked (Bart Viability)", StatusImplemented, "bart_fail_viability", StatusBlocked},
		{"Implemented to Blocked (Bart ADR)", StatusImplemented, "bart_fail_adr", StatusBlocked},
		{"Verified to Done", StatusVerified, "lovejoy_merge", StatusDone},
		{"Blocked to Ready", StatusBlocked, "lisa_redecide", StatusReady},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.from.Next(tt.signal); got != tt.want {
				t.Errorf("EpicStatus.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}
