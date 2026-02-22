# Feature Brief: Springfield Autonomous Control Plane (EPIC-010)

**Epic:** EPIC-010
**Owner:** @Lisa (Planning) / @Ralph (Build)
**Status:** 🏗️ In Definition

---

## 1. Problem Statement
Currently, Springfield agents are "smart scripts" that manage their own execution environment (git worktrees, branches) and communicate via loose string matching ("[[FINISH]]") or file presence (`TODO.md`).

This creates fragility:
1.  **Environment Risk:** Agents can accidentally delete the branch they are running on.
2.  **Concurrency Blocks:** Lisa handles one Epic at a time; parallelizing Ralph is difficult because environment management is coupled to agent prompts.
3.  **Signal Failures:** The "Just Lisa" timeout incident proved that string scanning (`[[FINISH]]`) is unreliable for termination.
4.  **State Drift:** The Orchestrator "polls" for files, but agents "push" state changes. This disconnect leads to race conditions.

## 2. The Solution: Autonomous Control Plane
We will transform Springfield from a "Task Runner" into a **Stateful Control Plane** that manages the lifecycle, environment, and state transitions of all agents.

**Core Shifts:**
1.  **Agents as Workers:** Agents are stateless. They wake up in a pre-provisioned worktree, do a job, and signal completion. They *never* run `git worktree add` or `git checkout`.
2.  **Sentinel Signalling:** Agents use a strict, typed CLI callback (`springfield signal`) to terminate their session. This acts as a cryptographic handshake (Sentinel Token) to prove authorization.
3.  **Daemonized Orchestrator:** The `springfield` binary runs as a daemon/service, monitoring `td` state and spawning agent processes in response to state transitions.

---

## 3. The State Machine & Signaling Protocol

The Control Plane manages the `td` Epic state. Agents provide **Signals** to request transitions.

| Current State | Worker | Signal Sent (`--status`) | Next State | Orchestrator Action |
| :--- | :--- | :--- | :--- | :--- |
| **Planned** | **Lisa** | `complete` | `Ready` | 1. Parse `td` for new Epics.<br>2. Provision Worktree `feat/epic-123`.<br>3. Spawn Ralph. |
| **Ready** | **Ralph** | `success` | `Implemented` | 1. Kill Ralph process.<br>2. Spawn Bart in same worktree. |
| **Implemented** | **Bart** | `success` (Verified) | `Verified` | 1. Kill Bart.<br>2. Spawn Lovejoy. |
| **Implemented** | **Bart** | `failed` (Bugs) | `In Progress` | 1. Kill Bart.<br>2. **Re-spawn Ralph** (Correction Loop). |
| **Implemented** | **Bart** | `blocked` (Bad Arch) | `Blocked` | 1. Kill Bart.<br>2. Spawn Lisa (Replanning). |
| **Verified** | **Lovejoy** | `released` | `Done` | 1. Merge PR.<br>2. Delete Worktree.<br>3. Close Epic in `td`. |

### The "Correction Loop" (Ralph vs. Lovejoy)
When Bart signals `failed`:
*   The Orchestrator reads the signal payload (e.g., "3 unit tests failed").
*   It does **not** invoke Lovejoy.
*   It immediately re-spawns **Ralph** in the *same* worktree.
*   Ralph wakes up, sees `FEEDBACK.md` (or `td` comments), fixes the code, and signals `success` again.

---

## 4. Technical Architecture

### A. The Sentinel Protocol
Every agent invocation receives a unique, ephemeral session token via injection into their prompt.

**Invocation:**
```bash
# Orchestrator spawns Ralph
springfield agent run --name ralph --sentinel "xp9-f2a" --worktree "/tmp/wt-123"
```

**Prompt Injection (Template):**
> "Your session sentinel is `xp9-f2a`. When you are finished, you MUST execute:
> `springfield signal --sentinel xp9-f2a --status <outcome>`"

**Termination:**
When the agent runs the signal command:
1.  Springfield intercepts the command.
2.  Verifies `xp9-f2a` matches the session.
3.  Updates the internal state machine.
4.  **Terminates the agent process immediately.** (No waiting for LLM chatter).

### B. Prompt Templating
We need a templating engine (Go `text/template`) for prompts to inject dynamic values:
*   `{{.Sentinel}}` - The session token.
*   `{{.Worktree}}` - The current path.
*   `{{.EpicID}}` - The context ID.

**Location:** `.pi/prompts/ralph.prompt.md.tpl`

### C. The Daemon Loop (`orchestrate --daemon`)
1.  **Poll `td`**: "Are there any Epics in `Ready` state that don't have a runner?"
2.  **Reconcile**:
    *   Found `Ready` Epic? -> Create Worktree -> Start Ralph.
    *   Found `Implemented` Epic? -> Start Bart.
    *   Found `Verified` Epic? -> Start Lovejoy.
3.  **Monitor**: Track PIDs of running agents. Handle timeouts/crashes.

---

## 5. Acceptance Criteria (BDD)

### Scenario: Happy Path Delivery
**Given** an Epic is in `Planned` state
**When** Lisa signals `complete`
**Then** the Orchestrator provisions a worktree
**And** Ralph is spawned in that worktree
**When** Ralph signals `success`
**Then** Bart is spawned in the same worktree
**When** Bart signals `success`
**Then** Lovejoy is spawned
**When** Lovejoy signals `released`
**Then** the worktree is deleted and Epic is `Done`

### Scenario: The Correction Loop
**Given** Ralph has signaled `success`
**And** Bart is reviewing
**When** Bart signals `failed` with reason "Unit tests broken"
**Then** the Epic transitions to `In Progress`
**And** Ralph is re-spawned in the same worktree to fix it

### Scenario: Hallucination Safeguard
**Given** Ralph is running with sentinel `abc-123`
**When** Ralph executes `springfield signal --sentinel xyz-999`
**Then** the command fails with "Unauthorized Sentinel"
**And** the agent process continues (allowing Ralph to retry)

---

## 6. Dependencies
- **td(1)**: Must support concurrent access (SQLite WAL mode or careful locking).
- **Git Worktrees**: Must be managed strictly by the binary.

---
