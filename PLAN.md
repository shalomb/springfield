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

### EPIC-011: Pipeline Redesign Bootstrapping (v0.8.0)
**Status:** Ready for Lisa
**Priority:** P0 (Highest)

**Objective:** Use the current Springfield agents (v0.7.0) to implement the deterministic pipeline redesign (v0.8.0), transitioning to strict Inner/Outer loops and LRM (Last Responsible Moment) planning. This epic is sequenced so that early tasks (like Bart's output) unblock the agents' ability to execute later tasks (like Lisa's ToT).

**Phase 1: The Feedback Loop (Bart → Lisa)**
- [x] Update `bart.prompt.md.tpl` to output structured `## Retrospective Signal`.
- [x] Wire Bart's retrospective output to `td log --decision <epic-id>`.
- [x] Author `docs/standards/epic-decomposition-protocol.md`.
- *Bootstrapping impact:* Lisa can now read cross-iteration learnings via `td` before planning her next task.

**Phase 2: Lisa's Last Responsible Moment (LRM)**
- [x] Update `lisa.prompt.md.tpl` to run Tree of Thoughts (ToT) and Self-Consistency loops before choosing an approach.
- [x] Implement the `TODO-{id}.md` Context/Constraint Layer handoff.
- [ ] Update orchestrator to natively enforce structural checks on `springfield signal`.
- *Bootstrapping impact:* Ralph receives a strict, context-driven boundary and tests, drastically reducing his context window while empowering architectural autonomy.

**Phase 3: Trigger & Escalation Tuning**
- [x] Update orchestrator to trigger Lovejoy *only* when `StatusAllEpicsDone` is true for the release.
- [x] Add OHECI escalation path to `ralph.prompt.md.tpl` for assumption breaks.
- [x] Add GECR loop to `marge.prompt.md.tpl` for Feature Brief generation.
- *Bootstrapping impact:* Eliminates wasted LLM tokens on per-epic changelogs and context-spiral debugging.

---

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
