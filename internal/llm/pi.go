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
		// Default executor calls real 'pi' CLI. Not tested in unit tests to
		// avoid subprocess calls.
		execFn = func(ctx context.Context, name string, arg ...string) ([]byte, error) {
			return exec.CommandContext(ctx, name, arg...).Output()
		}
	}

	out, err := execFn(ctx, "pi", args...)
	if err != nil {
		return Response{}, err
	}

	// For now, pi CLI doesn't return token usage, so we'll leave it at zero.
	return Response{Content: string(out)}, nil
}
