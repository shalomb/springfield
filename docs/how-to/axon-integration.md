# Axon Integration

Springfield's sandboxing layer (`internal/sandbox`) is backed by the [Axon](https://github.com/unop/axon) library. This provides a secure, observed, and context-aware execution environment for AI agents.

## Architecture

Springfield uses the `axon` CLI in `--format json` mode to execute commands. This enables:
- **Type Safety**: Structured output with exit codes, stdout, and stderr separated.
- **Context Awareness**: Axon automatically detects project type (Go, JavaScript, etc.), build tools, and git status.
- **Execution Metadata**: Real-time tracking of execution duration and working directories.
- **Tool Discovery**: Automatic detection of available DevOps tools (git, go, npm, etc.).

## Configuration

The `AxonSandbox` struct in `internal/sandbox/axon.go` handles the execution. It resolves the `axon` binary by checking:
1. The provided `AxonPath` (if any).
2. The system `$PATH`.
3. Fallback paths: `~/shalomb/axon/go-axon/bin/axon` and `~/shalomb/tide/go-axon/bin/axon`.

## Running Integration Tests

Integration tests require a live `axon` binary that supports the `--format json` flag.

To run only the Axon integration tests:
```bash
just test-integration-axon
```

If the `axon` binary is not found, these tests will be skipped automatically during the full test suite (`just test-integration`).

## Context Injection

The agent loop (`internal/agent/agent.go`) receives context metadata from each Axon execution. This metadata is fed back into the LLM conversation, allowing the agent to stay aware of the project state without explicitly running diagnostic commands.

Example metadata:
- Project type: go
- Build tool: go
- Git status: main branch, clean
