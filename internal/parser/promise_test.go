package parser

import (
	"testing"
)

func TestExtractPromise(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Promise
		wantErr  bool
	}{
		{
			name:     "Extract COMPLETE promise",
			input:    "Task done.<promise>COMPLETE</promise>",
			expected: PromiseComplete,
		},
		{
			name:     "Extract FAILED promise",
			input:    "Task failed.<promise>FAILED</promise>",
			expected: PromiseFailed,
		},
		{
			name:     "Return UNKNOWN if no promise present",
			input:    "Just text without tags.",
			expected: PromiseUnknown,
		},
		{
			name:     "Case-insensitive extraction",
			input:    "<promise>complete</promise>",
			expected: PromiseComplete,
		},
		{
			name:     "Promise with surrounding whitespace",
			input:    "<promise>  COMPLETE  </promise>",
			expected: PromiseComplete,
		},
		{
			name:     "Promise in code block (ignored)",
			input:    "```\n<promise>COMPLETE</promise>\n```",
			expected: PromiseUnknown,
		},
		{
			name:     "Multiple promise tags (use first)",
			input:    "<promise>COMPLETE</promise><promise>FAILED</promise>",
			expected: PromiseComplete,
		},
		{
			name:     "Invalid promise values (error)",
			input:    "<promise>OK</promise>",
			expected: PromiseUnknown,
			wantErr:  true,
		},
		{
			name:     "Invalid promise values (case mismatch and value mismatch)",
			input:    "<promise>done</promise>",
			expected: PromiseUnknown,
			wantErr:  true,
		},
		{
			name:     "Empty promise tag (error)",
			input:    "<promise></promise>",
			expected: PromiseUnknown,
			wantErr:  true,
		},
		{
			name:     "Promise at start of response",
			input:    "<promise>COMPLETE</promise> The task is finished.",
			expected: PromiseComplete,
		},
		{
			name:     "Promise at end of response",
			input:    "The task is finished. <promise>COMPLETE</promise>",
			expected: PromiseComplete,
		},
		{
			name:     "Promise with multiline content",
			input:    "<promise>\nCOMPLETE\n</promise>",
			expected: PromiseComplete,
		},
		{
			name:     "Malformed tag (unclosed)",
			input:    "<promise>COMPLETE",
			expected: PromiseUnknown,
		},
		{
			name:     "Mixed case with tags",
			input:    "<pRoMiSe>cOmPlEtE</pRoMiSe>",
			expected: PromiseComplete,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractPromise(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractPromise() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ExtractPromise() = %v, want %v", got, tt.expected)
			}
		})
	}
}
