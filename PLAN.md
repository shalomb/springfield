# PLAN.md - Epic Backlog

## EPIC-008: Knowledge Architecture (Diataxis)
**Value Statement:** For **Agents and Contributors**, who **need precise context without token overload**, the **Modular Documentation Structure** is a **knowledge retrieval system** that **minimizes hallucination and maximizes relevance**.

### 📝 Retrospective
- **Context Management:** Separating documentation into Diataxis categories significantly reduces noise for agents. Using a site-map (`AGENTS.md`) allows agents to selectively load only the relevant standards, saving tokens and improving focus.
- **Standards Enforcement:** Automated verification scripts (`scripts/test_*.sh`) are essential for maintaining the integrity of the new architecture as the project scales.
- **Legacy Cleanup:** Removing retired agents and orientation files reduces mental overhead for both human and AI contributors.
- **Improvements for Next Cycle:** Consider further atomization of `PLAN.md` to avoid merge contention (as noted in TR-001). Integrate documentation linting into the graduated test ladder (`just test`).

**The "Why":** A monolithic `AGENTS.md` is hard to maintain and expensive to process. By adopting Diataxis, we separate "How-to" (Guides) from "Explanation" (Context) and "Reference" (Standards), allowing agents to load only what they need.

**Scope:**
- ✅ Replace monolithic `AGENTS.md` with a structured index (Site Map).
- ✅ Ensure all core standards exist as separate files in `docs/standards/`.
- ✅ Ensure all architectural decisions exist in `docs/adr/`.
- ✅ Verify agents can traverse links or are configured to read specific paths.
- ✅ Sanitize repository of legacy orientation files and retired personas.

**Acceptance Criteria:**
- [x] `AGENTS.md` acts as a "Site Map" linking to `docs/`.
- [x] `docs/standards/coding-conventions.md` is the source of truth for code style.
- [x] `docs/standards/git-branching.md` is the source of truth for workflow.
- [x] `docs/reference/agents.md` defines the agent team roles.
- [x] `docs/features/README.md` explains the BDD process.

**Tasks (from TODO.md):**
- [x] **Task 1: Create `AGENTS.md` Site Map**
- [x] **Task 2: Define Atomic Commit Protocol (ACP)**
- [x] **Task 3: Establish Coding Conventions**
- [x] **Task 4: Create Features README**
- [x] **Task 5: Compliance and Safety (ADR-000)**
- [x] **Task 6: Cross-Reference and Finalize Index**

**Attributes:**
- **Status:** ✅ Done
- **Complexity:** Low
- **Urgency:** Medium
- **Dependencies:** None

---

## EPIC-007: Autonomous Development Loop ("just do")
**Value Statement:** For **Developers**, who **want to delegate end-to-end feature implementation**, the **Autonomous Development Loop** is a **workflow orchestrator** that **automates the cycle of planning, coding, reviewing, and refining**.

**The "Why":** Manual handoffs between agents (Planning -> Coding -> Review) are inefficient. We need a closed-loop system where agents collaborate iteratively to complete complex tasks without constant human interruption.

**Scope:**
- ✅ `just do` command as the entry point.
- ✅ Sequential agent chaining: Lisa -> Ralph -> Herb -> Bart.
- ✅ Context persistence: `TODO.md` (Plan) and `FEEDBACK.md` (Review).
- ✅ Dynamic branching: Lisa manages feature branches based on specs.
- ✅ Iteration logic: Loop repeats based on feedback severity.
- ✅ Exit criteria: Handover to Lovejoy for merging when "Done".

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
- **Status:** 📋 Planned
- **Complexity:** High
- **Urgency:** High
- **Dependencies:** EPIC-002 (Tmux), EPIC-003 (Logging)

---

## EPIC-001: Git Branching Standard
**Value Statement:** For **Developers**, who **struggle with inconsistent history and merge conflicts**, the **Git Branching Standard** is a **protocol** that **ensures clean collaboration and predictable releases**.

**The "Why":** Without a standardized model, we risk "merge hell", lost code, and unclear release points. We need a shared mental model for how code moves from laptop to production.
**Scope:**
- ✅ Trunk-based development definition
- ✅ Branch naming conventions (feat/, fix/)
- ✅ Merge strategy (Squash vs Merge Commit)
- ❌ Automated CI/CD pipeline implementation (future epic)

**Acceptance Criteria:**
- [x] `docs/standards/git-branching.md` exists and is ratified. ✅
- [ ] Team members can explain the lifecycle of a feature branch.
- [ ] Repository settings enforce the strategy (if applicable).
- [ ] **BDD Scenarios:** `features/git_branching.feature`

**Attributes:**
- **Status:** ✅ Done
- **Complexity:** Low
- **Urgency:** High (Foundational)
- **Dependencies:** None
- **ADRs:** `docs/adr/ADR-001-git-branching.md`

