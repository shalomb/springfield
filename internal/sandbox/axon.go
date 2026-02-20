package sandbox

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shalomb/axon/pkg/executor"
	"github.com/shalomb/axon/pkg/types"
)

// AxonSandbox implements the Sandbox interface using the axon library.
type AxonSandbox struct {
	exec *executor.Executor
}

// NewAxonSandbox creates a new AxonSandbox.
func NewAxonSandbox(configPath string) (*AxonSandbox, error) {
	if configPath == "" {
		// Environment variable discovery
		configPath = os.Getenv("SPRINGFIELD_CONFIG")
	}

	if configPath == "" {
		// Look for config.toml in current directory or parent directories
		dir, err := os.Getwd()
		if err == nil {
			for {
				path := filepath.Join(dir, "config.toml")
				if _, err := os.Stat(path); err == nil {
					configPath = path
					break
				}
				parent := filepath.Dir(dir)
				if parent == dir {
					break
				}
				dir = parent
			}
		}
	}

	var opts []executor.Option
	if configPath != "" {
		opts = append(opts, executor.WithConfigPath(configPath))
	}

	opts = append(opts,
		executor.WithContainerRuntime("podman"), // Springfield prefers podman
		executor.WithBaseImage("docker.io/library/debian:trixie-slim"),
		executor.WithGuardrails(true),
		executor.WithSecurityLevel("development"),
		executor.WithAgent("bash"),
		executor.WithCPULimit("0.5"),     // 50% of one core
		executor.WithMemoryLimit("512m"), // 512MB
	)

	ex, err := executor.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize axon executor: %w", err)
	}

	return &AxonSandbox{exec: ex}, nil
}

// Execute runs a command inside an axon sandbox and returns structured output.
func (s *AxonSandbox) Execute(ctx context.Context, command string) (*types.Result, error) {
	if s.exec == nil {
		return nil, fmt.Errorf("axon executor not initialized")
	}

	req := executor.Request{
		Command: command,
		Agent:   "bash",
	}

	return s.exec.Execute(ctx, req)
}
