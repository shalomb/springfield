# PLAN.md - Product Backlog

> **Marge's Note:** This backlog has been reprioritized using WSJF (Weighted Shortest Job First). Focus is on completing the Springfield Binary migration and retiring shell-based orchestration.
> *Last Updated: 2026-02-20 (Marge approved EPIC-010 for next sprint)*

## 🚀 Active Focus

### EPIC-010: Agent Command Migration (Justfile to Binary)
**WSJF Score: 3.5** (CoD: 21 / Size: 6)
**Value Statement:** For **the Agent Crew**, who **face cognitive load from dual orchestration patterns**, the **Agent Command Migration** is a **port of shell-based agent recipes** to **type-safe Go runners** within the Springfield binary.

**The "Why":** Shell-based agent invocation is fragile, hard to test, and relies on `npm exec` overhead. Moving to Go ensures consistency with EPIC-009 and enables robust unit testing of agent loops.

**Scope:**
- [ ] **Phase 1 (Ralph):** Move Ralph's multi-iteration loop to `internal/agent/ralph_runner.go`.
- [ ] **Phase 2 (Lisa):** Port Lisa's planning and feedback parsing to Go.
- [ ] **Phase 3 (Bart/Lovejoy):** Port Bart's verdict checking and Lovejoy's merge ceremony.
- [ ] **Phase 4 (Cleanup):** Retire all shell-based agent recipes and `npm exec` calls.

**Acceptance Criteria:**
- [ ] All agent prompts migrated to editable markdown in `.github/agents/prompt_{ralph,lisa,bart,lovejoy}.md`.
- [ ] `just ralph`, `just lisa`, etc. are thin wrappers for `springfield agent`.
- [ ] Unit test coverage ≥90% for new runner logic.
- [ ] **Marge Gate:** No regression in agent output or behavior compared to shell versions.

**Attributes:**
- **Status:** 🚧 In Progress
- **Complexity:** Medium
- **Urgency:** High
- **Dependencies:** EPIC-009 (Completed)
- **Depends On:** EPIC-009 is now complete

---

## 📚 Just Completed

### EPIC-009: Springfield Binary Orchestrator & td(1) Integration
**td:** `td-3cc3c3`
**WSJF Score: 4.5** (CoD: 36 / Size: 8)
**Value Statement:** For **the Agent Crew**, who **suffer from fragile orchestration and branch contention**, the **Springfield Binary** is a **type-safe Go orchestrator** that **replaces shell-based Justfile loops** and **uses td(1) for shared planning state**.

**The "Why":** Shell-based orchestration cannot be unit tested and is prone to string-matching errors. Shared state in `td` (SQLite) eliminates planning conflicts across git worktrees.
**Scope:**
- [x] **cmd/springfield:** Implement the orchestration state machine in Go.
- [x] **td(1) Integration:** Use `td` as the source of truth for Epic/Task state.
- [x] **Typed Signals:** Replace keyword grep in `FEEDBACK.md` with `td log --decision`.
- [x] **TODO-{id}.md:** Implement the handoff context deposit protocol.

**Acceptance Criteria:**
- [x] `just do` delegates entirely to `cmd/springfield`.
- [x] State transitions follow the table in ADR-008.
- [x] Multiple Ralph worktrees can run concurrently without `PLAN.md` conflicts.
- [x] **Marge Gate:** Orchestration logic is covered by 90%+ unit test coverage.

**Attributes:**
- **Status:** ✅ Completed (2026-02-20)
- **Complexity:** High
- **Urgency:** Critical
- **Dependencies:** EPIC-007 (Loop logic), ADR-008
- **ADRs:** `docs/adr/ADR-008-planning-state-td-springfield-orchestrator.md`

---

## 📋 Backlog (Prioritized)

### EPIC-005: Agent Governance & Selection
**WSJF Score: 3.25** (CoD: 26 / Size: 8)
**Value Statement:** For **Budget Owners & Developers**, who **need to manage costs and tailor agent behavior**, the **Agent Governance Layer** is a **configuration and control system** that **balances operational flexibility with financial safety**.

**The "Why":** "Infinite loops" in agent logic can bankrupt us. We need a way to say "use this model, for this task, within this budget."
**Scope:**
- [x] **Governance Framework:** ADRs (007, 008) and Quality Indices (Farley, Adzic).
- [ ] **Unified Config (`.springfield.yaml`)** & Global Fallback.
- [ ] **Budget Enforcer:** Per-session and per-day hard limits.
- [ ] **Model Selection Logic:** Swap models based on task complexity.
- [ ] **Tool/Sandbox Mapping:** Define accessible tools.

**Acceptance Criteria:**
- [x] Governance standards (Feedback, Farley, Adzic) are documented and applied.
- [ ] Every LLM call is logged with token count and cost.
- [ ] System rejects requests when budget is exceeded.
- [ ] Agents can be configured via a `.springfield.yaml`.
- [ ] **Marge Gate:** Budget thresholds are agreed upon.