**Tasks:**
- [x] Task 1: Create Git Branching Strategy Document ✅ @Ralph 2026-02-17 [Verified @Herb]
- [x] Task 2: Define ADR for Branching Strategy ✅ @Lisa 2026-02-17
- [x] Task 3: Configure Repository Protection Rules (Simulated) ✅ @Ralph 2026-02-17

---

## EPIC-002: Tmux Agent Orchestration
**Value Statement:** For **Developers/Operators**, who **need to run multiple agents simultaneously**, the **Tmux Orchestration Layer** is a **tooling set** that **allows concurrent execution without window clutter**.

**The "Why":** Running 5 agents (Lisa, Ralph, etc.) in separate terminals is unmanageable. We need a "command center" view.
**Scope:**
- ✅ Script to launch/attach named tmux sessions
- ✅ `just` command integration
- ✅ Detached mode support
- ✅ Smart session reuse (detect existing `$TMUX`)
- ❌ Web-based management UI

**Acceptance Criteria:**
- [x] `just flow` launches the full agent mesh in a tmux session. ✅
- [x] Users can toggle between agent views easily. ✅
- [x] Logs are preserved in detached panes. ✅
- [x] Windows are titled with agent names (e.g. `ralph-1`). ✅

**Attributes:**
- **Status:** ✅ Done
- **Complexity:** Medium
- **Urgency:** Medium
- **Dependencies:** None
- **ADRs:** `docs/adr/ADR-002-tmux-orchestration.md`

---

## EPIC-003: Logging & Observability
**Value Statement:** For **Operators**, who **cannot debug failed agent actions**, the **Structured Logging System** is a **framework** that **provides traceability and context for every action**.

**The "Why":** Debugging "why did Ralph do that?" is currently impossible with standard stdout. We need structured, grep-able logs.
**Scope:**
- ✅ JSON structured logging format
- ✅ Standardized log levels (INFO, DEBUG, TRACE)
- ✅ Agent Identity in log context
- ❌ ELK/Splunk integration

**Acceptance Criteria:**
- [x] All agents emit JSON logs to a central file/stream. ✅ [Verified @Herb]
- [x] Logs contain `agent_id`, `task_id`, and `timestamp`. ✅ [Verified @Herb]
- [x] CLI tool exists to tail/filter these logs. ✅ (`just logs`) [Verified @Herb]
- [x] **BDD Scenarios:** `features/logging.feature` ✅ [Verified @Herb]
- [x] **ADR:** `docs/adr/ADR-003-logging-standard.md` ✅ [Verified @Herb]

**Attributes:**
- **Status:** ✅ Done
- **Complexity:** Medium
- **Urgency:** High (Debugging)
- **Dependencies:** None

---

## EPIC-004: Agent Sandboxing
**Value Statement:** For **System Administrators**, who **fear agents destroying the host system**, the **Sandboxing Environment** is a **security boundary** that **ensures safe execution of arbitrary code**.

**The "Why":** Agents like Ralph execute code. Running this as root/user on the host is dangerous. We need containment.
**Scope:**
- [ ] Docker/Container-based execution context
- [ ] Workspace mounting strategy
- [ ] Resource constraints (CPU/Memory)
- ❌ Network restriction policies (Deferred)
- ❌ Full VM virtualization (Out of Scope)

**Acceptance Criteria:**
- [ ] Agents run inside a defined container image (or similar isolation).
- [ ] Agents cannot access host files outside the mounted workspace.
- [ ] Workspace state is preserved between runs.
- [ ] **BDD Scenarios:** `features/sandboxing.feature`
- [ ] **Marge Gate:** Performance impact is measured and accepted by stakeholders.
- [ ] **Marge Gate:** Security model is validated against common "jailbreak" patterns.

**Attributes:**
- **Status:** 🔍 Discovery
- **Complexity:** High
- **Urgency:** High (Security)
- **Dependencies:** EPIC-003 (Logging)
- **ADRs:** `docs/adr/ADR-004-agent-sandboxing.md` (Planned)

**Risks:**
- **TR-005:** `pi` environment constraints may prevent Docker-in-Docker or nested virtualization.
- **TR-006:** Filesystem mounting latency could impact Ralph's performance.

**Tasks:**
- [ ] Task 1: Research `pi` environment capabilities for isolation (Docker, podman, nsenter)
- [ ] Task 2: Draft ADR-004 with proposed isolation strategy
- [ ] Task 3: Create `features/sandboxing.feature`
- [ ] Task 4: Prototype isolation script

---

## EPIC-005: Agent Governance & Selection
**Value Statement:** For **Budget Owners & Developers**, who **need to manage costs and tailor agent behavior**, the **Agent Governance Layer** is a **configuration and control system** that **balances operational flexibility with financial safety**.

