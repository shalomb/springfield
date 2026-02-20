# ADR-008: Planning State Boundary, td(1) Adoption, and Springfield Binary Orchestration

**Date:** 2026-02-19
**Status:** Proposed

---

## Context

ADR-007 established the Just-In-Time Epic Refinement model and identified three
structural problems with the current planning loop:

1. Epic state is inferred from file presence (`TODO.md` exists?) and keyword
   grep on prose (`FEEDBACK.md` contains "critical"?) — fragile for the same
   reasons EPIC-007's retrospective identified with string matching.

2. Planning state lives in markdown files committed to git feature branches.
   When multiple Ralph worktrees run in parallel, `PLAN.md`, `TODO.md`, and
   `FEEDBACK.md` create branch contention — each worktree has its own copies,
   and no agent has a consistent view of the system's live state without
   performing git merges.

3. The Justfile is the orchestrator. It cannot be unit tested. Its decision
   logic — which agent to invoke, under what conditions, with what context —
   is encoded in shell conditionals and cannot be verified independently of
   end-to-end execution.

Three decisions are required:

- **Where does live planning state live?** (not in git-branched markdown)
- **What is the boundary between narrative documents and machine state?**
- **Which component owns the orchestration loop?**

---

## Decision

### 1. State Boundary: Narrative vs. Machine-Readable

Planning artifacts are split into two distinct categories that must not be
conflated:

**Narrative documents (markdown, git-committed to `main`):**
These exist for human understanding, agent context injection, and long-term
reasoning records. They are not parsed by the orchestrator. They change slowly.

| Document | Owner | Contains |
|:---|:---|:---|
| `PLAN.md` | Lisa | Epic value statements, the "Why", scope boundaries, ADR links, retrospective learnings, risk register prose, WSJF rationale |
| `TODO-{td-id}.md` | Lisa | Handoff context for one Epic — intent layer, approach layer, constraint layer (see ADR-007 §3). Immutable once deposited. |
| `FEEDBACK.md` | Bart | Code review narrative, quality analysis prose |
| `docs/adr/*.md` | Lisa | Architectural decision records |

**Machine-readable state (`td` — out of git, shared across all worktrees):**
These exist for the Springfield binary to drive orchestration. They are typed,
queryable, and shared across all agents regardless of which branch or worktree
they occupy.

| State | td Model | Owner |
|:---|:---|:---|
| Epic status and identity | `td epic` | Lisa |
| Task status and assignment | `td create` / `td start` / `td handoff` | Ralph |
| Task dependencies and blockers | `td dep add` | Lisa / Ralph |
| Ralph's working decomposition | `td` issues under Epic | Ralph |
| Bart's typed retrospective signal | `td log --decision` on Epic | Bart |
| Signal type (implementation failure / option viability failure / ADR verdict) | `td log --decision` structured entry | Bart |
| Ready queue for orchestrator | `td query "status = ready"` | Springfield binary |

**The PLAN.md / td link:**
Every Epic in PLAN.md carries its `td` Epic ID as a typed attribute:

```markdown
### EPIC-005: Agent Governance & Selection
**td:** `td-a3f8`
**WSJF Score: 3.25** ...
```

This is the only coupling between the two layers. The narrative references the
machine state by ID. The machine state does not reference the narrative.

### 2. TODO.md is Replaced by Named Handoff Documents

`TODO.md` is retired as a filename. It is replaced by `TODO-{td-id}.md` where
`{td-id}` is the td Epic ID for the Epic being handed off.

**Rationale:**
- A single `TODO.md` is ambiguous when multiple Ralph worktrees are running.
  `TODO-td-a3f8.md` is unambiguously scoped to one Epic.
- The filename is a machine-readable signal — the Springfield binary knows
  which context document to inject for which Epic without reading the file.
- The file's lifecycle is explicit: created by Lisa on `main`, deposited into
  the Epic's worktree branch once, immutable for the Epic's lifetime, deleted
  with the feature branch on merge.

**The three layers (per ADR-007 §3):**

```markdown
# TODO-td-a3f8.md — Agent Governance & Selection

## Intent (Immutable — from Marge's Feature Brief)
[User need, acceptance criteria, definition of done]

## Approach (Decided by Lisa at LRM — fixed for this iteration)
[Chosen option, rationale, alternatives considered and rejected]

## Constraints (Inherited — not negotiable)
[ADR links, Farley/Adzic thresholds, ACP requirements,
 tech debt items that must not be worsened]
```

Ralph's working task decomposition is **not** in this file. It lives in `td`.

### 3. td(1) is Adopted as the Planning State Store

