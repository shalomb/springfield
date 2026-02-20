package sandbox

import (
	"context"

	"github.com/shalomb/axon/pkg/types"
)

// Sandbox defines the interface for executing commands in an isolated environment.
type Sandbox interface {
	Execute(ctx context.Context, command string) (*types.Result, error)
}
