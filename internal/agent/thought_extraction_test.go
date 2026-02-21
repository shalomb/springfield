package agent

import (
	"testing"
)

func TestExtractThought(t *testing.T) {
	tests := []struct {
		name     string
		response string
		want     string
	}{
		{
			name:     "simple thought",
			response: "<thought>I need to think.</thought>",
			want:     "I need to think.",
		},
		{
			name:     "multiline thought",
			response: "<thought>\nLine 1\nLine 2\n</thought>",
			want:     "Line 1\nLine 2",
		},
		{
			name:     "thought with other text",
			response: "Explanation:\n<thought>Reasoning</thought>\nAction: ...",
			want:     "Reasoning",
		},
		{
			name:     "no thought",
			response: "Just some text",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractThought(tt.response)
			if got != tt.want {
				t.Errorf("extractThought() = %q, want %q", got, tt.want)
			}
		})
	}
}
