# =============================================================================
# DEFAULT RECIPE
# =============================================================================
default:
	@just list

# =============================================================================
# VARIABLES
# =============================================================================
BINARY_NAME := "springfield"
BUILD_DIR := "bin"
TEST_DIR := "tests"

# =============================================================================
# HELP & DOCUMENTATION
# =============================================================================

# Show grouped list of all recipes by function/area
list:
	@echo "🌸 Springfield - AI Agent Orchestration (Go Migration)"
	@echo "==================================================="
	@echo ""
	@echo "🚀 CORE COMMANDS:"
	@echo "  build                Build the application (bin/springfield)"
	@echo "  run                  Build and run the application"
	@echo "  clean                Clean build artifacts"
	@echo ""
	@echo "🧪 GRADUATED TEST LADDER (Run in order, stop on failure):"
	@echo "  test                 Run full graduated test ladder (structure -> lint -> unit -> integration)"
	@echo "  test-structure       Check code structure (fmt, vet)"
	@echo "  test-lint            Run linters (golangci-lint)"
	@echo "  test-unit            Run fast unit tests"
	@echo "  test-integration     Run integration tests (BDD)"
	@echo "  test-coverage        Generate coverage report"
	@echo ""
	@echo "🔧 DEVELOPMENT:"
	@echo "  deps                 Install dependencies (go mod download)"
	@echo "  fmt                  Format code (go fmt)"
	@echo "  lint                 Run linter (golangci-lint)"
	@echo "  install-tools        Install development tools (godog, golangci-lint)"
	@echo ""
	@echo "📚 AGENT MANAGEMENT:"
	@echo "  install              Install agent skills to ~/.pi/agent/skills/"
	@echo "  lisa                 Run the intelligent Planner agent (Lisa)"
	@echo "  ralph                Run the autonomous Build agent (Ralph)"
	@echo ""

# =============================================================================
# BUILD & RUN
# =============================================================================

# Install dependencies
deps:
	go mod download
	go mod tidy

# Build the application
build:
	@echo "Building {{BINARY_NAME}}..."
	@mkdir -p {{BUILD_DIR}}
	go build -o {{BUILD_DIR}}/{{BINARY_NAME}} ./cmd/springfield

# Build and run the application
run *args: build
	@./{{BUILD_DIR}}/{{BINARY_NAME}} {{args}}

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	go clean

# =============================================================================
# GRADUATED TEST LADDER (Fail-Fast Strategy)
# =============================================================================
# The `just test` target executes phases in strict order, stopping on first failure:
#   Phase 1 (test-structure):   Go syntax validation (fmt, vet) - INSTANT
#   Phase 2 (test-lint):        Code quality (golangci-lint) - FAST
#   Phase 3 (test-unit):        Unit tests (short) - FAST
#   Phase 4 (test-integration): Integration tests (BDD) - SLOW
# =============================================================================

test:
	@echo "🚀 Starting Graduated Test Ladder..."
	@just test-structure
	@just test-lint
	@just test-unit
	@just test-integration
	@echo "✅ COMPLETE: All test levels passed!"

# Phase 1: Validate Go structure and syntax (instant fail-fast)
test-structure:
	@echo "🔍 Validating Go structure (Phase 1)..."
	go fmt ./...
	go vet ./...
	@echo "✅ Structure validation passed"

# Phase 2: Check code quality with golangci-lint (fail-fast on quality issues)
test-lint:
	@echo "🔍 Checking code quality (Phase 2)..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "✅ Lint check passed"; \
	else \
		echo "⚠️  golangci-lint not found. Skipping lint check."; \
		echo "   Run 'just install-tools' to install."; \
	fi

# Phase 3: Run fast unit tests (short mode)
test-unit:
	@echo "🧪 Running unit tests (Phase 3)..."
	go test -v -short -race ./internal/... ./pkg/...
	@echo "✅ Unit tests passed"

# Phase 4: Run integration tests (BDD/Godog)
test-integration:
	@echo "🧪 Running integration tests (Phase 4)..."
	@if [ -n "$(shell axon --version 2>/dev/null)" ] || [ -e ~/shalomb/tide/go-axon/bin/axon ] || [ -e ~/shalomb/axon/go-axon/bin/axon ]; then \
		go test -v ./tests/integration; \
		echo "✅ Integration tests passed"; \
	else \
		echo "⚠️  axon binary not found. Skipping integration tests."; \
		echo "   Run 'just install-tools' or install axon locally."; \
	fi

# Run only Axon-related integration tests
test-integration-axon:
	@echo "🧪 Running Axon integration tests..."
	go test -v ./tests/integration --args --godog.tags='@requires_axon && ~@wip'
	@echo "✅ Axon integration tests complete"

