# TODO: EPIC-005 Phase 3 - Agent Cost Controls & Model Optimization

**Epic ID:** td-e1fd16
**Status:** ACTIVE
**Primary Task:** Model Selection per Agent (td-19b078)

## Intent
We need to move away from a "one-model-fits-all" approach. High-reasoning tasks like planning (Lisa) and code review (Bart) benefit from more capable models (Opus), while implementation (Ralph) is often better suited to faster, balanced models (Sonnet). Development remains on Haiku for cost efficiency.

## Approach
1.  **Config Update:** Enhance `internal/config` to support per-agent model overrides.
2.  **Agent Parameterization:** Update `internal/agent` to accept a `Model` parameter in its configuration.
3.  **Orchestrator Integration:** Ensure the `Orchestrator` passes the correct role-specific model when initializing runners.
4.  **Verification:** Add unit tests to verify the model flag is correctly passed to the `pi` CLI runner.

## Tasks (from td)
- [ ] **td-19b078**: Model Selection per Agent (tuning roles)
- [ ] **td-7d836a**: Per-Session Budget Enforcement (Blocked on JSON)
- [ ] **td-752039**: Per-Day Budget Tracking (Blocked on JSON)

## Constraints
- **ADR-005**: Follow coding standards.
- **Atomic Commit Protocol**: Every commit must be a coherent, testable unit.
- **Safety**: Maintain backward compatibility with the existing global model flag.

## Context
- `internal/config/config.go`: Current configuration logic.
- `internal/agent/agent.go`: Agent runner implementation.
- `internal/orchestrator/orchestrator.go`: Where agents are instantiated.

---
*Created by Lisa Simpson (Planning Agent)*
