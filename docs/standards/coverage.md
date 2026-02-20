# Code Coverage Standards

This document outlines the code coverage expectations and documented exceptions for the Springfield project.

## Targets

Springfield aims for **95%+ statement coverage** across all internal packages.

## Exceptions

The following areas are intentionally excluded from coverage requirements:

### 1. Main Entry Points

- `cmd/springfield/main.go`: The `main()` function is a simple wrapper around `runMain()` and `os.Exit()`. It is not unit tested to avoid complex process-level side effects during the test suite. `runMain()` is tested via `TestRunMain`.

### 2. External CLI Wrappers (Default Branches)

- `internal/llm/pi.go`: The default branch of `execFn` in `Chat()` calls the real `pi` CLI. This branch is not covered by unit tests to avoid making real subprocess calls during the test suite. This branch is considered structurally equivalent to a `main()` entry point for the CLI wrapper.

## Verification

Coverage is verified during the CI process and can be checked locally using:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```
