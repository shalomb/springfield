# TODO.md - Sprint for EPIC-003: Logging & Observability

## Context
Tracing agent actions is difficult without structured logs. 
We need a unified JSON logging format that includes agent context and timestamps.
See [PLAN.md#EPIC-003] for full value statement and scope.

## Strategy
- Implement a Python-based logger or a shell utility.
- Each agent write to `logs/[agent_name].log`.
- Log format: `{"timestamp": "...", "agent": "...", "epic": "...", "level": "...", "message": "..."}`.
- Update `scripts/tmux-orch.sh` to ensure log directory existence.

## Tasks
- [x] Task 1: Define Logging ADR (ADR-003) ✅ @Lisa 2026-02-18
- [x] Create BDD Scenarios for Logging ✅ @Lisa 2026-02-18
- [x] Task 2: Create structured logging utility ✅ @Ralph 2026-02-18
  - Assigned to: Ralph
  - Subtasks:
    - [x] Create `scripts/logger.py` ✅
    - [x] Support JSON output ✅
    - [x] Add basic unit tests. ✅ @Herb 2026-02-18
- [x] Task 3: Integrate logger into core scripts ✅ @Ralph 2026-02-18
  - Assigned to: Ralph
  - Subtasks:
    - [x] Update `scripts/tmux-orch.sh` to use the new logger for orchestration events. ✅
    - [x] Add `just logs` command to tail all logs. ✅

## BDD Scenarios
See `features/logging.feature` (to be created).
