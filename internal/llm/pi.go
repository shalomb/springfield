package llm

import (
	"context"
	"os/exec"
)

// PiLLM implements LLMClient by calling the 'pi' CLI.
type PiLLM struct {
	Model    string
	executor func(ctx context.Context, name string, arg ...string) ([]byte, error)
}

func (p *PiLLM) Chat(ctx context.Context, messages []Message) (Response, error) {
	args := []string{"-p", "--no-tools"}

	if p.Model != "" {
		args = append(args, "--model", p.Model)
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

	out, err := execFn(ctx, "pi", args...)
	if err != nil {
		return Response{}, err
	}

	// For now, pi CLI doesn't return token usage, so we'll leave it at zero.
	return Response{Content: string(out)}, nil
}

// executorWithFallback tries to run 'pi' directly, then falls back to 'npm exec'.
// This ensures the system works in environments where pi isn't in the PATH.
func (p *PiLLM) executorWithFallback(ctx context.Context, name string, arg ...string) ([]byte, error) {
	// Try 'pi' directly first
	cmd := exec.CommandContext(ctx, name, arg...)
	out, err := cmd.Output()
	if err == nil {
		return out, nil
	}

	// Check if error is "command not found" (pi not in PATH)
	// If so, fall back to npm exec
	if isCommandNotFound(err) {
		// Try 'npm exec' as fallback
		// 'npm exec @mariozechner/pi-coding-agent -- <args>'
		npmArgs := []string{"exec", "@mariozechner/pi-coding-agent", "--"}
		npmArgs = append(npmArgs, arg...)
		cmd = exec.CommandContext(ctx, "npm", npmArgs...)
		out, npmErr := cmd.Output()
		if npmErr == nil {
			return out, nil
		}
		// If npm also fails, return npm error
		return nil, npmErr
	}

	// If pi failed for reasons other than "not found", return that error
	return nil, err
}

// isCommandNotFound checks if an error is due to a command not being found.
func isCommandNotFound(err error) bool {
	if err == nil {
		return false
	}
	// Check for standard "executable file not found" message
	errMsg := err.Error()
	return errMsg == "executable file not found in $PATH" || 
	       errMsg == "command not found"
}
