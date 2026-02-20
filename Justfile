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
PI_FLAGS := "--verbose --mode json --no-session --thinking medium --no-extensions --provider google-gemini-cli --model gemini-3-flash-preview"

# Prompts defined as variables to avoid quoting hell in recipes
PROMPT_RALPH := "Assume the role of .github/agents/ralph.md. If TODO.md exists, pick the highest priority task and work on it. If there are uncommitted changes but no tasks left in TODO.md, create a clean completion git commit and 'git rm TODO.md' if it still exists. Strictly adhere to the Atomic Commit Protocol (docs/standards/atomic-commit-protocol.md). Employ TDD processes (RED -> GREEN -> REFACTOR) and ensure that every commit is an indivisible unit containing BDD specs, TDD tests, minimal implementation, and documentation. Ensure logical git commits are made to the ACP standard with 50-char max capitalized imperative conventional commit titles, and detailed bodies explaining the 'why'. Ensure that the codebase is in a working state after each commit. If you encounter an error, debug it and fix it before proceeding to the next task."

PROMPT_LISA := "Assume the role of Lisa Simpson (.github/agents/lisa.md). Your mission is to translate high-level intent from PLAN.md into executable tasks for Ralph. 1. Reflect & Learn: Analyze recent commits and branch state. Identify learnings, technical debt, or necessary reprioritizations. Update PLAN.md with a 'Retrospective' section for the completed epic if appropriate. 2. Analyze Feedback: If FEEDBACK.md exists, analyze it against PLAN.md. If errors are critical (breaking functionality, security, crash), create specific corrective tasks in TODO.md. If errors are minor (style, non-blocking edge cases), log them in PLAN.md under 'Known Issues' and clear FEEDBACK.md. DO NOT loop if you have already tried to fix this twice. 3. Technical Breakdown: Identify the next high-priority Epic from PLAN.md. Translate it into a technical breakdown in a new TODO.md. Ensure tasks follow the Atomic Commit Protocol (docs/standards/atomic-commit-protocol.md) - each task should ideally map to one or more atomic commits. 4. Moral Compass: Ensure the plan adheres to Enterprise compliance and safety standards (ADR-000 Building Blocks, RBAC, audit logging). 5. Autonomous Setup: Detect the current branch. If on 'main', create a new git branch for the epic named 'feat/epic-{name}'. Add the TODO.md and updated PLAN.md to this branch. 6. Atomic Handover: Commit the plan with a clear message following ACP standards. You are the intelligent pre-processor. You provide the logic Ralph needs to succeed without eating the paste. Ensure TODO.md tasks are atomic, testable, and include success criteria."

PROMPT_BART := "Assume the role of Bart Simpson (Quality Agent). Your mission is to verify and break the code. 1. Static Review: Review the code for SOLID principles, Clean Code standards, Go best practices, and Atomic Commit Protocol adherence. 2. Dynamic Verification: Run 'just test' to verify the test ladder and BDD scenarios. 3. Adversarial Testing: Think of edge cases Ralph might have missed. 4. Feedback: Document all static issues, test failures, bugs, or missing coverage in FEEDBACK.md. Flag critical issues that block release. Exit with a non-zero status if any test fails or critical bugs are discovered."

PROMPT_LOVEJOY := "Assume the role of Reverend Lovejoy (Release). Your mission is to perform the release ceremony. 1. Readiness Check: Verify that TODO.md is empty and FEEDBACK.md contains no blocking issues. 2. Merge: Merge the feature branch into main using a squash merge with a clean, descriptive message. 3. Documentation: Update CHANGELOG.md and capture any major learnings for the next cycle. 4. Cleanup: Delete the local and remote feature branch after a successful merge."

# =============================================================================
# DEFAULT
# =============================================================================
default:
    @just list

