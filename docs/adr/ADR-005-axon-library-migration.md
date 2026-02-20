# ADR-005: Axon Library Migration

**Date:** 2026-02-18
**Status:** Accepted
**Replaces:** ADR-004: Axon Library Integration Strategy (Phase 1/Option A)

## Context
Initial integration between Springfield and Axon (EPIC-004) was performed via a CLI-subprocess bridge. Springfield would invoke the `axon` binary, pass a command, and parse the resulting JSON output.

This approach was a "Phase 1" necessity because Axon's core logic was not yet exported as a Go library. However, a review of the Axon codebase (`~/shalomb/axon/go-axon`) reveals that `pkg/executor` and `pkg/types` have now been successfully extracted and matured.

## Decision
We will migrate Springfield's `internal/sandbox` implementation from a CLI-based execution model to a direct Go library integration with Axon.

### 1. Implementation Detail
- Springfield will `require` the `github.com/shalomb/axon` module.
- The `AxonSandbox` struct in `internal/sandbox/axon.go` will be refactored to use `executor.New()` and `exec.Execute(ctx, req)`.
- We will use Axon's canonical types from `pkg/types` (e.g., `types.Result`) instead of duplicating them or parsing JSON strings.

### 2. Configuration Integration
- Springfield will pass its configuration into Axon's `executor.New(opts...)` to control:
    - Container runtime (Podman).
    - Security levels and guardrails.
    - Context intelligence settings.

## Consequences

### Positive
- **Performance:** Eliminates the overhead of spawning a shell process to invoke the CLI.
- **Type Safety:** Compile-time verification of the integration contract; no more fragile JSON parsing.
- **Robustness:** Direct access to exit codes, stdout, and stderr buffers without intermediate string contamination.
- **Intelligence:** Native access to Axon's context metadata (project type, git status) without additional parsing logic.

### Negative
- **Dependency Coupling:** Springfield now depends on the Axon module at compile time.
- **Build Complexity:** Springfield's build pipeline must now be able to resolve and pull the Axon Go module.

## Consequences for Ralph (Execution)
- Ralph must refactor `internal/sandbox/axon.go`.
- Ralph must update `go.mod` to include Axon.
- Ralph must update BDD tests to verify that the library integration correctly captures all edge cases (including the exit code issues identified in FEEDBACK.md).
