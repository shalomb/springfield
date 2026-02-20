package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// QuotaExceededError represents an API quota/rate limit error
type QuotaExceededError struct {
	Message  string
	Original error
}

func (e *QuotaExceededError) Error() string {
	return fmt.Sprintf("API quota exceeded: %s", e.Message)
}

// IsQuotaExceededError checks if an error is a QuotaExceededError
func IsQuotaExceededError(err error) bool {
	_, ok := err.(*QuotaExceededError)
	return ok
}

// PiLLM implements LLMClient by calling the 'pi' CLI.
type PiLLM struct {
	Model    string
	executor func(ctx context.Context, name string, arg ...string) ([]byte, error)
}

func (p *PiLLM) Chat(ctx context.Context, messages []Message) (Response, error) {
	logger := GetLogger("PiLLM.Chat")

	// Log call details
	logger.Debugf("Starting LLM call with %d messages", len(messages))
	for i, msg := range messages {
		logger.Debugf("  Message %d (role=%s): %d chars", i, msg.Role, len(msg.Content))
	}

	args := []string{"-p", "--no-tools"}

	// Pass the model if configured
	// The pi CLI does recognize "provider/model" format
	if p.Model != "" {
		args = append(args, "--model", p.Model)
		logger.Debugf("Using model: %s", p.Model)
	} else {
		logger.Debugf("No model specified, using pi CLI defaults")
	}

	// For now, pi CLI doesn't seem to have a temperature flag in this mock implementation
	// but we could add it if it did.

	var systemPrompt string
	var otherMessages []string

	for _, msg := range messages {
		if msg.Role == "system" {
			systemPrompt = msg.Content
		} else {
			otherMessages = append(otherMessages, msg.Content)
		}
	}

	if systemPrompt != "" {
		args = append(args, "--system-prompt", systemPrompt)
	}
	args = append(args, otherMessages...)

	execFn := p.executor
	if execFn == nil {
		// Default executor tries 'pi' CLI first, then falls back to npm exec.
		// This ensures Springfield works even if pi isn't in PATH.
		execFn = p.executorWithFallback
	}

	logger.Debugf("Executing pi CLI...")
	out, err := execFn(ctx, "pi", args...)
	if err != nil {
		logger.WithError(err).Errorf("pi CLI execution failed")
		return Response{}, err
	}

	// For now, pi CLI doesn't return token usage, so we'll leave it at zero.
	response := Response{Content: string(out)}
	logger.Debugf("LLM call completed. Response: %d chars", len(response.Content))
	return response, nil
}

// executorWithFallback tries to run 'pi' directly, then falls back to 'npm exec'.
// This ensures the system works in environments where pi isn't in the PATH.
func (p *PiLLM) executorWithFallback(ctx context.Context, name string, arg ...string) ([]byte, error) {
	logger := GetLogger("executorWithFallback")

	// Try 'pi' directly first
	logger.Debugf("Attempting to execute: %s with %d arguments", name, len(arg))
	cmd := exec.CommandContext(ctx, name, arg...)
	out, err := cmd.Output()
	if err == nil {
		logger.Debugf("Successfully executed %s, got %d bytes", name, len(out))
		return out, nil
	}

	// Check if error is "command not found" (pi not in PATH)
	// If so, fall back to npm exec
	if isCommandNotFound(err) {
		logger.Debugf("%s not found in PATH, falling back to npm exec", name)

		// Try 'npm exec' as fallback
		// 'npm exec @mariozechner/pi-coding-agent -- <args>'
		npmArgs := []string{"exec", "@mariozechner/pi-coding-agent", "--"}
		npmArgs = append(npmArgs, arg...)
		logger.Debugf("Executing: npm exec @mariozechner/pi-coding-agent with %d arguments", len(arg))

		// Capture stdout and stderr separately for better error reporting
		// but also stream them to console for real-time visibility
		cmd := exec.CommandContext(ctx, "npm", npmArgs...)
		var stdout, stderr bytes.Buffer

		// Pipe stdout and stderr to both buffer and console
		cmd.Stdout = io.MultiWriter(&stdout, os.Stdout)
		cmd.Stderr = io.MultiWriter(&stderr, os.Stderr)

		npmErr := cmd.Run()
		// Flush output streams to ensure they display immediately
		os.Stdout.Sync()
		os.Stderr.Sync()

		stdoutBytes := stdout.Bytes()
		stderrBytes := stderr.Bytes()
		stdoutStr := string(stdoutBytes)
		stderrStr := string(stderrBytes)

		if npmErr == nil {
			logger.Debugf("npm exec succeeded, got %d bytes stdout", len(stdoutBytes))
			// Filter out npm warnings and only return actual output from pi
			filtered := filterNpmOutput(stdoutBytes)
			logger.Debugf("After filtering npm output: %d bytes", len(filtered))
			return filtered, nil
		}

		// npm failed - provide detailed error information

		logger.Debugf("npm exec failed with stderr: %s", stderrStr)
		logger.Debugf("npm exec stdout: %s", stdoutStr)

		// Check for quota/rate limit errors (these are terminal conditions)
		if isQuotaExceeded(stderrStr) {
			errMsg := formatExecutionError("npm exec", npmErr, stderrStr, stdoutStr)
			logger.WithError(npmErr).Errorf("QUOTA EXCEEDED: %s", errMsg)
			return nil, &QuotaExceededError{
				Message:  errMsg,
				Original: npmErr,
			}
		}

		// Build a detailed error message
		errMsg := formatExecutionError("npm exec", npmErr, stderrStr, stdoutStr)
		logger.WithError(npmErr).Errorf("npm exec failed: %s", errMsg)

		return nil, fmt.Errorf("npm exec failed: %s", errMsg)
	}

	// If pi failed for reasons other than "not found", return that error
	logger.WithError(err).Errorf("command execution failed")
	return nil, err
}

