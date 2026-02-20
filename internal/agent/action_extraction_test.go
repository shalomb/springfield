package agent

import (
	"testing"
)

func TestAgent_ActionExtraction_Safety(t *testing.T) {
	tests := []struct {
		name     string
		response string
		wantCmd  string
	}{
		{
			name:     "action on new line",
			response: "THOUGHT: check files.\nACTION: ls -la",
			wantCmd:  "ls -la",
		},
		{
			name:     "action at beginning",
			response: "ACTION: ls -la",
			wantCmd:  "ls -la",
		},
		{
			name:     "action with trailing text",
			response: "ACTION: echo hello\nMore text here",
			wantCmd:  "echo hello",
		},
		{
			name:     "empty action",
			response: "ACTION: ",
			wantCmd:  "",
		},
		{
			name:     "no action",
			response: "Just some text",
			wantCmd:  "",
		},
		{
			name:     "ACTION: inside a sentence should be ignored",
			response: "I will perform the ACTION: ls now.",
			wantCmd:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractAction(tt.response)
			if got != tt.wantCmd {
				t.Errorf("extractAction() = %q, want %q", got, tt.wantCmd)
			}
		})
	}
}

func TestAgent_ActionExtraction_ShellMetacharacters(t *testing.T) {
	// Block obviously malicious sequential commands that try to escape the intended action.
	badActions := []string{
		"ls; rm -rf /",
		"echo `rm -rf /`",
		"$(rm -rf /)",
	}

	for _, action := range badActions {
		if !isUnsafeAction(action) {
			t.Errorf("expected %q to be identified as unsafe", action)
		}
	}

	goodActions := []string{
		"ls | grep go",
		"make build && make test",
	}
	for _, action := range goodActions {
		if isUnsafeAction(action) {
			t.Errorf("expected %q to be identified as safe", action)
		}
	}
}
