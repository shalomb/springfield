# TODO.md - Sprint for EPIC-002: Tmux Agent Orchestration

## Context
Running 5 agents (Marge, Lisa, Ralph, Bart, Lovejoy) in separate terminals is unmanageable. 
We need a "command center" view using Tmux to allow concurrent execution and easy session management.
See [PLAN.md#EPIC-002] for full value statement and scope.

## Strategy
- Use a single Tmux session named `springfield`.
- Each agent gets its own named window/pane.
- Use `just` as the entry point for orchestration.
- Support "detached" mode so agents can work in the background.

## Tasks
- [x] Task 0: Define ADR for Tmux Orchestration ✅ @Lisa 2026-02-17
- [x] Task 1: Create BDD Scenarios for Tmux Orch ✅ @Lisa 2026-02-17
- [x] Task 2: Create `scripts/tmux-orch.sh` prototype ✅ @Ralph 2026-02-17
- [x] Task 3: Implement `just flow` command ✅ @Ralph 2026-02-17

## BDD Scenarios
See `features/tmux_orchestration.feature` (to be created).