# Generate coverage report
test-coverage:
	@echo "📊 Generating coverage report..."
	mkdir -p coverage
	go test -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@echo "Coverage report generated: coverage/coverage.html"

# =============================================================================
# DEVELOPMENT TOOLS
# =============================================================================

# Install development tools (godog, golangci-lint)
install-tools:
	@echo "Installing development tools..."
	go install github.com/cucumber/godog/cmd/godog@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Tools installed to $(go env GOPATH)/bin"

# =============================================================================
# AGENT SKILLS
# =============================================================================

# Install the agent skills to the local pi environment
install:
	@echo "Installing Springfield agents..."
	@mkdir -p ~/.pi/agent/skills
	@cp -r .pi/agent/skills/* ~/.pi/agent/skills/
	@echo "Installed: bart, lisa, lovejoy, marge, ralph to ~/.pi/agent/skills/"

# =============================================================================
# AGENT WORKFLOWS
# =============================================================================

# Run a Springfield agent
agent name task:
    ./bin/springfield --agent {{name}} --task "{{task}}"

# Run the autonomous "Ralph" loop (Agent working off TODO.md)
ralph *args='':
    #!/usr/bin/env bash
    set -euo pipefail
    printf "🤖 Starting Ralph Loop...\n"
    EXTRA_PROMPT="{{args}}"
    while :; do
        if [[ ! -e TODO.md ]]; then
            printf "✅ No TODO.md found. Work complete!\n"
            break
        fi

        printf "📋 Tasks found. Engaging Agent...\n"
        # Note: This assumes npm and the agent package are available
        npm exec @mariozechner/pi-coding-agent -- \
            --verbose \
            --mode json \
            --no-session \
            --thinking medium \
            --no-extensions \
            --provider google-gemini-cli \
            --model gemini-3-flash-preview \
            @TODO.md \
            -p "If there no tasks to work on - git rm TODO.md and create a completion git commit. If there are tasks to work on - assume the role of .github/agents/ralph.md and pick the highest priority task and work on it. Strictly adhere to the Atomic Commit Protocol (docs/standards/atomic-commit-protocol.md). Employ TDD processes (RED -> GREEN -> REFACTOR) and ensure that every commit is an indivisible unit containing BDD specs, TDD tests, minimal implementation, and documentation. Ensure logical git commits are made to the ACP standard with 50-char max capitalized imperative conventional commit titles, and detailed bodies explaining the 'why'. Ensure that the codebase is in a working state after each commit. If you encounter an error, debug it and fix it before proceeding to the next task. ${EXTRA_PROMPT:+USER INSTRUCTION: $EXTRA_PROMPT}"

        echo ''
        echo '********'
        sleep 1
    done

# Run the intelligent "Lisa" loop (Planner preparing work for Ralph)
lisa *args='':
    #!/usr/bin/env bash
    set -euo pipefail
    printf "📚 Starting Lisa Simpson (Intelligent Planner)...\n"
    EXTRA_PROMPT="{{args}}"

    npm exec @mariozechner/pi-coding-agent -- \
        --verbose \
        --mode json \
        --no-session \
        --thinking medium \
        --no-extensions \
        --provider google-gemini-cli \
        --model gemini-3-flash-preview \
        -p "Assume the role of Lisa Simpson (.github/agents/lisa.md). Your mission is to translate high-level intent from PLAN.md into executable tasks for Ralph. \
        1. Reflect & Learn: Analyze recent commits and branch state. Identify learnings, technical debt, or necessary reprioritizations. Update PLAN.md with a 'Retrospective' section for the completed epic if appropriate. \
        2. Technical Breakdown: Identify the next high-priority Epic from PLAN.md. Translate it into a technical breakdown in a new TODO.md. Ensure tasks follow the Atomic Commit Protocol (docs/standards/atomic-commit-protocol.md) - each task should ideally map to one or more atomic commits. \
        3. Moral Compass: Ensure the plan adheres to Takeda's compliance and safety standards (ADR-0001 Building Blocks, RBAC, audit logging). \
        4. Autonomous Setup: Create a new git branch for the epic named 'feat/epic-{name}'. Add the TODO.md and updated PLAN.md to this branch. \
        5. Atomic Handover: Commit the plan with a clear message following ACP standards. \
        You are the intelligent pre-processor. You provide the logic Ralph needs to succeed without eating the paste. Ensure TODO.md tasks are atomic, testable, and include success criteria. ${EXTRA_PROMPT:+USER INSTRUCTION: $EXTRA_PROMPT}"
