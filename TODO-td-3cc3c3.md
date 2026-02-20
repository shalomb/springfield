# TODO-td-3cc3c3.md — Springfield Binary Orchestrator & td(1) Integration

## Intent (Immutable — from Marge's Feature Brief)
**User Need:** As an agent crew, we need a reliable, testable orchestration layer that doesn't suffer from branch contention or fragile string-matching on markdown files.
**Acceptance Criteria:**
- Springfield binary (`cmd/springfield`) owns the state machine.
- `td(1)` (SQLite) is the shared source of truth for all agents/worktrees.
- `just do` is a thin wrapper around the binary.
- Multiple Ralph worktrees can run simultaneously without `PLAN.md` conflicts.
- Orchestration logic is unit-testable in Go.

## Approach (Decided by Lisa at LRM — fixed for this iteration)
- **State Machine:** Implement the `EpicStatus` enum and transition table from ADR-008 §4.
- **td(1) Integration:** Shell out to `td` as a subprocess. Use `td query` to find the next ready Epic.
- **Signal Handling:** Bart will now use `td log --decision {type}` to signal progress. The orchestrator must parse this signal.
- **Handoff Deposit:** Implement the `git worktree add` and `cp TODO-{id}.md` protocol from ADR-008 §5.
- **Justfile:** Refactor `do` recipe to call `./bin/springfield orchestrate`.

## Constraints (Inherited — not negotiable)
- **ADR-008:** Follow the state machine and architecture defined in the ADR.
- **Atomic Commit Protocol:** Every task completion must be an atomic commit.
- **Farley Index:** Orchestrator code must be Fast, Maintainable, and Repeatable.
- **Adzic Index:** Success criteria must be Intention-Revealing and Business-Readable.
- **Tech Debt:** Fix `tests/integration/feedback_loop_test.go:46, 136` (unchecked errors) as part of the test cleanup in this epic.
