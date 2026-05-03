package agent

import (
	"testing"
)

// TestIsSignalAction tests the signal action detection logic.
func TestIsSignalAction(t *testing.T) {
	tests := []struct {
		name   string
		action string
		want   bool
	}{
		{
			name:   "valid signal command",
			action: "springfield signal --sentinel abc123 --status success",
			want:   true,
		},
		{
			name:   "signal command with reason",
			action: "springfield signal --sentinel abc123 --status success --reason \"All tests pass\"",
			want:   true,
		},
		{
			name:   "signal command at start of multiline (should reject)",
			action: "springfield signal\n--sentinel abc123 --status success",
			want:   false,
		},
		{
			name:   "signal command with embedded newline (should reject)",
			action: "springfield signal --sentinel abc123\n--status success",
			want:   false,
		},
		{
			name:   "signal as substring (should reject)",
			action: "some springfield signal command",
			want:   false,
		},
		{
			name:   "not a signal command",
			action: "ls -la /tmp",
			want:   false,
		},
		{
			name:   "empty action",
			action: "",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isSignalAction(tt.action)
			if got != tt.want {
				t.Errorf("isSignalAction(%q) = %v, want %v", tt.action, got, tt.want)
			}
		})
	}
}

// TestExtractSentinel tests sentinel extraction with strict UUIDv4 validation.
func TestExtractSentinel(t *testing.T) {
	tests := []struct {
		name         string
		action       string
		wantSentinel string
		wantOk       bool
		description  string
	}{
		{
			name:         "valid UUIDv4 sentinel",
			action:       "springfield signal --sentinel d330877b-7f7b-4e3a-ac5d-2e8c3f9a1b2c --status success",
			wantSentinel: "d330877b-7f7b-4e3a-ac5d-2e8c3f9a1b2c",
			wantOk:       true,
			description:  "Standard UUIDv4 format should match",
		},
		{
			name:         "valid UUIDv4 with only sentinel flag",
			action:       "springfield signal --sentinel 550e8400-e29b-41d4-a716-446655440000",
			wantSentinel: "550e8400-e29b-41d4-a716-446655440000",
			wantOk:       true,
			description:  "UUIDv4 alone should match",
		},
		{
			name:         "sentinel followed by double dash without value (should fail)",
			action:       "springfield signal --sentinel  --status success",
			wantSentinel: "",
			wantOk:       false,
			description:  "Empty sentinel value should not match",
		},
		{
			name:         "sentinel followed by another flag (should fail)",
			action:       "springfield signal --sentinel --status success",
			wantSentinel: "",
			wantOk:       false,
			description:  "Flag name instead of value should not match",
		},
		{
			name:         "sentinel with arbitrary text (should fail)",
			action:       "springfield signal --sentinel some-random-text --status success",
			wantSentinel: "",
			wantOk:       false,
			description:  "Non-UUIDv4 text should not match",
		},
		{
			name:         "command substitution attempt (should fail)",
			action:       "springfield signal --sentinel $(whoami) --status success",
			wantSentinel: "",
			wantOk:       false,
			description:  "Command substitution should not match",
		},
		{
			name:         "file path injection attempt (should fail)",
			action:       "springfield signal --sentinel /etc/passwd --status success",
			wantSentinel: "",
			wantOk:       false,
			description:  "File path should not match",
		},
		{
			name:         "UUIDv4 with uppercase (should match - case insensitive)",
			action:       "springfield signal --sentinel D330877B-7F7B-4E3A-AC5D-2E8C3F9A1B2C --status success",
			wantSentinel: "D330877B-7F7B-4E3A-AC5D-2E8C3F9A1B2C",
			wantOk:       true,
			description:  "Uppercase UUIDv4 should match",
		},
		{
			name:         "UUIDv4 without hyphens (should fail)",
			action:       "springfield signal --sentinel d330877b7f7b4e3aac5d2e8c3f9a1b2c --status success",
			wantSentinel: "",
			wantOk:       false,
			description:  "UUIDv4 without hyphens should not match",
		},
		{
			name:         "partial UUIDv4 (should fail)",
			action:       "springfield signal --sentinel d330877b-7f7b-4e3a --status success",
			wantSentinel: "",
			wantOk:       false,
			description:  "Incomplete UUIDv4 should not match",
		},
		{
			name:         "no --sentinel flag",
			action:       "springfield signal --status success",
			wantSentinel: "",
			wantOk:       false,
			description:  "Missing --sentinel flag should not match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSentinel, gotOk := extractSentinel(tt.action)
			if gotSentinel != tt.wantSentinel || gotOk != tt.wantOk {
				t.Errorf("extractSentinel(%q)\n  got sentinel=%q ok=%v\n  want sentinel=%q ok=%v\n  description: %s",
					tt.action, gotSentinel, gotOk, tt.wantSentinel, tt.wantOk, tt.description)
			}
		})
	}
}
