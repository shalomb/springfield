package parser

import (
	"testing"
)

func TestMarkdownSanitizer_StripCodeBlocks(t *testing.T) {
	s := NewMarkdownSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Strip single code block",
			input:    "Before\n```bash\nrm -rf /\n```\nAfter",
			expected: "Before\nAfter",
		},
		{
			name:     "Preserve text outside code blocks",
			input:    "Text 1\n```\ncode\n```\nText 2\n~~~go\nmore code\n~~~\nText 3",
			expected: "Text 1\nText 2\nText 3",
		},
		{
			name:     "Handle nested code blocks (ignore different fences inside)",
			input:    "Outer\n```\nInner\n~~~\nStill inner?\n~~~\n```\nEnd",
			expected: "Outer\nEnd",
		},
		{
			name: "Indented code blocks stripped correctly",
			input: `Some text

    ls -la
    pwd

More text`,
			expected: `Some text

More text`,
		},
		{
			name:     "Tilde code blocks stripped correctly",
			input:    "Start\n~~~\nInside\n~~~\nEnd",
			expected: "Start\nEnd",
		},
		{
			name:     "Multiple code blocks in single response",
			input:    "One\n```\n1\n```\nTwo\n```\n2\n```\nThree",
			expected: "One\nTwo\nThree",
		},
		{
			name:     "Code block without closing marker (unclosed)",
			input:    "Before\n```\nInside\nAnd still inside",
			expected: "Before",
		},
		{
			name:     "Code block at start of response",
			input:    "```\nStart\n```\nText",
			expected: "Text",
		},
		{
			name:     "Code block at end of response",
			input:    "Text\n```\nEnd\n```",
			expected: "Text",
		},
		{
			name:     "Empty code block",
			input:    "Before\n```\n```\nAfter",
			expected: "Before\nAfter",
		},
		{
			name:     "Tags within code blocks are ignored",
			input:    "Before\n```\n<action>echo hello</action>\n```\nAfter",
			expected: "Before\nAfter",
		},
		{
			name:     "Mixed fenced and indented",
			input:    "Fenced:\n```\ncode\n```\nIndented:\n\n    more code\n\nFinal",
			expected: "Fenced:\nIndented:\n\nFinal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.StripCodeBlocks(tt.input)
			if got != tt.expected {
				t.Errorf("StripCodeBlocks() = %q, want %q", got, tt.expected)
			}
		})
	}
}
