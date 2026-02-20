# Go Agent Development Standards

## 1. Project Structure

We follow the standard Go project layout:

```text
springfield/
├── cmd/                # Main applications (entry points)
│   └── springfield/    # Main application name
├── internal/           # Private application and library code
│   ├── agent/          # Core agent logic
│   ├── llm/            # LLM integration
│   ├── memory/         # Context and memory management
│   └── sandbox/        # Execution environment (Axon integration)
├── pkg/                # Library code that's ok to use by external apps
└── tests/              # Additional external test apps and test data
```

## 2. Coding Style

- **Formatting:** Always run `go fmt` (or `gofmt -s -w`).
- **Linting:** We use `golangci-lint` with strict settings. Zero warnings policy.
- **Naming:**
    - `MixedCaps` for exported names (public).
    - `mixedCaps` for unexported names (private).
    - Short, meaningful variable names (`i`, `err` are fine in short scopes).
    - Package names should be singular and lowercase (`agent`, not `Agents`).

## 3. Error Handling

- **Explicit Checks:** Check errors immediately.
    ```go
    f, err := os.Open("filename.ext")
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    ```
- **Wrapping:** Use `%w` to wrap errors for context.
- **No Panics:** Avoid `panic()` in libraries. Return errors instead. Only `main` or initialization code should panic if recovery is impossible.

## 4. Concurrency

- **Channels over Shared Memory:** "Do not communicate by sharing memory; instead, share memory by communicating."
- **Context:** Pass `context.Context` as the first argument to functions that involve I/O or long-running operations.
- **Goroutine Leaks:** Ensure every goroutine has a clear termination condition (usually via context cancellation).

## 5. Testing Strategy (The Ladder)

We use a **Graduated Test Ladder** to catch issues early. Run `just test` to execute the full suite.

1.  **Structure:** `go vet`, `go fmt` (Check correctness)
2.  **Lint:** `golangci-lint` (Check quality)
3.  **Unit:** `go test -short ./...` (Check logic, mocked dependencies)
4.  **Integration:** `godog` (Check behavior, real dependencies/sandbox)
5.  **Acceptance:** End-to-end user scenarios.

### Writing Tests
- **Table-Driven Tests:** Preferred for unit tests.
    ```go
    func TestAdd(t *testing.T) {
        tests := []struct {
            name     string
            a, b     int
            expected int
        }{
            {"positive", 1, 2, 3},
            {"negative", -1, -1, -2},
        }
        for _, tt := range tests {
            t.Run(tt.name, func(t *testing.T) {
                assert.Equal(t, tt.expected, Add(tt.a, tt.b))
            })
        }
    }
    ```
- **BDD (Godog):** Feature files describe behavior in Gherkin syntax. Step definitions implement the logic.

## 6. Dependency Injection

- **Interfaces:** Define interfaces where you use them (consumer-side).
- **Structs:** Accept dependencies in the constructor/factory function.
    ```go
    type Agent struct {
        llm LLMClient
    }
    
    func NewAgent(llm LLMClient) *Agent {
        return &Agent{llm: llm}
    }
    ```

## 7. Configuration

- **Environment Variables:** Configuration via `env` vars (12-factor app).
- **Flags:** CLI flags for runtime overrides (using `cobra`/`viper`).
- **Defaults:** Sensible defaults for optional settings.

## 8. Task Management (td)

We use `td(1)` for shared planning state across git worktrees. This ensures a consistent view of the project's live state without branch contention.

- **Discovery:** Run `td usage --new-session` at the start of a session to understand current focus.
- **Progress:** Use `td start {id}` when beginning a task and `td handoff --done` when finished.
- **Queries:** Use `td query "status = ready"` to see what's ready for work.
- **Signals:** Bart uses `td log --decision {type}` to signal Epic outcomes to the orchestrator.