`td` (https://github.com/marcus/td) is adopted as the shared, machine-readable
planning state store for all agents.

**Why td over alternatives:**

- **SQLite off-branch:** `.todos/issues.db` lives under `.todos/` in the project
  root, shared across all git worktrees on the same host. Branch contention on
  planning state is eliminated — no merges, no conflicts, no reconciliation logic.
- **Native git worktree support (confirmed from source):** `td`'s `workdir`
  package explicitly resolves the database root for git worktrees via
  `git rev-parse --git-common-dir`. When Ralph runs in
  `worktrees/epic-td-a3f8`, td automatically resolves up to the main
  worktree's `.todos/issues.db` — no `.td-root` file or manual configuration
  required. All Ralph worktrees share one database transparently. This was
  verified by reading `internal/workdir/workdir.go` in the td source.
- **Adequate capability:** Epic tracking, dependency graph, `td query` language,
  session-aware handoffs (`td handoff --done --remaining --decision --uncertain`),
  `td usage` for cold-start context injection. Covers all Springfield needs today.
- **CLI only — subprocess model confirmed:** `td`'s meaningful packages
  (`internal/db`, `internal/models`, `internal/session`, `internal/query`,
  `internal/workflow`) are all under `internal/` — not importable by external
  Go packages by compiler enforcement. The Springfield binary must shell out
  to `td` as a subprocess. This is acceptable — the dependency is operational,
  not compile-time — but it means `td` must be present in the runtime
  environment and integration tests require a real `td` instance against a
  test-scoped working directory.
- **Simplicity over power:** `td` is chosen over `beads` (which offers
  Dolt-backed distributed merging, cell-level conflict resolution, and richer
  graph semantics) because Springfield's current execution model is single-host
  worktrees. The additional capability of `beads` is not yet needed, and the
  Dolt dependency is a larger operational bet than SQLite. This decision is
  revisited if Springfield scales to distributed agents on separate machines.
- **Go + SQLite stack:** Consistent with Springfield's existing technology choices.
- **LRM principle applied to tooling:** Adopt the adequate tool now; upgrade when
  the constraint that requires upgrade actually exists.

**The branch contention problem solved:**

```
main / coordination
├── feat/epic-td-a3f8   ← Ralph-1 worktree (code only)
├── feat/epic-td-b2c1   ← Ralph-2 worktree (code only)
│
.todos/issues.db         ← ALL agents, ALL worktrees, one shared store
                            No branching. No merging. No conflicts.
                            Resolved automatically by td's workdir package.
```

**td as the orchestration signal source:**

| Old signal | Fragility | New signal | Strength |
|:---|:---|:---|:---|
| `TODO.md` exists | File presence on one branch | `td query "status = ready"` | Typed, cross-branch |
| `FEEDBACK.md` contains "critical" | Keyword grep on prose | `td log --decision` signal type field | Typed enumeration |
| Iteration counter `max_iterations=2` | Magic number, no semantics | Epic status transitions | Explicit state machine |
| Lisa runs unconditionally | No trigger condition | `td query "status = verified"` → Lisa's LRM | Event-driven |

### 4. The Springfield Binary Owns the Orchestration Loop

The Justfile `do` loop is retired as the orchestrator. The Springfield binary
(`cmd/springfield`) takes ownership of the orchestration state machine.

**Rationale:**
- Shell conditionals cannot be unit tested. The orchestration logic — which
  agent to invoke, under what Epic state, with what signal — must be testable
  in isolation.
- Epic states become a typed Go enum. Invalid transitions are compile errors,
  not runtime surprises caught in production.
- The state machine is auditable, documented, and covered by the same TDD
  discipline Ralph applies to production code.

**Epic state enumeration:**

```go
type EpicStatus string

const (
    StatusPlanned     EpicStatus = "planned"
    StatusReady       EpicStatus = "ready"
    StatusInProgress  EpicStatus = "in_progress"
    StatusImplemented EpicStatus = "implemented"
    StatusVerified    EpicStatus = "verified"
    StatusDone        EpicStatus = "done"
    StatusDeferred    EpicStatus = "deferred"
    StatusBlocked     EpicStatus = "blocked"
)
```

**State transition table (the `do` loop, typed):**

| Current State | Signal | Next State | Agent Invoked |
|:---|:---|:---|:---|
| `planned` | Lisa's LRM decision | `ready` | — (Lisa deposits TODO-{id}.md) |
| `ready` | Springfield binary tick | `in_progress` | Ralph |
| `in_progress` | Ralph completes | `implemented` | Bart |
| `implemented` | Bart: no critical issues | `verified` | Lovejoy |
| `implemented` | Bart: implementation failure | `in_progress` | Ralph (iterate) |
| `implemented` | Bart: option viability failure | `blocked` | Lisa (re-evaluate) |
| `implemented` | Bart: ADR invalidated | `blocked` | Lisa (re-evaluate + update ADR) |
| `verified` | Lovejoy merges | `done` | — |
| `blocked` | Lisa re-decides | `ready` | Ralph (new approach) |
| `any` | Lisa: defer | `deferred` | — (moves to debt register) |

**The Justfile becomes a thin invocation layer:**

```bash
# just do — delegates entirely to the binary
do *args:
    ./bin/springfield orchestrate {{args}}
```

The binary reads td for state, invokes agents, manages worktrees, and drives
transitions. All decision logic is in Go, not shell.

### 5. Lisa's Execution Context and the Worktree Deposit Protocol

Lisa operates on `main` (or `coordination` in protected branch environments).
She never executes planning work on a feature branch.

**The handoff sequence for each Epic:**

```
1. Lisa on main:
   a. Runs ToT option exploration, makes LRM decision
   b. Drafts ADR if needed (committed to main, status: Proposed)
   c. Creates td Epic:
        td epic create "Agent Governance" --priority P1
        → returns td-a3f8
   d. Creates initial td tasks from acceptance criteria
   e. Writes TODO-td-a3f8.md (narrative handoff — intent, approach, constraints)

2. Lisa creates the execution environment:
        git worktree add worktrees/epic-td-a3f8 -b feat/epic-td-a3f8

3. Lisa deposits the handoff into the worktree:
        cp TODO-td-a3f8.md worktrees/epic-td-a3f8/
        cd worktrees/epic-td-a3f8
        git add TODO-td-a3f8.md
        git commit -m "plan(epic-td-a3f8): deposit handoff context for Ralph"
        git push

4. Lisa updates td Epic status:
        td update td-a3f8 --status ready

5. Springfield binary detects td-a3f8 status=ready:
        → Creates worktree if not exists
        → Invokes Ralph in worktrees/epic-td-a3f8
        → Injects @TODO-td-a3f8.md as context
        → Ralph runs td usage to load live task state
```

**Agent branch awareness:**

| Agent | Branch | td access | Writes to |
|:---|:---|:---|:---|
| Lisa | `main` / `coordination` | Full — plans and creates | PLAN.md, ADRs, td Epics, TODO-{id}.md (then deposits) |
| Ralph | `feat/epic-{td-id}` | Full — reads and updates | Code, tests, td tasks, td handoffs |
| Bart | `feat/epic-{td-id}` | Full — reads and signals | FEEDBACK.md, td `--decision` logs |
| Lovejoy | `main` (via squash merge) | Full — closes | PLAN.md retrospective, td Epic → `done` |

---

## Consequences

### Positive

- **Branch contention eliminated.** Planning state in td is never on a feature
  branch. Multiple Ralph worktrees run concurrently without planning state
  conflicts.

- **Orchestration is testable.** The state machine in the Springfield binary is
  covered by table-driven unit tests. The `do` loop's behaviour is verifiable
  without end-to-end execution.

- **Typed signals replace keyword grep.** Bart's signal type (implementation
  failure, option viability failure, ADR verdict) is a typed field in td, not
  a keyword matched in prose. The orchestrator cannot misclassify a signal.

