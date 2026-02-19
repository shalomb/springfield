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

# Common arguments for the pi-coding-agent (defined as a list for proper expansion)
PI_AGENT_FLAGS := "--verbose --mode json --no-session --thinking medium --no-extensions --provider google-gemini-cli --model gemini-3-flash-preview"

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
	@echo "  herb                 Run the Quality Review agent (Herb)"
	@echo "  bart                 Run the Quality Verification agent (Bart)"
	@echo "  lovejoy              Run the Release agent (Lovejoy)"
	@echo ""
	@echo "🤖 AUTONOMOUS LOOP:"
	@echo "  plan                 Alias for 'just lisa'"
	@echo "  do                   Run the autonomous loop (Lisa -> Ralph -> Bart -> Herb)"
	@echo ""
	@echo "🌿 GIT WORKFLOW:"
	@echo "  start-feature name   Start a new feature branch from main"
	@echo "  start-fix name       Start a new fix branch from main"
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

test:
	@echo "🚀 Starting Graduated Test Ladder..."
	@just test-structure
	@just test-lint
	@just test-unit
	@just test-integration
	@echo "✅ COMPLETE: All test levels passed!"

# Phase 1: Validate Go structure and syntax
test-structure:
	@echo "🔍 Validating Go structure (Phase 1)..."
	go fmt ./...
	go vet ./...
	@echo "✅ Structure validation passed"

# Phase 2: Check code quality
test-lint:
	@echo "🔍 Checking code quality (Phase 2)..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "✅ Lint check passed"; \
	else \
		echo "⚠️  golangci-lint not found. Skipping lint check."; \
		echo "   Run 'just install-tools' to install."; \
	fi

# Phase 3: Run fast unit tests
test-unit:
	@echo "🧪 Running unit tests (Phase 3)..."
	go test -v -short -race ./internal/... ./pkg/...
	@echo "✅ Unit tests passed"

# Phase 4: Run integration tests
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

install-tools:
	@echo "Installing development tools..."
	go install github.com/cucumber/godog/cmd/godog@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Tools installed to $(go env GOPATH)/bin"

# =============================================================================
# AGENT SKILLS
# =============================================================================

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
	#!/usr/bin/env bash
	set -euo pipefail
	./bin/springfield --agent {{name}} --task "{{task}}"

# Run the autonomous "Ralph" loop (Agent working off TODO.md)
ralph *args='':
	#!/usr/bin/env bash
	set -euo pipefail
	printf "🤖 Starting Ralph Loop...\n"
	EXTRA_PROMPT="{{args}}"
	while :; do
		UNCOMMITTED=$(git status --porcelain --untracked-files=no)
		if [[ ! -e TODO.md ]] && [[ -z "$UNCOMMITTED" ]]; then
			printf "✅ No TODO.md found and no uncommitted changes. Work complete!\n"
			break
		fi

		if [[ -n "$UNCOMMITTED" ]]; then
			printf "📝 Uncommitted changes detected. Engaging Ralph to finalize...\n"
		else
			printf "📋 Tasks found in TODO.md. Engaging Ralph...\n"
		fi

		# Note: This assumes npm and the agent package are available
		ARGS=()
		if [[ -e TODO.md ]]; then
			ARGS+=("@TODO.md")
		fi

		npm exec @mariozechner/pi-coding-agent -- {{PI_AGENT_FLAGS}} \
			"${ARGS[@]}" \
			-p "Assume the role of .github/agents/ralph.md. 
			If TODO.md exists, pick the highest priority task and work on it. 
			If there are uncommitted changes but no tasks left in TODO.md, create a clean completion git commit and 'git rm TODO.md' if it still exists.
			Strictly adhere to the Atomic Commit Protocol (docs/standards/atomic-commit-protocol.md). 
			Employ TDD processes (RED -> GREEN -> REFACTOR) and ensure that every commit is an indivisible unit containing BDD specs, TDD tests, minimal implementation, and documentation. 
			Ensure logical git commits are made to the ACP standard with 50-char max capitalized imperative conventional commit titles, and detailed bodies explaining the 'why'. 
			Ensure that the codebase is in a working state after each commit. 
			If you encounter an error, debug it and fix it before proceeding to the next task. ${EXTRA_PROMPT:+USER INSTRUCTION: $EXTRA_PROMPT}"

		echo ''
		echo '********'
		sleep 1
	done

# Alias for lisa
plan *args='':
	@just lisa {{args}}

# Run the autonomous "just do" orchestrator
do *args='':
	@just plan {{args}}
	@just ralph {{args}}
	@just bart {{args}}
	@just herb {{args}}

