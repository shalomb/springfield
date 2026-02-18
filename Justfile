# Springfield Protocol Justfile

# Launch the agent orchestration mesh
flow:
    @mkdir -p logs
    @bash scripts/tmux-orch.sh

# Attach to the current orchestration session
attach:
    @bash scripts/tmux-orch.sh attach

# Clean up agents and tmux windows
stop:
    @bash scripts/tmux-orch.sh stop

# Tail the combined structured log stream
logs:
    @mkdir -p logs
    @touch logs/springfield.log
    @tail -f logs/springfield.log

# Run all tests (Unit + BDD)
test:
    @PYTHONPATH=. pytest tests/unit tests/integration

# Run unit tests with coverage report
test-cov:
    @PYTHONPATH=. pytest --cov=scripts --cov-report=term-missing tests/unit

# Git Branching Commands
start-feature name:
    @if [ -z "$(echo {{name}} | grep -E '^[a-z0-9-]+$')" ]; then \
        echo "Error: Branch name must be lowercase-kebab-case"; \
        exit 1; \
    fi
    @git checkout -b feat/{{name}}
    @echo "Switched to branch feat/{{name}}"

start-fix name:
    @if [ -z "$(echo {{name}} | grep -E '^[a-z0-9-]+$')" ]; then \
        echo "Error: Branch name must be lowercase-kebab-case"; \
        exit 1; \
    fi
    @git checkout -b fix/{{name}}
    @echo "Switched to branch fix/{{name}}"
