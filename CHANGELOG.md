# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.0] - 2026-02-20

### Added
- **EPIC-005: Agent Governance & Selection**
  - Unified configuration system via `.springfield.toml` for budget control and model selection.
  - Session-based budget enforcement with hard limits (per-session, per-day, per-request).
  - LLM model fallback logic for graceful degradation under constraints.
  - Token usage and cost logging integration across all LLM calls.

- **Planning Architecture (ADRs 007 & 008)**
  - Fidelity gradient for Epic maturity: far-term stubs → near-term options → ready.
  - Lisa's LRM (Logical Reasoning Model) for option evaluation using Tree-of-Thought and Self-Consistency.
  - Task Decomposition strategies: by workflow step, by business rule, by data variation.
  - State machine for Epic lifecycle with typed transitions (proposed→ready→implemented→done/verified).

- **Quality Standards & Indices**
  - Farley Index: 8 properties for test quality (Fast, Isolated, Repeatable, Self-Verifying, Independent, Focused, Deterministic, Maintainable).
  - Adzic Index: 8 properties for BDD scenario quality (Business-Readable, Intention-Revealing, Atomic, Data-Driven, Executable, Non-Redundant, Focused, Maintainable).
  - Feedback standard with typed signal output (✅ Approved, ⚠️ Rework, ❌ Blocker, ❓ Viability Failure).
  - Shift-left quality gates: Ralph checks Farley in code review, Marge checks Adzic in scenario design.

- **Skills Infrastructure**
  - New agent skills: `impersonate` (find and load agent contexts), `farley-index`, `adzic-index`.
  - Skill mirrors in `.github/skills/` for non-pi-SDK tooling to discover agent capabilities.
  - Enhanced agent definitions: aligned responsibilities with governance (Ralph → Task Decomposition, Lisa → LRM, Marge → Adzic, Bart → Typed Feedback).

- **Documentation & Reference**
  - Comprehensive ADRs with amendment protocol (ADR-007 Amendment A: ADR Lifecycle).
  - Agent responsibility alignment with standards (`docs/standards/AGENTS.md`).
  - Reference guides: `farley-index.md`, `adzic-index.md`, `task-decomposition.md`.
  - Discovery documentation: `sandbox-audit.md` for environment capability evaluation.

- **Testing & Integration**
  - 52 unit tests covering agent loop, LLM fallback, config parsing, logger concurrency.
  - 16 BDD scenarios (Gherkin/Godog) for agent governance, feedback loop, sandboxing.
  - Integration tests for governance policy enforcement and task decomposition.

### Fixed
- Fixed Justfile `PI_FLAGS` quoting to handle empty or whitespace-only environment variables.
- Improved FINISH marker detection robustness (word boundary, line-specific matching).
- Shell redirection guardrails refined: now allow pipes and legitimate shell patterns while blocking exploits.
- Standardized error handling across core modules: no ignored errors on logger.Log, os.Chdir, etc.

### Changed
- Agent definitions reorganized to reflect governance responsibilities:
  - Ralph now owns Task Decomposition strategies (upfront decomposition before implementation).
  - Lisa now owns Logical Reasoning Model for option generation and evaluation.
  - Marge now owns Adzic Index application in Feature Brief design.
  - Bart now generates typed Feedback signals with explicit "viability failure" escalation path.

- Planning architecture evolved:
  - Epic Intent Layer (Marge's Feature Brief) is immutable once decided.
  - Epic Approach Layer (Lisa's LRM decision) is fixed for iteration but immutable in handoff.
  - Task Decomposition (Ralph's TODO-{td-id}.md) is working layer, evolved during implementation.
  - Constraints layer (inherited ADRs) is non-negotiable.

### Deprecated
- Direct string matching on agent prose for orchestration (replaced by Springfield binary with typed state machine).
- Manual TODO.md management (replaced by td(1) and Springfield binary integration).

### See Also
- [ADR-007: Epic Refinement and Lisa's LRM](docs/adr/ADR-007-epic-refinement-and-lisa-lrm.md)
- [ADR-008: Planning State (td) and Springfield Orchestrator](docs/adr/ADR-008-planning-state-td-springfield-orchestrator.md)
- [Farley Index](docs/reference/farley-index.md)
- [Adzic Index](docs/reference/adzic-index.md)
- [Task Decomposition Guide](docs/standards/task-decomposition.md)
- [Feedback Standard](docs/standards/feedback.md)

## [0.3.0] - 2026-02-19

### Added
- **EPIC-008: Knowledge Architecture (Diataxis)**
  - Reorganized project documentation using the Diataxis framework (Concepts, How-To, Reference, Standards).
  - Replaced monolithic `AGENTS.md` with a structured Site Map to minimize agent token usage.
  - Established `docs/standards/atomic-commit-protocol.md` for consistent git history.
  - Established `docs/standards/coding-conventions.md` for Go and Scripting standards.
  - Created `docs/features/README.md` to define the BDD-driven development process.
  - Implemented `docs/adr/ADR-000-compliance-and-safety.md` for safety and compliance.
  - Added automated verification scripts for all newly established standards.

### Fixed
- Resolved potential infinite loop in `ralph` target by ignoring untracked files in git status check.
- Hardened bash scripts with `set -euo pipefail`.
- Updated `.gitignore` with Go and Python standard patterns.

### Changed
- Refactored `Justfile` to consolidate Quality Review role: Bart now handles both static analysis (code review, ACP verification) and dynamic verification (test execution, BDD validation).
- Updated `PLAN.md` to reflect completion of EPIC-008.

### Removed
- Sanitized repository of legacy orientation files (`INDEX.md`, `MANIFEST.txt`, `00_READ_ME_FIRST.txt`).
- Retired deprecated agent personas (`frink`, `herb`, `troy-mcclure`).

## [0.2.0] - 2026-02-18

### Added
- **EPIC-002: Tmux Agent Orchestration**
  - Implemented `scripts/tmux-orch.sh` for multi-agent session management.
  - Added `just flow`, `just attach`, and `just stop` commands.
  - Established `ADR-002` for Tmux orchestration strategy.
  - Created BDD scenarios for orchestration in `features/tmux_orchestration.feature`.
  - Added support for titled windows and agent suffixes (`ralph-1`).
  - Implemented 80/20 window split for execution and log tailing.

## [0.1.0] - 2026-02-17

### Added
- **EPIC-001: Git Branching Standard**
  - Defined Trunk-Based Coordination model for multi-agent parallel workflows.
  - Implemented BDD scenarios for branch management in `features/git_branching.feature`.
  - Established Architecture Decision Record `ADR-001` for the branching strategy.
  - Defined repository protection rules and merge strategies.
- **Springfield Protocol Foundation**
  - Bootstrapped core 5-agent team: Marge (Product), Lisa (Planning), Ralph (Build), Bart (Quality), Lovejoy (Release).
  - Created `PLAN.md` for value-driven roadmap management.
  - Established the `docs/` structure (features, standards, adr, reference).
  - Created a comprehensive Glossary of terms.
