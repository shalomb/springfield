package parser

import (
	"strings"
)

// MarkdownSanitizer provides functionality to strip code blocks from Markdown content.
type MarkdownSanitizer struct {
}

// NewMarkdownSanitizer creates a new MarkdownSanitizer instance.
func NewMarkdownSanitizer() *MarkdownSanitizer {
	return &MarkdownSanitizer{}
}

// StripCodeBlocks removes all content within Markdown code blocks.
// This includes triple-backtick blocks (```), tilde blocks (~~~), and indented code blocks.
// It returns the sanitized content with code blocks removed.
func (s *MarkdownSanitizer) StripCodeBlocks(input string) string {
	lines := strings.Split(input, "\n")
	var output []string
	inCodeBlock := false
	fenceChar := ""

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if inCodeBlock {
			// Check for closing fence
			if strings.HasPrefix(trimmed, fenceChar) && len(trimmed) >= 3 {
				inCodeBlock = false
				fenceChar = ""
				continue
			}
			// Skip content inside code block
			continue
		}

		// Check for opening fence
		if (strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~")) && len(trimmed) >= 3 {
			inCodeBlock = true
			if strings.HasPrefix(trimmed, "```") {
				fenceChar = "```"
			} else {
				fenceChar = "~~~"
			}
			continue
		}

		// Check for indented code block (4 spaces or 1 tab)
		if isIndented(line) {
			if i == 0 || strings.TrimSpace(lines[i-1]) == "" {
				// Skip all subsequent indented lines
				for i < len(lines) && (isIndented(lines[i]) || strings.TrimSpace(lines[i]) == "") {
					i++
				}
				i-- // Backtrack because the outer loop will increment
				continue
			}
		}

		output = append(output, line)
	}

	return strings.Join(output, "\n")
}

func isIndented(line string) bool {
	if line == "" {
		return false
	}
	return strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t")
}
