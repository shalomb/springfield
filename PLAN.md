# PLAN.md - Product Backlog

> **Marge's Note:** This backlog has been reprioritized using WSJF (Weighted Shortest Job First). Focus is on finishing the Autonomous Loop (EPIC-007) and then moving to Governance (EPIC-005) to control costs before we scale.
> *Last Updated: 2026-02-19*

## 🚀 Active Focus

### EPIC-007: Autonomous Development Loop ("just do")
**WSJF:** ∞ (In Flight / Critical Path)
**Value Statement:** For **Developers**, who **want to delegate end-to-end feature implementation**, the **Autonomous Development Loop** is a **workflow orchestrator** that **automates the cycle of planning, coding, reviewing, and refining**.

**The "Why":** Manual handoffs between agents (Planning -> Coding -> Review) are inefficient. We need a closed-loop system where agents collaborate iteratively to complete complex tasks without constant human interruption.

**Scope:**
- [ ] `just do` command as the entry point.
- [ ] Sequential agent chaining: Lisa -> Ralph -> Herb -> Bart.
- [ ] Context persistence: `TODO.md` (Plan) and `FEEDBACK.md` (Review).
- [ ] Dynamic branching: Lisa manages feature branches based on specs.
- [ ] Iteration logic: Loop repeats based on feedback severity.
- [ ] Exit criteria: Handover to Lovejoy for merging when "Done".

**Acceptance Criteria:**
- [ ] `just do` initiates the loop in the current context.
- [ ] **Lisa (Planner):**
    - Parses `PLAN.md` and `FEEDBACK.md`.
    - Generates BDD scenarios in `docs/features/`.
    - Creates/Updates `TODO.md` with prioritized tasks (TDD first, Refactor last).
    - Manages git branches (creates `feat/xxx` if on `main`).
- [ ] **Ralph (Builder):**
    - Executes tasks from `TODO.md`.
    - Continues working as long as `TODO.md` exists OR uncommitted changes remain.
    - Finalizes work by committing remaining changes and removing `TODO.md`.
    - Updates task status in real-time.
- [ ] **Herb & Bart (Reviewers):**
    - Herb reviews code changes (Static Analysis/Style).
    - Bart verifies functionality against BDD scenarios.
    - Both populate `FEEDBACK.md` with findings.
- [ ] **Orchestrator:**
    - Detects loop continuation (Is `TODO.md` empty? Is `FEEDBACK.md` critical?).
    - Invokes `just lovejoy` for merge when cycle is complete.

**Attributes:**
- **Status:** 🏗️ In Progress (Recovering)
- **Complexity:** High
- **Urgency:** High
- **Dependencies:** EPIC-002 (Tmux), EPIC-003 (Logging)

---

## 📋 Backlog (Prioritized)

### EPIC-005: Agent Governance & Selection
**WSJF Score: 3.25** (CoD: 26 / Size: 8)
**Value Statement:** For **Budget Owners & Developers**, who **need to manage costs and tailor agent behavior**, the **Agent Governance Layer** is a **configuration and control system** that **balances operational flexibility with financial safety**.

**The "Why":** "Infinite loops" in agent logic can bankrupt us. We need a way to say "use this model, for this task, within this budget."
**Scope:**
- [ ] **Unified Config (`.springfield.yaml`)** & Global Fallback.
- [ ] **Budget Enforcer:** Per-session and per-day hard limits.
- [ ] **Model Selection Logic:** Swap models based on task complexity.
- [ ] **Tool/Sandbox Mapping:** Define accessible tools.

**Acceptance Criteria:**
- [ ] Every LLM call is logged with token count and cost.
- [ ] System rejects requests when budget is exceeded.
- [ ] Agents can be configured via a `.springfield.yaml`.
- [ ] **Marge Gate:** Budget thresholds are agreed upon.

**Attributes:**
- **Status:** 📋 Ready
- **Complexity:** Medium
- **Urgency:** Medium
- **Dependencies:** EPIC-003 (Logging)

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

### EPIC-004: Agent Sandboxing
**WSJF Score: 1.65** (CoD: 33 / Size: 20)
**Value Statement:** For **System Administrators**, who **fear agents destroying the host system**, the **Sandboxing Environment** is a **security boundary** that **ensures safe execution of arbitrary code**.

**The "Why":** Agents like Ralph execute code. Running this as root/user on the host is dangerous.
**Scope:**
- [ ] Docker/Container-based execution context.
- [ ] Workspace mounting strategy.
- [ ] Resource constraints (CPU/Memory).

**Acceptance Criteria:**
- [ ] Agents run inside a defined container image.
- [ ] Agents cannot access host files outside the workspace.
- [ ] **Marge Gate:** Security model is validated against "jailbreak" patterns.

**Attributes:**
- **Status:** 🔍 Discovery
- **Complexity:** High
- **Urgency:** High (Security)
- **Dependencies:** EPIC-003 (Logging)

---

## ✅ Completed History

### EPIC-008: Knowledge Architecture (Diataxis)
- **Status:** ✅ Done
- **Outcome:** Replaced monolithic `AGENTS.md` with a structured index. Established `docs/standards/` and `docs/adr/`.

### EPIC-001: Git Branching Standard
- **Status:** ✅ Done
- **Outcome:** Defined `feat/` and `fix/` conventions. Ratified `docs/standards/git-branching.md`.

### EPIC-002: Tmux Agent Orchestration
- **Status:** ✅ Done
- **Outcome:** `just flow` launches agent mesh. Named windows and detached logging implemented.

### EPIC-003: Logging & Observability
- **Status:** ✅ Done
- **Outcome:** JSON structured logging with `agent_id` and `task_id`. `just logs` created.

---

## 🚩 Technical Debt & Risks

### TR-001: PLAN.md Merge Contention
- **WSJF:** High (Quick Win)
- **Risk:** High-concurrency merges will cause conflicts in the single `PLAN.md` file.
- **Mitigation:** Future epic to split status into individual files (e.g., `docs/plans/EPIC-XXX.status`).

### TR-002: Coordination Branch Race Conditions
- **Risk:** Lisa's planning commits may conflict with automated downstream syncs from `main`.
- **Mitigation:** Future investigation into "Planning Locks".

### TR-003: Worktree Lifecycle Management
- **Risk:** Crashed agents leave "ghost" worktrees.
- **Mitigation:** Future task for `just gc-worktrees`.

### TR-004: Roadmap/Code Decoupling
- **Risk:** PR gates on `main` prevent timely roadmap updates.
- **Mitigation:** Future ADR to decide if `PLAN.md` should move to a separate coordination repo.

## EPIC-XXX: Continuous Improvement
Recurring epic for retrospective and process iteration.
