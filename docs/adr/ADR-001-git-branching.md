# ADR-001: Git Branching Model

**Date:** 2026-02-17
**Status:** Proposed
**Feature:** docs/features/git-branching-model.md
**BDD:** features/git_branching.feature

## Problem
Currently, there is no standardized git branching model defined for the project. This leads to potential conflicts, inconsistent history, and difficulties in collaboration or automation.

## Decision
We will adopt a **Trunk-Based Development** model with short-lived feature branches.

## Rationale
- **Simplicity:** Fewer long-lived branches reduce merge complexity.
- **Speed:** Encourages frequent integration and smaller PRs.
- **Tooling:** Compatible with standard CI/CD practices (GitHub Flow).

## Consequences
- **Requires Discipline:** Developers must merge frequently (daily if possible).
- **Feature Flags:** Incomplete features must be hidden behind flags, not long-lived branches.
- **Release Strategy:** Releases are tagged directly from `main` (no `develop` branch).

## Alternatives Considered
- **GitFlow:** Too complex for our team size; excessive branching overhead.
- **GitHub Flow:** Very similar, but we want explicit release tagging.

## References
- [Trunk Based Development](https://trunkbaseddevelopment.com/)
