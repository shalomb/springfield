# =============================================================================
# SETTINGS & CONFIGURATION
# =============================================================================
set shell := ["bash", "-uc"]
set positional-arguments := true
set ignore-comments := true

# =============================================================================
# VARIABLES
# =============================================================================
BINARY_NAME := "springfield"
BUILD_DIR := "bin"

# =============================================================================
# DEFAULT
# =============================================================================
default:
    @just help

# =============================================================================
# HELP
# =============================================================================
help:
    @printf "🌸 Springfield - AI Agent Orchestration (Go Migration)\n"
    @printf "===================================================\n\n"
    @printf "🚀 CORE COMMANDS:\n"
    @printf "  %-20s %s\n" "build" "Build the application (bin/springfield)"
    @printf "  %-20s %s\n" "run" "Build and run the application"
    @printf "  %-20s %s\n" "clean" "Clean build artifacts"
    @printf "\n🧪 GRADUATED TEST LADDER:\n"
    @printf "  %-20s %s\n" "test" "Run full graduated test ladder"
    @printf "  %-20s %s\n" "test-structure" "Go syntax validation (fmt, vet)"
    @printf "  %-20s %s\n" "test-lint" "Code quality (golangci-lint)"
    @printf "  %-20s %s\n" "test-unit" "Fast unit tests"
    @printf "  %-20s %s\n" "test-integration" "Integration tests (BDD)"
    @printf "\n📚 AGENT MANAGEMENT:\n"
    @printf "  %-20s %s\n" "lisa" "Planner Agent"
    @printf "  %-20s %s\n" "ralph" "Build Agent"
    @printf "  %-20s %s\n" "do" "Autonomous Loop (Lisa->Ralph->Bart)"
    @printf "\n🌿 GIT WORKFLOW:\n"
    @printf "  %-20s %s\n" "start-feature" "Start a new feature branch"

# =============================================================================
# BUILD & RUN
# =============================================================================

deps:
    go mod download
    go mod tidy

build:
    @printf "Building %s...\n" "{{BINARY_NAME}}"
    @mkdir -p "{{BUILD_DIR}}"
    go build -o "{{BUILD_DIR}}/{{BINARY_NAME}}" ./cmd/springfield

run *args: build
    @"{{BUILD_DIR}}/{{BINARY_NAME}}" {{args}}

clean:
    @printf "Cleaning build artifacts...\n"
    rm -rf "{{BUILD_DIR}}"
    go clean

# =============================================================================
# TESTING
# =============================================================================

test:
    @printf "🚀 Starting Graduated Test Ladder...\n"
    @just test-structure
    @just test-lint
    @just test-unit
    @just test-integration
    @printf "✅ COMPLETE: All test levels passed!\n"

test-structure:
    @printf "🔍 Validating Go structure (Phase 1)...\n"
    go fmt ./...
    go vet ./...
    @printf "✅ Structure validation passed\n"

test-lint:
    @printf "🔍 Checking code quality (Phase 2)...\n"
    @if command -v golangci-lint &>/dev/null; then \
        golangci-lint run ./...; \
        printf "✅ Lint check passed\n"; \
    else \
        printf "⚠️  golangci-lint not found. Skipping.\n"; \
    fi

test-unit:
    @printf "🧪 Running unit tests (Phase 3)...\n"
    go test -v -short -race ./internal/... ./pkg/...
    @printf "✅ Unit tests passed\n"

test-integration:
    @printf "🧪 Running integration tests (Phase 4)...\n"
    @# Check if podman is available (required by axon library)
    @if command -v podman &>/dev/null; then \
        go test -v ./tests/integration; \
        printf "✅ Integration tests passed\n"; \
    else \
        printf "⚠️  podman not found. Skipping integration tests (Axon library requires podman).\n"; \
    fi

test-coverage:
    @printf "📊 Generating coverage report...\n"
    @mkdir -p coverage
    go test -coverprofile=coverage/coverage.out ./...
    go tool cover -html=coverage/coverage.out -o coverage/coverage.html
    @printf "Coverage report generated: coverage/coverage.html\n"

install-tools:
    @printf "Installing development tools...\n"
    go install github.com/cucumber/godog/cmd/godog@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    @printf "✅ Tools installed to %s/bin\n" "$(go env GOPATH)"

# =============================================================================
# AGENTS
# =============================================================================

