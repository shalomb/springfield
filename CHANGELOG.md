# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
