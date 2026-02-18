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
- **Concurrency:** **Git Worktrees** prevent file locking and allow one host to support multiple parallel agents (Ralphs) without context contamination.
- **Throughput:** **Merge Queues** resolve the "concurrency collision" problem when multiple agents finish work simultaneously.
- **Flexibility:** The **Coordination Branch/PR** pattern allows the roadmap to be updated dynamically even when code branches (`main`) are protected by rigid human-centric rules.

## Consequences
- **Requires Discipline:** Developers must merge frequently (daily if possible).
- **Tooling Overhead:** The harness must manage worktree lifecycle (cleanup is critical).
- **CI/CD Complexity:** Merge queues require specific provider support (GitHub, GitLab, etc.).
- **Lisa's Access:** In protected repos, Lisa needs a path to update `PLAN.md` (either a separate branch or an bypass-permission).

## Alternatives Considered
- **GitFlow:** Too complex for our team size; excessive branching overhead.
- **GitHub Flow:** Very similar, but we want explicit release tagging.

## References
- [Trunk Based Development](https://trunkbaseddevelopment.com/)