install:
    @printf "Installing Springfield agents...\n"
    @mkdir -p ~/.pi/agent/skills
    @cp -r .pi/agent/skills/* ~/.pi/agent/skills/
    @printf "Installed agents to ~/.pi/agent/skills/\n"

# Generic agent runner
agent name task:
    #!/usr/bin/env bash
    set -euo pipefail
    ./bin/springfield --agent "$1" --task "$2"

# Ralph: The Builder
ralph *args:
    #!/usr/bin/env bash
    set -euo pipefail
    printf "🤖 Starting Ralph Loop...\n"

    # Build the task instruction
    task_instruction=""
    if [[ -n "{{args}}" ]]; then
        task_instruction="{{args}}"
    fi

    # Use the Go-based Springfield binary instead of npm/pi-coding-agent
    ./bin/springfield --agent ralph --task "$task_instruction"

# Lisa: The Planner
lisa *args:
    #!/usr/bin/env bash
    set -euo pipefail
    printf "📚 Starting Lisa Simpson (Intelligent Planner)...\n"

    # Build the task instruction
    task_instruction=""
    if [[ -n "{{args}}" ]]; then
        task_instruction="{{args}}"
    fi

    # Use the Go-based Springfield binary instead of npm/pi-coding-agent
    ./bin/springfield --agent lisa --task "$task_instruction"

# Aliases
plan *args:
    @just lisa {{args}}

# Orchestrator
do *args:
    @./bin/springfield orchestrate {{args}}

# Reviewers
bart *args:
    #!/usr/bin/env bash
    set -euo pipefail
    printf "🛹 Starting Bart Simpson (Quality Agent)...\n"

    # Build the task instruction
    task_instruction=""
    if [[ -n "{{args}}" ]]; then
        task_instruction="{{args}}"
    fi

    # Use the Go-based Springfield binary instead of npm/pi-coding-agent
    ./bin/springfield --agent bart --task "$task_instruction"

    # Post-Execution Assertion: Fail if Bart found critical issues
    if [[ -f FEEDBACK.md ]] && grep -qE "Status:.*(REJECTED|BLOCKED)|❌.*Verdict" FEEDBACK.md; then
        printf "❌ Bart found critical issues in FEEDBACK.md. Exiting with error.\n" >&2
        exit 1
    fi

lovejoy *args:
    #!/usr/bin/env bash
    set -euo pipefail
    printf "⛪ Starting Reverend Lovejoy (Release Ceremony)...\n"

    # Build the task instruction
    task_instruction=""
    if [[ -n "{{args}}" ]]; then
        task_instruction="{{args}}"
    fi

    # Use the Go-based Springfield binary instead of npm/pi-coding-agent
    ./bin/springfield --agent lovejoy --task "$task_instruction"

    # Post-Execution Assertion: Fail if release blocked
    if [[ -f TODO.md ]]; then
        printf "❌ TODO.md is not empty. Cannot release.\n" >&2
        exit 1
    fi
    if [[ -f FEEDBACK.md ]] && grep -qE "Status:.*(REJECTED|BLOCKED)|❌.*Verdict" FEEDBACK.md; then
        printf "❌ FEEDBACK.md contains blocking issues. Cannot release.\n" >&2
        exit 1
    fi

# =============================================================================
# GIT
# =============================================================================

start-feature name:
    #!/usr/bin/env bash
    set -euo pipefail

    # Greg: Check input first. Lhunath: Quote variables.
    if [[ ! "$1" =~ ^[a-z0-9-]+$ ]]; then
        printf >&2 "Error: Branch name '%s' must be in lowercase-kebab-case.\n" "$1"
        exit 1
    fi

    git checkout main

    # Check if upstream is configured before pulling
    if git rev-parse --abbrev-ref @{u} &>/dev/null; then
        git pull
    fi

    git checkout -b "feat/$1"

start-fix name:
    #!/usr/bin/env bash
    set -euo pipefail

    if [[ ! "$1" =~ ^[a-z0-9-]+$ ]]; then
        printf >&2 "Error: Branch name '%s' must be in lowercase-kebab-case.\n" "$1"
        exit 1
    fi

    git checkout main
    if git rev-parse --abbrev-ref @{u} &>/dev/null; then
        git pull
    fi

    git checkout -b "fix/$1"
