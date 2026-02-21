package orchestrator

import (
	"testing"
)

func TestEpicStatusTransitions(t *testing.T) {
	tests := []struct {
		current EpicStatus
		signal  string
		want    EpicStatus
	}{
		{StatusPlanned, "lisa_ready", StatusReady},
		{StatusReady, "tick", StatusInProgress},
		{StatusInProgress, "ralph_done", StatusImplemented},
		{StatusImplemented, "bart_ok", StatusVerified},
		{StatusImplemented, "bart_fail_implementation", StatusInProgress},
		{StatusImplemented, "bart_fail_viability", StatusBlocked},
		{StatusImplemented, "bart_fail_adr", StatusBlocked},
		{StatusVerified, "lovejoy_merge", StatusDone},
		{StatusBlocked, "lisa_redecide", StatusReady},
	}

	for _, tt := range tests {
		got, err := tt.current.Transition(tt.signal)
		if err != nil {
			t.Errorf("%s + %s: unexpected error: %v", tt.current, tt.signal, err)
			continue
		}
		if got != tt.want {
			t.Errorf("%s + %s: got %s, want %s", tt.current, tt.signal, got, tt.want)
		}
	}
}

func TestInvalidTransitions(t *testing.T) {
	_, err := StatusDone.Transition("any")
	if err == nil {
		t.Errorf("expected error transitioning from StatusDone, got nil")
	}
}