// formatExecutionError creates a detailed error message from command execution failure
func formatExecutionError(cmdName string, err error, stderr, stdout string) string {
	var details string

	if stderr != "" {
		details = stderr
	} else if stdout != "" {
		details = stdout
	} else {
		details = err.Error()
	}

	// Try to extract Anthropic API error message if present
	if strings.Contains(details, "rate_limit_error") {
		if extracted := extractAnthropicErrorMessage(details); extracted != "" {
			details = extracted
		}
	}

	return details
}

// extractAnthropicErrorMessage parses Anthropic API error JSON and extracts the message
func extractAnthropicErrorMessage(stderr string) string {
	// Try to find JSON in the error message
	// Anthropic errors look like: Error: 429 {"type":"error",...}
	startIdx := strings.Index(stderr, `{"type":"error"`)
	if startIdx == -1 {
		return ""
	}

	// Find the end of the JSON object
	endIdx := strings.LastIndex(stderr, "}")
	if endIdx == -1 || endIdx <= startIdx {
		return ""
	}

	jsonStr := stderr[startIdx : endIdx+1]

	// Parse the JSON
	var errObj map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &errObj); err != nil {
		return ""
	}

	// Extract nested error message
	if errData, ok := errObj["error"].(map[string]interface{}); ok {
		if errType, ok := errData["type"].(string); ok {
			if message, ok := errData["message"].(string); ok {
				return fmt.Sprintf("Anthropic API error (%s): %s", errType, message)
			}
		}
	}

	return ""
}

// isQuotaExceeded checks if the error is due to API quota/rate limiting
func isQuotaExceeded(stderr string) bool {
	// Check for common quota/rate limit error patterns
	quotaPatterns := []string{
		"429",                     // HTTP 429 Too Many Requests
		"exhausted your capacity", // Google Gemini
		"rate limit",              // Generic rate limit
		"quota",                   // Generic quota error
		"too many requests",       // Generic rate limit message
		"request limit exceeded",  // Some APIs
		"billing_exception",       // Anthropic billing
		"401",                     // Unauthorized (may include quota)
		"403",                     // Forbidden (may include quota)
	}

	stderrLower := strings.ToLower(stderr)
	for _, pattern := range quotaPatterns {
		if strings.Contains(stderrLower, pattern) {
			return true
		}
	}
	return false
}

// isCommandNotFound checks if an error is due to a command not being found.
func isCommandNotFound(err error) bool {
	if err == nil {
		return false
	}
	// Check for "executable file not found" or "command not found" messages
	// These can be prefixed with "exec: " depending on context
	errMsg := err.Error()
	return strings.Contains(errMsg, "executable file not found") ||
		strings.Contains(errMsg, "command not found") ||
		strings.Contains(errMsg, "no such file or directory")
}

// filterNpmOutput removes npm warnings from the output while preserving actual content.
func filterNpmOutput(out []byte) []byte {
	lines := strings.Split(string(out), "\n")
	var result []string
	for _, line := range lines {
		// Skip npm warnings
		if !strings.HasPrefix(strings.TrimSpace(line), "npm warn") {
			result = append(result, line)
		}
	}
	return []byte(strings.Join(result, "\n"))
}
