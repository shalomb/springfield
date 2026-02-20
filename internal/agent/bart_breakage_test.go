package agent

import (
	"testing"
)

func TestIsUnsafeAction_Adversarial(t *testing.T) {
	tests := []struct {
		action string
		want   bool
	}{
		{"ls", false},
		{"ls && cat foo", false},
		{"ls ; rm -rf /", true},
		{"ls || echo fail", true},
		{"$(rm -rf /)", true},
		{"`rm -rf /`", true},
		{"echo \"hello\" & rm -rf /", false},    // This should probably be blocked too!
		{"echo \"hello\"\nrm -rf /", false},     // Regex only catches first line, so this is "safe" but tricky
		{"cat /etc/shadow > output.txt", false}, // Allowed per PLAN.md
	}

	for _, tt := range tests {
		got := isUnsafeAction(tt.action)
		if got != tt.want {
			t.Errorf("isUnsafeAction(%q) = %v; want %v", tt.action, got, tt.want)
		}
	}
}
