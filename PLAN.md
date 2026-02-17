# PLAN.md - Epic Backlog

## EPIC-001: Git Branching Standard
**Value Statement:** For **Developers**, who **struggle with inconsistent history and merge conflicts**, the **Git Branching Standard** is a **protocol** that **ensures clean collaboration and predictable releases**.

**The "Why":** Without a standardized model, we risk "merge hell", lost code, and unclear release points. We need a shared mental model for how code moves from laptop to production.
**Scope:**
- ✅ Trunk-based development definition
- ✅ Branch naming conventions (feat/, fix/)
- ✅ Merge strategy (Squash vs Merge Commit)
- ❌ Automated CI/CD pipeline implementation (future epic)

**Acceptance Criteria:**
- [ ] `docs/standards/git-branching.md` exists and is ratified.
- [ ] Team members can explain the lifecycle of a feature branch.
- [ ] Repository settings enforce the strategy (if applicable).
- [ ] **BDD Scenarios:** `features/git_branching.feature`

**Attributes:**
- **Status:** 📋 Ready
- **Complexity:** Low
- **Urgency:** High (Foundational)
- **Dependencies:** None
- **ADRs:** `docs/adr/ADR-001-git-branching.md`

---

## EPIC-002: Tmux Agent Orchestration
**Value Statement:** For **Developers/Operators**, who **need to run multiple agents simultaneously**, the **Tmux Orchestration Layer** is a **tooling set** that **allows concurrent execution without window clutter**.

**The "Why":** Running 5 agents (Lisa, Ralph, etc.) in separate terminals is unmanageable. We need a "command center" view.
**Scope:**
- ✅ Script to launch/attach named tmux sessions
- ✅ `just` command integration
- ✅ Detached mode support
- ❌ Web-based management UI

**Acceptance Criteria:**
- [ ] `just flow` launches the full agent mesh in a tmux session.
- [ ] Users can toggle between agent views easily.
- [ ] Logs are preserved in detached panes.

**Attributes:**
- **Status:** 📋 Ready
- **Complexity:** Medium
- **Urgency:** Medium
- **Dependencies:** None

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
- [ ] All agents emit JSON logs to a central file/stream.
- [ ] Logs contain `agent_id`, `task_id`, and `timestamp`.
- [ ] CLI tool exists to tail/filter these logs.

**Attributes:**
- **Status:** 📋 Ready
- **Complexity:** Medium
- **Urgency:** High (Debugging)
- **Dependencies:** None

---

## EPIC-004: Agent Sandboxing
**Value Statement:** For **System Administrators**, who **fear agents destroying the host system**, the **Sandboxing Environment** is a **security boundary** that **ensures safe execution of arbitrary code**.

**The "Why":** Agents like Ralph execute code. Running this as root/user on the host is dangerous. We need containment.
**Scope:**
- ✅ Docker/Container-based execution context
- ✅ Workspace mounting strategy
- ✅ Network restriction policies
- ❌ Full VM virtualization

**Acceptance Criteria:**
- [ ] Agents run inside a defined container image.
- [ ] Agents cannot access host files outside the mounted workspace.
- [ ] Workspace state is preserved between runs.

**Attributes:**
- **Status:** 📋 Ready
- **Complexity:** High
- **Urgency:** High (Security)
- **Dependencies:** EPIC-003 (Logging)

---

## EPIC-005: Cost Control & Agent Selection
**Value Statement:** For **Budget Owners**, who **need to prevent runaway LLM costs**, the **Cost Control Middleware** is a **governance tool** that **tracks and limits token usage**.

**The "Why":** "Infinite loops" in agent logic can bankrupt us. We need a kill-switch and visibility.
**Scope:**
- ✅ Token counting middleware
- ✅ Budget limits (per session/day)
- ✅ "Cheapest Model" selection logic
- ❌ Real-time billing API integration

**Acceptance Criteria:**
- [ ] Every LLM call is logged with token count and estimated cost.
- [ ] System rejects requests when budget is exceeded.
- [ ] Reporting command shows daily spend.

**Attributes:**
- **Status:** 📋 Ready
- **Complexity:** Medium
- **Urgency:** Medium
- **Dependencies:** EPIC-003 (Logging)

---

## EPIC-006: Legacy Agent Compatibility
**Value Statement:** For **Adopters**, who **have existing agent definitions**, the **Compatibility Layer** is a **bridge** that **allows Springfield to run legacy/external agent structures**.

**The "Why":** We shouldn't force a rewrite of all existing `.github/agents` definitions. We should embrace them.
**Scope:**
- ✅ Support for `.github/agents`, `.claude/agents`, etc.
- ✅ Precedence logic (Repo > Default)
- ❌ Conversion/Migration tools

**Acceptance Criteria:**
- [ ] Springfield detects and loads agents from existing folder structures.
- [ ] Repo-defined agents override defaults.

**Attributes:**
- **Status:** 📋 Ready
- **Complexity:** Medium
- **Urgency:** Low
- **Dependencies:** None