**The "Why":** "Infinite loops" in agent logic can bankrupt us, and generic agent prompts don't always fit specific project needs. We need a way to say "use this model, for this task, within this budget."

**Scope:**
- [ ] **Unified Config (`.springfield.yaml`):** Repo-level overrides for agent behavior and selection.
- [ ] **Global Configuration Fallback:** Support for `~/.config/springfield/config.yaml` as a base layer.
- [ ] **Token counting middleware:** Track usage across different LLM providers (Pi, Claude, Copilot, Gemini).
- [ ] **Budget Enforcer:** Per-session and per-day hard limits to prevent runaway costs.
- [ ] **Model Selection Logic:** Ability to swap models (e.g., GPT-4 for planning, GPT-3.5 for simple tasks) based on task complexity.
- [ ] **Prompt Engineering Injection:** Support for project-specific system prompts and identity definitions.
- [ ] **Tool/Sandbox Mapping:** Define which directories and tools are accessible to specific agents.
- [ ] **Output Stream Handling:** Intercepting agent output for logging and cost analysis.
- ❌ Real-time billing API integration (Out of Scope)

**Technical Requirements & Discovery:**
- **Provider Surface Areas:** Examine CLI interfaces for `pi`, `claude` (CLI), `copilot` (prerelease), and `gemini-cli`.
- **Controllable Aspects:**
    - Model selection (e.g., GPT-4 vs GPT-3.5).
    - Prompt identity (e.g., "You are Ralph, an expert in debugging... use TDD").
    - Formatting context (e.g., "Always include a JSON blob...").
    - Identity invocation (e.g., @ralph mentions vs configuration loading).
    - Resource isolation (mounting safe directories into sandboxes).

**Acceptance Criteria:**
- [ ] Every LLM call is logged with token count and estimated cost.
- [ ] System rejects requests when budget is exceeded.
- [ ] Reporting command (`just budget`) shows daily spend.
- [ ] Agents can be configured via a `.springfield.yaml` in the repo root.
- [ ] **Marge Gate:** The `.springfield.yaml` schema is validated as "human-friendly" (easy for a dev to write without a manual).
- [ ] **Marge Gate:** Budget thresholds are agreed upon by stakeholders; "Fail-safe" mode is implemented (agents stop gracefully when budget is hit).
- [ ] **Marge Gate:** Privacy check: ensure sensitive project prompts aren't leaked in global logs.

**Attributes:**
- **Status:** 📋 Ready
- **Complexity:** Medium
- **Urgency:** Medium
- **Dependencies:** EPIC-003 (Logging), EPIC-004 (Sandboxing)
- **ADRs:** `docs/adr/ADR-005-agent-governance.md` (Planned)

---

## EPIC-006: Existing Agent Compatibility
**Value Statement:** For **Adopters**, who **have existing agent definitions**, the **Compatibility Layer** is a **bridge** that **allows Springfield to run legacy/external agent structures**.

**The "Why":** We shouldn't force a rewrite of all existing `.github/agents` definitions. We should embrace them.
**Scope:**
- [ ] Support for `.github/agents`, `.claude/agents`, etc.
- [ ] Precedence logic (Repo > Default)
- ❌ Conversion/Migration tools (Out of Scope)

**Acceptance Criteria:**
- [ ] Springfield agents are primed to load from existing folder structures.
- [ ] Repo-defined agents override defaults.
- [ ] **Marge Gate:** Identified legacy agents (e.g., from `pi` defaults) map successfully to Springfield roles.

**Attributes:**
- **Status:** 📋 Ready
- **Complexity:** Medium
- **Urgency:** Low
- **Dependencies:** None

---

## Technical Debt & Risks (Backlog)

### 🚩 TR-001: PLAN.md Merge Contention
- **Risk:** High-concurrency merges will cause conflicts in the single `PLAN.md` file.
- **Mitigation:** Future epic to split status into individual files (e.g., `docs/plans/EPIC-XXX.status`).

### 🚩 TR-002: Coordination Branch Race Conditions
- **Risk:** Lisa's planning commits may conflict with automated downstream syncs from `main`.
- **Mitigation:** Future investigation into "Planning Locks" or atomic reconciliation logic.

### 🚩 TR-003: Worktree Lifecycle Management
- **Risk:** Crashed agents leave "ghost" worktrees and fill up disk space.
- **Mitigation:** Future task for `just gc-worktrees` cleanup routine.

### 🚩 TR-004: Roadmap/Code Decoupling
- **Risk:** PR gates on `main` prevent timely roadmap updates.
- **Mitigation:** Future ADR to decide if `PLAN.md` should move to a separate coordination repo.

---

## EPIC-XXX: Continuous Improvement

Allow for the system to do a retrospective after each major release to identify both technical and process improvements. This will be a recurring epic that ensures we are always iterating on our practices and tooling based on data from the agents.
