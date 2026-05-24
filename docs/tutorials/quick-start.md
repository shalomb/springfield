# Quick Start Guide

A rapid walkthrough of common Springfield workflows. For deeper understanding, see [Getting Started](docs/how-to/getting-started.md) and [Standard Workflows](docs/how-to/workflows.md).

## Installation & Setup

First time? Initialize Springfield:
```bash
just install-tools  # Install Go development tools
just build          # Build the Springfield binary
springfield init    # Set up td database, config, and .gitignore
```

## Run Tests

```bash
just test                # Run full graduated test ladder
just test-structure      # Format and vet checks
just test-lint           # Linter checks
just test-unit           # Unit tests (fast)
just test-integration    # BDD integration tests
```

## Start a Feature

Create and build a new feature:
```bash
just start-feature my-awesome-feature  # Create a branch
# ... write feature.md ...
just lisa "Analyze feature.md and create a TODO.md"
just do                                 # Run: Lisa → Ralph → Bart (loop)
```

Or build manually:
```bash
just ralph   # Implement one task from TODO.md
just test    # Verify it works
```

When ready to ship:
```bash
just lovejoy  # Release ceremony (merge, cleanup, celebrate)
```

## Common Commands

| Command | Purpose |
| --- | --- |
| `just build` | Compile the binary |
| `just run` | Build and run locally |
| `just clean` | Remove build artifacts |
| `just lisa [task]` | Run planner agent |
| `just ralph [task]` | Run build agent |
| `just bart [task]` | Run quality review agent |
| `just lovejoy` | Run release ceremony |
| `just do` | Autonomous loop (Lisa → Ralph → Bart) |

## Troubleshooting

**"td not found"**
```bash
# td is a dependency for Springfield's task management system
# Install with: go install github.com/shalomb/td@latest
# Then re-run: springfield init
```

**"Tests failing"**
- Run `just test-unit` to isolate the issue
- Check `TODO.md` and `FEEDBACK.md` for context
- Fix and commit atomically

**"Agent seems stuck"**
- Kill the process (Ctrl+C)
- Review the failing task in `TODO.md`
- Clarify ambiguous instructions and retry

## Next Steps

- **Deep dive:** Read [Standard Workflows](docs/how-to/workflows.md)
- **Understand concepts:** See [Master Model](docs/concepts/model.md)
- **Reference agents:** Check [Agent Reference](docs/reference/agents.md)
- **Troubleshoot:** Review [Debugging Guide](docs/how-to/debugging-and-observability.md)
