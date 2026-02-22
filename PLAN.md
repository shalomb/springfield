# PLAN.md - Springfield Product Backlog

**Last Updated:** 2026-02-22 13:00 GMT+1  
**Status:** EPIC-010 (Autonomous Control Plane) Active  

---

## 🚀 Current Release: v0.7.0 (Development Cycle)

### Autonomous Control Plane ⭐ ACTIVE
**td:** `td-db921a`
**Status:** In Progress
**Priority:** P1

**Objective:** Transform Springfield from a "Task Runner" into a **Stateful Control Plane** that manages the lifecycle, environment, and state transitions of all agents via a strict signaling protocol.

**Scope:**
- [ ] **td-b6dae6**: Implement `springfield signal` command and validation logic.
- [ ] **td-d59930**: Implement Sentinel Token generation and tracking in Agent Runner.
- [ ] **td-278c10**: Implement Daemon Orchestrator loop (poll `td`, spawn/reap workers).
- [ ] **td-e53407**: Migrate Prompt Files to Templates (`.tpl`) and inject Sentinel.

**Acceptance Criteria:**
- [ ] Agents terminate sessions *only* via `springfield signal`.
- [ ] Orchestrator rejects signals with invalid Sentinels.
- [ ] Ralph is automatically re-spawned in the same worktree if Bart signals failure.
- [ ] Lisa can signal "Session Complete" without binding to a specific Epic ID.
- [ ] All environment management (git worktree add/rm) is handled by the Daemon, not the Agent.

**Feature Briefs:**
- [Features & Architecture](docs/features/autonomous-control-plane.md)
- [Signaling Protocol](docs/standards/signaling-protocol.md)
- [BDD Specs](tests/integration/features/control_plane.feature)

---

## 📋 Backlog (Prioritized)

### EPIC-005 Phase 3: Agent Cost Controls & Model Optimization
**td:** `td-e1fd16`
**Status:** Paused (Pre-requisite for robust scaling, but EPIC-010 architecture comes first)
- [ ] Per-Session Budget Enforcement
- [ ] Per-Day Budget Tracking

---

## ✅ Completed History

### EPIC-005 Phase 2: Robust Structured Output Parsing
**td:** `td-7b0cb8`
- **Status:** ✅ Done (2026-02-21)
- **Outcome:** Lexical Sanitizer (MarkdownSanitizer) integrated. Promise Semantic Contract enforced.

---

*Maintained by Lisa Simpson (Planning Agent) with input from the team.*
