# PLAN.md - Product Backlog

> **Marge's Note:** This backlog has been reprioritized using WSJF (Weighted Shortest Job First). Focus is on finishing the Autonomous Loop (EPIC-007) and then moving to Governance (EPIC-005) to control costs before we scale.
> *Last Updated: 2026-02-19*

## 🚀 Active Focus

### EPIC-004: Agent Sandboxing
**WSJF Score: 1.65** (CoD: 33 / Size: 20)
**Value Statement:** For **System Administrators**, who **fear agents destroying the host system**, the **Sandboxing Environment** is a **security boundary** that **ensures safe execution of arbitrary code**.

**The "Why":** Agents like Ralph execute code. Running this as root/user on the host is dangerous. We need containment.
**Scope:**
- [x] Axon library-based execution context (Migrated from CLI)
- [x] Workspace mounting strategy (Workspace isolation)
- [x] Resource constraints (CPU/Memory)
- ❌ Network restriction policies (Deferred)
- ❌ Full VM virtualization (Out of Scope)

**Acceptance Criteria:**
- [x] Agents run inside an isolated Axon container (via `pkg/executor`).
- [x] Agents cannot access host files outside the mounted workspace.
- [x] Workspace state is preserved between runs.
- [x] **BDD Scenarios:** `features/sandboxing.feature`
- [x] **Marge Gate:** Performance impact is measured and accepted by stakeholders.
- [x] **Marge Gate:** Security model is validated against common "jailbreak" patterns.

**Attributes:**
- **Status:** ✅ Done
- **Complexity:** High
- **Urgency:** High (Security)
- **Dependencies:** EPIC-003 (Logging)
- **ADRs:** `docs/adr/ADR-004-agent-sandboxing.md`, `docs/adr/ADR-005-axon-library-migration.md`

**Retrospective (2026-02-19):**
- **Learning:** Security guardrails (`isUnsafeAction`) were too restrictive, blocking standard shell redirection (`>`) which Ralph needs.
- **Learning:** Simple string matching for `[[FINISH]]` triggers prematurely if not bounded to the end or its own line.
- **Action:** Move safety logic refinement to Known Issues for further optimization, but immediately fix blockers.

**Risks:**
- **TR-005:** `pi` environment constraints may prevent Docker-in-Docker or nested virtualization.
- **TR-006:** Filesystem mounting latency could impact Ralph's performance.

**Corrective Actions (Priority):**
- [x] **CA-1: Robust `FINISH` Detection.** Use `[[FINISH]]` marker or similar to avoid false positives.
- [x] **CA-2: Explicit Error Handling.** Ensure `logger.Log` and `os.Chdir` errors are not ignored.
- [x] **CA-3: Regex Action Extraction.** Use `(?m)^ACTION:\s*(.+)$` for reliable extraction.
- [x] **CA-4: Strengthen Safety Guardrails.** Refine `isUnsafeAction` or migrate to Axon-native allowlist.
- [x] **CA-5: Repair Test Infrastructure.** Fix `tests/unit/test_logger_concurrency.py` and missing scripts.

**Tasks:**
- [x] Task 1: Research `pi` environment capabilities for isolation (Docker, podman, nsenter)
- [x] Task 2: Draft ADR-004 with proposed isolation strategy
- [x] Task 3: Create `features/sandboxing.feature`
- [x] Task 4: CLI Prototype (Superseded by library integration)
- [x] Task 10: Integrate Axon Library (`pkg/executor`)
- [x] Task 11: Implement Workspace Isolation via Axon Volume Mounting
- [x] Task 12: Implement Resource Constraints (CPU/Memory) in `internal/sandbox/axon.go`

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

---

## ✅ Completed History

### EPIC-007: Autonomous Development Loop ("just do")
- **Status:** ✅ Done
- **Outcome:** Implemented sequential agent chaining (Lisa -> Ralph -> Bart) with `TODO.md` and `FEEDBACK.md` context persistence. Consolidated Quality Review role into Bart (static + dynamic verification). `just do` entrypoint stabilized.
- **Retrospective (2026-02-19):** 
    - **Learning:** Simple string matching for `FINISH` and `ACTION:` is too fragile for LLM responses.
    - **Learning:** Ignoring errors in logging/filesystem calls leads to silent failures and QA rejection.
    - **Signals:** Bart's pessimism is a necessary filter for "happy path" implementation.

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

## 🚩 Technical Debt, Risks & Known Issues

### ⚠️ Known Issues (Minor Feedback)
- **Logger Inefficiency:** Current `pkg/logger` opens and closes two log files for every entry. Needs optimization for high-throughput (e.g., buffered writer).
- **Ghost Feature:** `docs/features/automated_feedback_loop.feature` exists but has no tests. (Moved to TODO for implementation).
- **Linting Error:** `internal/sandbox/axon_test.go:88:16` - unchecked `os.Chdir`.
- **Inconsistent safety guardrails:** `isUnsafeAction` blocks `;` and `||` but allows `&&`. Both `&&` and `;` can be used to chain malicious commands. Blocking `||` (logical OR) prevents common fallback patterns in shell scripts.
- **Multi-action inefficiency:** `extractAction` only extracts the first `ACTION:` from an LLM response. If the LLM provides multiple actions, they must be processed one-by-one in subsequent loop iterations.

---

## EPIC-XXX: Continuous Improvement
Recurring epic for retrospective and process iteration.