- **TODO naming is unambiguous.** `TODO-td-a3f8.md` is self-identifying.
  Multiple simultaneous Epics have no naming collision. The Springfield binary
  knows which context document to inject without filesystem search.

- **PLAN.md becomes a better document.** Stripped of execution tracking (status
  emoji, checkbox task lists, corrective action checklists), it becomes a pure
  reasoning record — value statements, retrospectives, risk register, ADR links.
  It reads as an argument, not a half-executed plan.

- **Ralph's cold-start is reliable.** `td usage` gives any new Ralph session
  complete situational awareness — current focus, pending tasks, recent
  decisions, blockers — without depending on the state of a markdown file that
  may be stale from the last session.

### Negative

- **td is an external dependency.** The Springfield binary must shell out to
  `td` or integrate its SQLite schema directly. The former is simpler but adds
  process overhead. The latter couples Springfield to td's schema.

- **Single-host constraint.** `.todos/issues.db` is local. Distributed agents
  on separate machines cannot share it. This is acceptable today; it is a
  known scaling constraint.

- **Lisa must maintain two representations.** PLAN.md and td must be kept
  consistent for each Epic. The td ID in PLAN.md is the bridge — if it drifts,
  the link breaks. Tooling or convention must enforce consistency.

- **Migration cost.** Existing PLAN.md Epics must be seeded into td. Existing
  `TODO.md` must be renamed and restructured. This is a one-time cost but
  requires a dedicated migration Epic.

---

## What This Does Not Decide

- Whether the Springfield binary shells out to `td` or integrates its SQLite
  schema directly (implementation detail of the orchestration Epic).
- The precise `td handoff` format for Bart's retrospective signal to Lisa
  (specified in the Epic Decomposition Protocol).
- How Lisa manages multiple simultaneous option sets in td when parallel Epics
  are in flight.
- Migration path for existing PLAN.md Epics into td (dedicated migration Epic).
- Upgrade path from `td` to `beads` if distributed execution becomes a
  requirement (revisited when the constraint exists).

---

## References

- ADR-007 — Epic Refinement and Lisa's LRM Role (motivation for this ADR)
- `docs/standards/git-branching.md` — Worktree and branch lifecycle
- `docs/concepts/model.md` — Discovery/Delivery Diamond, Ralph Wiggum Loop
- `docs/reference/loop-catalog.md` — Plan-and-Execute, Manager-Worker Loop
- `td(1)` — https://github.com/marcus/td
- `beads` — https://github.com/steveyegge/beads (evaluated, deferred)
- EPIC-005 (Agent Governance) — prerequisite for budget-governed multi-agent
  orchestration