# =============================================================================
# HELP
# =============================================================================
list:
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

    # Export prompt to avoid shell escaping issues
    export PROMPT="{{PROMPT_RALPH}}"
    export EXTRA_INSTRUCTION="{{args}}"

    while :; do
        # Use git status --porcelain for reliable parsing
        if [[ ! -e TODO.md ]] && [[ -z "$(git status --porcelain --untracked-files=no)" ]]; then
            printf "✅ No TODO.md found and no uncommitted changes. Work complete!\n"
            break
        fi

        if [[ -n "$(git status --porcelain --untracked-files=no)" ]]; then
            printf "📝 Uncommitted changes detected. Engaging Ralph to finalize...\n"
        else
            printf "📋 Tasks found in TODO.md. Engaging Ralph...\n"
        fi

        # Build args array safely
        cmd_args=({{PI_FLAGS}})
        if [[ -e TODO.md ]]; then
            cmd_args+=("@TODO.md")
        fi

        # Append user instruction if present
        full_prompt="$PROMPT"
        if [[ -n "$EXTRA_INSTRUCTION" ]]; then
            full_prompt="${full_prompt} USER INSTRUCTION: $EXTRA_INSTRUCTION"
        fi

        # Execute
        npm exec @mariozechner/pi-coding-agent -- "${cmd_args[@]}" -p "$full_prompt"

        printf '\n********\n'
        sleep 1
    done

# Lisa: The Planner
lisa *args:
    #!/usr/bin/env bash
    set -euo pipefail
    printf "📚 Starting Lisa Simpson (Intelligent Planner)...\n"

    export PROMPT="{{PROMPT_LISA}}"
    export EXTRA_INSTRUCTION="{{args}}"

    cmd_args=({{PI_FLAGS}})
    if [[ -e FEEDBACK.md ]]; then
        cmd_args+=("@FEEDBACK.md")
    fi

    full_prompt="$PROMPT"
    if [[ -n "$EXTRA_INSTRUCTION" ]]; then
        full_prompt="${full_prompt} USER INSTRUCTION: $EXTRA_INSTRUCTION"
    fi

    npm exec @mariozechner/pi-coding-agent -- "${cmd_args[@]}" -p "$full_prompt"

# Aliases
plan *args:
    @just lisa {{args}}

# Orchestrator
do *args:
    #!/usr/bin/env bash
    set -euo pipefail

    # Greg says: Don't use uppercase for private vars
    max_iterations=2
    iteration=1

    while (( iteration <= max_iterations )); do
        printf "\n🚀 [Iteration %d/%d] Starting Autonomous Loop...\n" "$iteration" "$max_iterations"

        just plan {{args}}
        just ralph {{args}}
        just bart {{args}}

        if [[ ! -e FEEDBACK.md ]]; then
            printf "\n✅ No FEEDBACK.md found. Cycle complete!\n"
            break
        fi

        # Grep quiet, case-insensitive, extended regex
        if grep -qiE "critical|blocker|rejected|fail" FEEDBACK.md; then
            printf "\n⚠️ Critical issues found in FEEDBACK.md. Re-looping for corrective planning...\n"
            ((iteration++))
        else
            printf "\n✅ Only minor issues found in FEEDBACK.md. Ending loop.\n"
            break
        fi
    done

    if (( iteration > max_iterations )); then
        printf "\n🛑 [HARD STOP] Loop limit reached (%d iterations). Human review required!\n" "$max_iterations"
        exit 1
    fi

# Reviewers
bart *args:
    #!/usr/bin/env bash
    set -euo pipefail
    printf "🛹 Starting Bart Simpson (Quality Agent)...\n"

    full_prompt="{{PROMPT_BART}}"
    if [[ -n "{{args}}" ]]; then
        full_prompt="${full_prompt} USER INSTRUCTION: {{args}}"
    fi

    npm exec @mariozechner/pi-coding-agent -- "{{PI_FLAGS}}" -p "$full_prompt"

    # Post-Execution Assertion: Fail if Bart found critical issues
    if [[ -f FEEDBACK.md ]] && grep -qiE "critical|blocker|rejected|fail" FEEDBACK.md; then
        printf "❌ Bart found critical issues in FEEDBACK.md. Exiting with error.\n" >&2
        exit 1
    fi

lovejoy *args:
    #!/usr/bin/env bash
    set -euo pipefail
    printf "⛪ Starting Reverend Lovejoy (Release Ceremony)...\n"

    full_prompt="{{PROMPT_LOVEJOY}}"
    if [[ -n "{{args}}" ]]; then
        full_prompt="${full_prompt} USER INSTRUCTION: {{args}}"
    fi

    npm exec @mariozechner/pi-coding-agent -- "{{PI_FLAGS}}" -p "$full_prompt"

    # Post-Execution Assertion: Fail if release blocked
    if [[ -f TODO.md ]]; then
        printf "❌ TODO.md is not empty. Cannot release.\n" >&2
        exit 1
    fi
    if [[ -f FEEDBACK.md ]] && grep -qiE "critical|blocker|rejected|fail" FEEDBACK.md; then
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