**Attributes:**
- **Status:** 📋 Ready (Phase 1 Done)
- **Complexity:** Medium
- **Urgency:** Medium
- **Dependencies:** EPIC-009 (Orchestrator), EPIC-003 (Logging)

### EPIC-006: Existing Agent Compatibility
**WSJF Score: 2.0** (CoD: 10 / Size: 5)
**Value Statement:** For **Adopters**, who **have existing agent definitions**, the **Compatibility Layer** is a **bridge** that **allows Springfield to run legacy/external agent structures**.

**The "Why":** We shouldn't force a rewrite of all existing `.github/agents` definitions.
**Scope:**
- [ ] Support for `.github/agents`, `.claude/agents`, etc.
- [ ] Precedence logic (Repo > Default)

**Acceptance Criteria:**
- [ ] Springfield agents are primed to load from existing folder structures.
- [ ] **Marge Gate:** Identified legacy agents map successfully.

**Attributes:**
- **Status:** 📋 Ready
- **Complexity:** Medium
- **Urgency:** Low
- **Dependencies:** None

---

## ✅ Completed History

### EPIC-005 (Phase 1): Governance Framework
- **Status:** ✅ Done (2026-02-20)
- **Outcome:** Established ADR-007 (Planning Loop) and ADR-008 (State Boundary). Implemented Farley and Adzic quality indices. Created feedback standard.
- **Retrospective (2026-02-20):**
    - **Learning:** Governance is empirical. ADR-007 Amendment A ensures ADRs are treated as hypotheses to be verified by Ralph's implementation.
    - **Learning:** "Premature refinement is waste" (EPIC-007) is solved by the Fidelity Gradient (Stub -> Options -> Ready).
    - **Signals:** Bart's approval of the framework validates the "Shift-Left" quality model.

### EPIC-004: Agent Sandboxing
- **Status:** ✅ Done (2026-02-19)
- **Outcome:** Agents run inside isolated Axon containers with workspace mounting and resource constraints.
- **Retrospective (2026-02-19):**
    - **Learning:** Security guardrails (`isUnsafeAction`) were too restrictive, blocking standard shell redirection (`>`) which Ralph needs.
    - **Learning:** Simple string matching for `[[FINISH]]` triggers prematurely if not bounded to the end or its own line.

### EPIC-007: Autonomous Development Loop ("just do")
- **Status:** ✅ Done (2026-02-19)
- **Outcome:** Implemented sequential agent chaining (Lisa -> Ralph -> Bart) with `TODO.md` and `FEEDBACK.md` context persistence.
- **Retrospective (2026-02-19):**
    - **Learning:** Simple string matching for `FINISH` and `ACTION:` is too fragile for LLM responses.
    - **Learning:** Ignoring errors in logging/filesystem calls leads to silent failures and QA rejection.

### EPIC-008: Knowledge Architecture (Diataxis)
- **Status:** ✅ Done (2026-02-18)
- **Outcome:** Replaced monolithic `AGENTS.md` with a structured index. Established `docs/standards/` and `docs/adr/`.

---

## 🚩 Technical Debt, Risks & Known Issues

### ⚠️ Known Issues (Minor Feedback)
- **Handoff Filename (EPIC-009):** `TODO-{id}.md` is renamed to `TODO.md` inside the worktree for simplification. This is acceptable but slightly deviates from the ADR.
- **Infinite Loop Protection:** Lisa (Planning Agent) is responsible for handling `StatusBlocked` and ensuring tasks are atomic and non-circular.
- **Justfile Fragility (FIXED IN TODO):** Greedy grep on `FEEDBACK.md` causes false positives.
- **Logger Inefficiency:** Current `pkg/logger` opens and closes two log files for every entry. Needs optimization for high-throughput (e.g., buffered writer).
- **Ghost Feature:** `docs/features/automated_feedback_loop.feature` exists but has no tests.
- **Linting Error:** `internal/sandbox/axon_test.go:88:16` - unchecked `os.Chdir`.
- **Inconsistent safety guardrails:** `isUnsafeAction` blocks `;` and `||` but allows `&&`.
- **Integration Test Debt:** `tests/integration/feedback_loop_test.go:46, 136` - unchecked `os.WriteFile` errors. (Found by Bart).

### ⚡ Risks
- **TR-007:** Springfield binary implementation slip (delays type-safe orchestration).
- **TR-008:** Lisa's ToT/Self-Consistency logic cost/latency (bottleneck).
- **TR-009:** `td(1)` data loss or corruption (single-host SQLite risk).

## 📔 Retrospective (2026-02-20)
- **Signal:** Bart approved EPIC-009 implementation with a PASS verdict.
- **Learning:** The state machine transitions (Planned -> Ready -> InProgress -> Implemented -> Verified) are now correctly handled by the Springfield binary orchestrator.
- **Learning:** Worktree management is robust, preventing stale directory issues during concurrent agent execution.
- **Learning:** Context injection via `TODO.md` (renamed from `TODO-{id}.md` in the worktree) simplifies agent logic.
- **Action:** Move to EPIC-010 to port agent-specific loops (Ralph's multi-iteration, Lisa's planning) from Justfile to Go.