# Run the intelligent "Lisa" loop (Planner preparing work for Ralph)
lisa *args='':
	#!/usr/bin/env bash
	set -euo pipefail
	printf "📚 Starting Lisa Simpson (Intelligent Planner)...\n"
	EXTRA_PROMPT="{{args}}"

	ARGS=()
	if [[ -e FEEDBACK.md ]]; then
		ARGS+=("@FEEDBACK.md")
	fi

	npm exec @mariozechner/pi-coding-agent -- {{PI_AGENT_FLAGS}} \
		"${ARGS[@]}" \
		-p "Assume the role of Lisa Simpson (.github/agents/lisa.md). Your mission is to translate high-level intent from PLAN.md into executable tasks for Ralph. 
		1. Reflect & Learn: Analyze recent commits and branch state. Identify learnings, technical debt, or necessary reprioritizations. Update PLAN.md with a 'Retrospective' section for the completed epic if appropriate. 
		2. Analyze Feedback: If FEEDBACK.md exists, analyze it against PLAN.md. If errors are critical (breaking functionality, security, crash), create specific corrective tasks in TODO.md. If errors are minor (style, non-blocking edge cases), log them in PLAN.md under 'Known Issues' and clear FEEDBACK.md. DO NOT loop if you have already tried to fix this twice.
		3. Technical Breakdown: Identify the next high-priority Epic from PLAN.md. Translate it into a technical breakdown in a new TODO.md. Ensure tasks follow the Atomic Commit Protocol (docs/standards/atomic-commit-protocol.md) - each task should ideally map to one or more atomic commits. 
		4. Moral Compass: Ensure the plan adheres to Enterprise compliance and safety standards (ADR-000 Building Blocks, RBAC, audit logging). 
		5. Autonomous Setup: Detect the current branch. If on 'main', create a new git branch for the epic named 'feat/epic-{name}'. Add the TODO.md and updated PLAN.md to this branch. 
		6. Atomic Handover: Commit the plan with a clear message following ACP standards. 
		You are the intelligent pre-processor. You provide the logic Ralph needs to succeed without eating the paste. Ensure TODO.md tasks are atomic, testable, and include success criteria. ${EXTRA_PROMPT:+USER INSTRUCTION: $EXTRA_PROMPT}"

# Run the Quality Review loop (Herb)
herb *args='':
	#!/usr/bin/env bash
	set -euo pipefail
	printf "🧐 Starting Herb Powell (Quality Review)...\n"
	EXTRA_PROMPT="{{args}}"

	npm exec @mariozechner/pi-coding-agent -- {{PI_AGENT_FLAGS}} \
		-p "Assume the role of Herb Powell (Quality/Review). Your mission is to perform a strict code review of the recent changes in this branch. 
		1. Static Analysis: Review the code for SOLID principles, Clean Code standards, and Go/Python best practices. 
		2. Atomic Commit Check: Verify that commits follow the Atomic Commit Protocol. 
		3. Feedback: Document any issues, technical debt, or refactoring suggestions in FEEDBACK.md. If the code is excellent, state that it is ready for verification. Exit with a non-zero status if critical issues or ACP violations are found. ${EXTRA_PROMPT:+USER INSTRUCTION: $EXTRA_PROMPT}"

# Run the Quality Verification loop (Bart)
bart *args='':
	#!/usr/bin/env bash
	set -euo pipefail
	printf "🛹 Starting Bart Simpson (Quality Verification)...\n"
	EXTRA_PROMPT="{{args}}"

	npm exec @mariozechner/pi-coding-agent -- {{PI_AGENT_FLAGS}} \
		-p "Assume the role of Bart Simpson (Quality/Verification). Your mission is to break the code and find bugs. 
		1. Test Execution: Run 'just test' to verify the entire test ladder. 
		2. BDD Validation: Verify that the implemented code matches the Gherkin scenarios in docs/features/. 
		3. Adversarial Testing: Think of edge cases Ralph might have missed. 
		4. Feedback: Document all failures, bugs, or missing coverage in FEEDBACK.md. Flag critical issues that block release. Exit with a non-zero status if any test fails or critical bugs are discovered. ${EXTRA_PROMPT:+USER INSTRUCTION: $EXTRA_PROMPT}"

# Run the Release Ceremony loop (Lovejoy)
lovejoy *args='':
	#!/usr/bin/env bash
	set -euo pipefail
	printf "⛪ Starting Reverend Lovejoy (Release Ceremony)...\n"
	EXTRA_PROMPT="{{args}}"

	npm exec @mariozechner/pi-coding-agent -- {{PI_AGENT_FLAGS}} \
		-p "Assume the role of Reverend Lovejoy (Release). Your mission is to perform the release ceremony. 
		1. Readiness Check: Verify that TODO.md is empty and FEEDBACK.md contains no blocking issues. 
		2. Merge: Merge the feature branch into main using a squash merge with a clean, descriptive message. 
		3. Documentation: Update CHANGELOG.md and capture any major learnings for the next cycle. 
		4. Cleanup: Delete the local and remote feature branch after a successful merge. ${EXTRA_PROMPT:+USER INSTRUCTION: $EXTRA_PROMPT}"

# =============================================================================
# GIT WORKFLOW
# =============================================================================

# Start a new feature branch from main
start-feature name:
	git checkout main
	git pull
	git checkout -b feat/{{name}}

# Start a new fix branch from main
start-fix name:
	git checkout main
	git pull
	git checkout -b fix/{{name}}
