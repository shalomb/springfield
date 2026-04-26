# Pipeline Redesign: Deterministic Control Plane + Focused Workers

**Status**: Concept / Pre-RFC  
**Created**: 2026-04-26  
**Last updated**: 2026-04-26  
**Session context**: Emerged from two parallel threads — (1) operational
experience running agent-mux against the terrapyne TODO.md backlog (~20
tasks, ~15 merged PRs, one human orchestrator watching a tmux pane), and
(2) a deep read of the Springfield codebase and its ADRs, particularly
ADR-007 (Epic Refinement, Lisa's LRM Role) and ADR-008 (Planning State
Boundary, td adoption, Springfield binary orchestration).

**Resume from**: Section 6 — Implementation Priorities. Everything up to
and including Section 5 is analysis. Section 6 is where work begins.

---

## 1. The Core Thesis

An LLM is a poor supervisor. It is an excellent worker.

The failure mode we repeatedly hit: an LLM orchestrator compensating for
missing infrastructure — polling panes for shell prompts, guessing whether
agents are alive, managing state transitions with string parsing, timing out
agents that were actually running fine.

None of that required intelligence. All of it required reliability.

**The redesign separates two concerns that are currently fused:**

- **Control plane**: deterministic, crash-safe, no LLM — owns state
  transitions, job routing, process lifecycle, stall detection
- **Data plane**: LLM workers that are pure functions with a narrow
  input/output contract — own only generative reasoning

This is not a new idea. It is how every production job queue, CI pipeline,
and data processing system is built. The fact that the workers happen to be
LLMs does not change the control plane architecture at all.

---

## 2. Evidence Base: Observed Frictions

From the terrapyne agent-mux session. These are infrastructure failures,
not LLM failures.

| Failure | Root cause | Category |
|---|---|---|
| Orchestrator sends Bart while Ralph still runs | tmux send-keys is fire-and-forget — no backpressure | Transport |
| `pane_is_ready()` times out, sends anyway | Monitor UI ≠ shell prompt — pane state undetectable | Sensing |
| Bart JSONL never created | Command queued behind running agent in same pane | Resource contention |
| Poller loops 200+ times on completed file | `sys.exit()` inside `for` loop — outer `while True` restarts | Logic bug |
| Timeout kills agents that are running fine | Wall-clock timeout ≠ stall; long agents are valid | Wrong signal |
| Dispatch state diverges from TODO.md | Two sources of truth updated independently | Data sync |
| Orchestrator crash loses in-flight state | Stateful Python subprocess, no persistence | Durability |
| wave3a-queue.sh and orchestrator.py collide | Two processes, one resource, no mutex | Concurrency |
| Ralph re-derives conventions every session | Stable knowledge re-computed per-task | Context bloat |
| Ralph investigates pre-existing failures | Baseline not communicated | Missing input |
| A10 / D7 / C1 timeout despite running fine | Timeout based on wall-clock, not output growth | Wrong signal |

**The constraint is agent API capacity, not the pane.** The tmux pane is a
display surface, not a job queue. Treating it as a job queue was the
architectural error. Every friction above flows from that mistake.

---

## 3. Springfield Convergence: What Already Exists

Before proposing changes, what Springfield has already independently built
that validates the thesis:

**Already correct:**
- **Sentinel protocol** — agents terminate via `springfield signal --sentinel
  <uuid>`. Cryptographic handshake, not file-watching or string parsing.
  Authorises the signal (right agent, right session) and terminates the
  process. This is strictly better than the file-drop approach in the original
  concept — adopt it.
- **Deterministic Go orchestrator** — `orchestrator.Tick()` is a pure function.
  Queries `td`, checks state, spawns processes, advances state machine. Zero
  LLM involvement.
- **`td` as the job store** — SQLite off-branch, shared across all worktrees,
  queryable, typed. Survives crashes. Supports `td query`, `td log --decision`,
  `td handoff`. Solves branch contention entirely.
- **Stateless workers in pre-provisioned worktrees** — `WorktreeManager.
  DepositHandoff()` copies `TODO-{id}.md` into the worktree before Ralph
  starts. Ralph wakes in a clean environment with context already placed.
- **`td usage` for cold-start context** — Ralph runs `td usage` at session
  start to get current task state, recent decisions, blockers. No markdown
  file archaeology.
- **ADR-007: LRM model and retrospective signal** — fully specified in the
  ADR. Lisa uses ToT + Self-Consistency before committing to an approach.
  Bart produces a structured retrospective for Lisa. The feedback loop is
  architecturally designed. (It is not yet implemented in the prompts —
  see Section 4.)
- **Feedback persistence via `td log --decision`** — Bart's learnings survive
  across agent sessions. Lisa queries with a watermark to read only what's
  new. This is stronger than a flat FEEDBACK.md accumulator.

**What Springfield has that the original concept doc missed:**
The Planner stage I described exists — it is Lisa's job, as specified in
ADR-007. The scope manifest, test scaffold, and guardrail assertions are
the *content* Lisa should be producing in `TODO-{id}.md`. The architecture
supports it. The prompts don't ask for it yet.

---

## 4. The Agent Distribution Problem

### 4.1 The two-loop structure

Springfield's five agents operate at **two different loop levels** that are
currently conflated:

```
Outer loop (per release — milestone boundary):
  Marge ──────────────────────────────────────── Lovejoy
  "define the release scope"          "ship the release"
        │                                        ▲
        ▼                                        │
  Inner loop (per epic — delivery boundary):     │
    Lisa → Ralph → Bart → (repeat per epic) ─────┘
```

**Marge** fires once at the start of a Springfield invocation to define
release scope and produce Feature Briefs. She does not fire per-epic.

**Lovejoy** fires once at the end of a Springfield invocation when all
epics in the release are `done`. He does not fire per-epic. A release is
the *aggregate* of all epics in a PLAN.md milestone — the natural unit
for a CHANGELOG entry, semver bump, and release notes.

**Lisa, Ralph, Bart** form the inner loop — one cycle per epic, with
Ralph→Bart retry on rejection.

This reframe eliminates:
- Marge's per-epic merge gate (was compensating for Bart lacking product
  context — give Bart the Feature Brief and it's unnecessary)
- Lovejoy's per-epic trigger (changelog for one bug fix is absurd;
  changelog for a milestone is meaningful)

### 4.2 What each agent actually does

| Agent | Loop level | Cognitive work | Loop assigned |
|---|---|---|---|
| Marge | Per release (start) | Problem framing, option exploration, Feature Brief authorship, BDD scenario quality | ToT, GECR, Self-Consistency |
| Lisa | Per epic | Option generation, LRM decision, handoff contract production | ToT, Self-Consistency, PAE, OLA |
| Ralph | Per epic (retry on rejection) | TDD execution against contract | Ralph Wiggum, ReAct, OHECI (escalation) |
| Bart | Per epic | Quality review, product conformance, retrospective for Lisa | OHECI, TALAR, RDS |
| Lovejoy | Per release (end) | Release narrative, CHANGELOG, semver, tag, cleanup | GECR, RDS |

### 4.3 Is Lovejoy necessary?

At per-epic granularity: no — a Go function in the orchestrator can merge,
tag, and update CHANGELOG from `td log` entries deterministically. No LLM
needed.

At per-release granularity: yes — synthesising a coherent human-readable
narrative across a whole milestone requires generative reasoning. GECR
(generate draft → critique → refine) produces a better CHANGELOG than a
bullet list generated from commit messages.

**Decision**: Lovejoy fires per-release, not per-epic. Per-epic merge is
handled by the orchestrator (deterministic: `gh pr merge --squash --auto`
after Bart signals `success`).

### 4.4 Agent count assessment

Five agents, two loop levels. No agent is redundant *at the right
granularity*. The smell was per-epic Lovejoy and per-epic Marge merge
gate — wrong granularity, not wrong agents.

---

## 5. The Loop Assignment Gap

### 5.1 What each agent's prompt actually wires

Cross-referencing the loop catalog against the prompt templates:

| Agent | Loops called for | Loops actually wired | Gap |
|---|---|---|---|
| Marge | ToT, GECR, Self-Consistency, VS | SPA (read → decide → write) | **Critical** — first-pass Feature Briefs only |
| Lisa | ToT, Self-Consistency, OLA, PAE | Linear PAE only | **Critical** — LRM model from ADR-007 unimplemented |
| Ralph | Ralph Wiggum, ReAct, OHECI (escalation) | Ralph Wiggum implicit, ReAct implicit | **Medium** — OHECI escalation path missing |
| Bart | OHECI, TALAR, RDS | Static review + SPA | **Critical** — retrospective signal for Lisa absent |
| Lovejoy | GECR, RDS | Mechanical checklist (SPA) | **Low** — acceptable at per-epic; matters at per-release |

### 5.2 The three critical missing connections

**Connection 1: Bart → Lisa retrospective signal (highest priority)**

This is the most consequential gap in the whole system. ADR-007 is built
entirely around this feedback loop. Without it, every Lisa session starts
cold. She reads prose in FEEDBACK.md but receives no *typed* signal about
what the previous iteration revealed about the next epic.

What Bart must produce (in addition to existing FEEDBACK.md):

```markdown
## Retrospective Signal (structured — for Lisa's LRM decision)

next_epic_constraints:
  - "StateVersionsAPI.list() has a round-trip that will break if
     filter[workspace][id] is added to F3 workspace scope"

approach_layer_verdict: confirmed | partially_confirmed | invalidated
  # confirmed = Lisa's chosen approach held under execution
  # partially_confirmed = held with caveats (document them)
  # invalidated = fundamental assumption broke — Lisa must re-enter ToT

adr_verdict: confirmed | rejected | amended
  # which ADR was exercised, what the evidence says

learnings_for_next_epic:
  - "paginate_with_meta mock gap: always mock in conftest fixtures"
  - "worktree-locked gh pr merge: use --repo flag, omit --delete-branch"

tech_debt_signals:
  - severity: medium
    item: "no multi-page test for paginate() delegation path"
    suggested_epic: "test coverage sweep"
```

This goes into `td log <epic-id> --decision` so Lisa can query it. It is
not committed to the feature branch. It survives Lovejoy's cleanup.

**Connection 2: Lisa's ToT + Self-Consistency at LRM (highest priority)**

ADR-007 specifies this explicitly. Lisa's prompt skips it entirely. She
goes directly from "identify next epic" to "create tasks and deposit
handoff." The option generation + Self-Consistency validation step that
ADR-007 calls her primary value-add is absent.

What Lisa must do before committing to an approach:

```
1. Generate 2-3 candidate approaches for the epic
   - For each: evaluate against ADR constraints, Farley, Adzic, tech debt
     signals from td log (Bart's retrospective, watermarked)
   - Record alternatives considered and why each was accepted/rejected

2. Self-Consistency validation on the chosen option:
   - Run the constraint evaluation twice independently
   - If the conclusion is not stable (different evaluation each time),
     the option has hidden ambiguity — do not commit, re-examine

3. Only then: produce TODO-{id}.md with the three layers
   - Intent layer (from Marge's Feature Brief — immutable)
   - Approach layer (Lisa's LRM decision — fixed for this iteration)
   - Constraint layer (ADR links, Farley/Adzic thresholds, guardrails)

4. Ralph's working layer is NOT pre-written by Lisa
   - Ralph decomposes bottom-up via TDD
   - Lisa's TODO-{id}.md sets the contract; Ralph fills the working layer
```

**Connection 3: Marge's ToT/GECR in Discovery (medium priority)**

Marge currently produces a Feature Brief on first-pass interpretation.
The exploratory reasoning that validates the framing — is this the right
problem? are we solving the root cause? — is absent.

What Marge must do before finalising the Feature Brief:

```
1. Generate 2-3 problem framings (not solutions — framings)
   - What is the user's actual pain?
   - What would solving vs. not solving look like?
   - What adjacent problems might this be confused with?

2. Score each framing against: user need, roadmap fit, risk acknowledgment

3. GECR: draft the Feature Brief, critique it against the winning framing,
   refine, one more pass

4. Log marge_approved to td — this becomes the immutable Intent layer
   for Lisa and Ralph
```

### 5.3 What does NOT need changing

- **Ralph's core loop**: Ralph Wiggum + ReAct is correct. The stateless
  worktree model is the right answer. The TDD discipline is the right
  answer. The `td usage` cold-start is the right answer.
- **Bart's quality review**: OHECI for edge cases, adversarial testing — 
  this is working. The gap is the *additional* structured output for Lisa.
- **The sentinel protocol**: correct and should be preserved.
- **`td` as state store**: correct. Survives crashes, cross-worktree, typed.
- **The orchestrator as Go binary**: correct. Testable, typed state machine.

---

## 6. Implementation Priorities

Ordered by impact. Each item is independent enough to be a standalone epic.

### Priority 1 — Bart's retrospective signal (unblocks everything else)

**What**: Add structured `## Retrospective Signal` section to Bart's
prompt. Output goes to `td log --decision`, not FEEDBACK.md.

**Why first**: ADR-007's entire LRM model depends on this. Lisa cannot
run ToT with last-iteration's information until Bart produces it. This
is the single highest-leverage change in the system — it activates the
Bart→Lisa feedback loop that makes the whole pipeline adaptive.

**Effort**: S. One prompt section, one `td log` command, one schema.

**Files to touch**:
- `.pi/prompts/bart.prompt.md.tpl` — add retrospective signal section
- `.pi/agents/bart.md` — document the dual-output responsibility
- `docs/standards/feedback.md` — add retrospective signal schema
- `tests/integration/features/automated_feedback_loop.feature` — BDD spec

---

### Priority 2 — Lisa's ToT + Self-Consistency at LRM

**What**: Add option generation, constraint evaluation, and Self-Consistency
validation to Lisa's prompt before she commits to an approach.

**Why second**: Requires Bart's retrospective signal to be useful (Lisa
reads `td log` for learnings). With Priority 1 done, Lisa can run properly
informed option evaluation. Without Priority 1, Lisa's ToT is uninformed.

**Effort**: M. Prompt restructure, new section, `td query` for watermark.

**Files to touch**:
- `.pi/prompts/lisa.prompt.md.tpl` — restructure workflow to add ToT phase
- `.pi/agents/lisa.md` — update responsibilities, remove "task-decomposer"
- `docs/adr/ADR-007-epic-refinement-and-lisa-lrm.md` — mark as Accepted
  once implementation matches spec
- `tests/integration/features/control_plane.feature` — BDD for LRM flow

---

### Priority 3 — Lisa's handoff contract tightening

**What**: Extend `TODO-{id}.md` to include machine-readable scope manifest,
test scaffold, baseline snapshot, and guardrail script — not just narrative
prose. Ralph reads structured facts, not prose he must interpret.

**Why third**: Requires Priority 2 (Lisa's ToT) to be in place, because
the scope manifest is produced during the LRM decision, not independently.

**Effort**: M. Template extension, baseline capture script, guardrail
script generation.

**What TODO-{id}.md gains**:

```toml
# scope.toml — machine-readable, generated by Lisa during LRM
[scope]
modify  = ["src/terrapyne/core/backend.py"]
create  = ["tests/test_core/test_backend_cloud.py"]
avoid   = ["src/terrapyne/cli/*"]
hubs    = []   # files imported by 3+ others — treat with care

[baseline]
test_count     = 634
known_failures = []
known_warnings = ["DeprecationWarning in conftest.py:47"]

[guardrails]
script = "guardrails-{id}.sh"   # generated, Ralph runs before signalling
```

Plus the test scaffold (skeleton test functions Ralph must make pass).

**Files to touch**:
- `.pi/prompts/lisa.prompt.md.tpl` — add scope manifest generation step
- `internal/orchestrator/worktree.go` — DepositHandoff reads scope.toml
- `docs/adr/ADR-007` — document the three-layer + machine-readable layer

---

### Priority 4 — Marge's ToT/GECR in Discovery

**What**: Add option-exploration loop to Marge's Feature Brief production.
Generate problem framings, score, GECR on the Brief, then finalise.

**Why fourth**: Improves the quality of input to Lisa's LRM decision, but
the inner loop (Lisa → Ralph → Bart) can function without it. The cost of
a wrong Feature Brief is rework in a future epic, not a broken pipeline.

**Effort**: S–M. Prompt addition, no structural changes.

**Files to touch**:
- `.pi/prompts/marge.prompt.md.tpl` — add framing exploration section
- `.pi/agents/marge.md` — document the GECR workflow

---

### Priority 5 — Ralph's OHECI escalation path

**What**: When Ralph hits an unexpected failure that doesn't resolve in 2
attempts, give him a structured reasoning loop: observe the failure mode,
hypothesize whether it's an implementation problem or an assumption break,
run one targeted experiment, conclude. If assumption break: escalate via
`springfield signal --status blocked` immediately rather than continuing.

**Why fifth**: Ralph's current behaviour (investigate inline, burn context)
is functional but wasteful. OHECI as an explicit escalation path prevents
context spiral and surfaces assumption breaks to Lisa faster.

**Effort**: S. One section in Ralph's prompt.

**Files to touch**:
- `.pi/prompts/ralph.prompt.md.tpl` — add OHECI escalation section
- `.pi/agents/ralph.md` — document when to escalate vs. iterate

---

### Priority 6 — Lovejoy at release boundary

**What**: Change Lovejoy's trigger from per-epic (`bart_approved`) to
per-release (all epics in PLAN.md milestone reach `done`). Add GECR loop
for release narrative generation. Add RDS for retrospective capture.

**Why sixth**: The per-epic merge is handled deterministically by the
orchestrator (`gh pr merge --squash --auto`). Lovejoy's value is the
human-readable release story, which only exists at milestone boundary.

**Effort**: M. Orchestrator trigger change, prompt restructure.

**Files to touch**:
- `internal/orchestrator/orchestrator.go` — change Lovejoy trigger
- `internal/orchestrator/status.go` — add `StatusAllEpicsDone` check
- `.pi/prompts/lovejoy.prompt.md.tpl` — restructure for release scope
- `tests/integration/features/control_plane.feature` — BDD for release

---

### Priority 7 — Deterministic control plane (agent-mux specific)

**What**: For the agent-mux / terrapyne workflow (outside Springfield),
replace the Python orchestrator with a file-watching supervisor that
reads job tickets. Two slots: Ralph-slot and Bart-slot, concurrent.
State in sqlite, not `/tmp/` JSON files.

**Why seventh**: Springfield already has this solved with the Go
orchestrator + td + sentinel. The agent-mux workflow is a simpler
variant of the same problem. This priority applies only if agent-mux
continues to be used independently of Springfield.

**Effort**: M. ~200 lines Python supervisor, sqlite schema, inotify loop.

**Files to create** (in agent-mux skill):
- `scripts/supervisor.py`
- `scripts/jobs.db` (schema)
- Updates to `dispatch.py` (strip tmux/pane code, keep job factory)
- Updates to Ralph/Bart prompt templates (exit contract section)

---

## 7. The Redesigned Architecture (Full Picture)

```
Springfield invocation:
═══════════════════════════════════════════════════════════════════

OUTER LOOP (per release)
─────────────────────────────────────────────────────────────────

  PLAN.md "v0.7.0: Autonomous Control Plane"
       │
       ▼
  [Marge — fires ONCE at invocation start]
    Loop: ToT (problem framings) → GECR (Feature Brief)
    Produces: Feature Briefs in PLAN.md, marge_approved in td
    Signals: springfield signal --status success
       │
       ▼
  INNER LOOP (per epic — repeats for each epic in release)
  ─────────────────────────────────────────────────────────

    [Lisa — fires once per epic]
      Loop: OLA (query td for Bart's retrospective signals, watermarked)
            → ToT (2–3 candidate approaches)
            → Self-Consistency (validate winner, 2 independent passes)
            → PAE (produce three-layer TODO-{id}.md + scope.toml)
      Produces: TODO-{id}.md (intent + approach + constraints + scaffold)
                scope.toml (machine-readable manifest)
                guardrails-{id}.sh (assertions Ralph runs)
                td Epic → status: ready
      Signals: springfield signal --status complete
          │
          ▼
    [Orchestrator — deterministic, no LLM]
      Detects status=ready → provisions worktree → deposits handoff
      Spawns Ralph with sentinel
          │
          ▼
    [Ralph — fires once per epic, retry on Bart rejection]
      Loop: Ralph Wiggum (stateless TDD per task)
              Inner: ReAct (thought → action → observation)
              Escalation: OHECI (if unexpected failure after 2 attempts
                          → observe → hypothesize → experiment → conclude
                          → if assumption break: signal blocked, not failed)
      Produces: commits, PR, sentinel signal
      Exit contract: writes bart-{id}.json OR failed-{id}.json
                     (job ticket — supervisor routes it)
          │
          ▼ (job ticket appears)
          │
    [Orchestrator — deterministic]
      Reads job ticket → spawns Bart with sentinel
      (Ralph-slot is now free — next epic's Ralph can start here
       if available — Ralph and Bart run concurrently)
          │
          ▼
    [Bart — fires once per epic, may trigger Ralph retry]
      Loop: OHECI (adversarial hypothesis → experiment → conclude)
            TALAR (run tests → analyze → learn → adjust → retest)
            RDS (reflect → document retrospective → write to td log)
      Reviews: code quality, product conformance (vs. Feature Brief),
               ADR validity
      Produces: FEEDBACK.md (for Ralph if rejected)
                td log --decision (structured retrospective for Lisa)
                sentinel signal (success/failed/blocked)
      Exit contract: writes approved-{id}.json OR rejected-{id}.json

    [Orchestrator on approved]:
      gh pr merge --squash --auto (deterministic, no agent)
      td Epic → status: done
      [if more epics remain: loop back to Lisa for next epic]

    [Orchestrator on rejected]:
      Ralph re-queued with FEEDBACK.md context (same worktree, attempt++)
      [if max_retries exceeded: td Epic → blocked → surface to human]

    [Orchestrator on blocked (Lisa chose wrong option)]:
      td Epic → blocked
      Lisa re-invoked with new constraint information
      Lisa re-enters ToT with failed option recorded as constraint

  END INNER LOOP (all epics done)
  ─────────────────────────────────────────────────────────

       │
       ▼
  [Lovejoy — fires ONCE when all epics in release are done]
    Loop: GECR (generate release narrative → critique → refine)
          RDS (capture release learnings → write to PLAN.md retrospective)
    Produces: CHANGELOG.md entry (human-readable milestone narrative)
              semver tag, git tag, release publication
              PLAN.md retrospective section
              td Epic statuses → closed
    Signals: springfield signal --status released

END OUTER LOOP
═══════════════════════════════════════════════════════════════════
```

---

## 8. Handoff Contract: What Ralph Receives

The single most important change for Ralph's efficiency. Currently Ralph
receives a narrative brief and derives the contract. Under the redesign,
Lisa produces the contract and Ralph executes it.

**`TODO-{id}.md` (three layers — unchanged from ADR-007):**

```markdown
# TODO-td-a3f8.md — Agent Governance & Selection

## Intent (immutable — from Marge's Feature Brief)
User need, acceptance criteria, definition of done.
Ralph cannot renegotiate this.

## Approach (decided by Lisa at LRM — fixed for this iteration)
Chosen option, rationale, alternatives considered and rejected,
ADR links (status: Proposed until Bart confirms).
Ralph cannot renegotiate this. If he discovers it is wrong,
he signals blocked, not failed.

## Constraints (inherited — not negotiable)
ADR links, Farley/Adzic thresholds, ACP requirements,
tech debt items that must not be worsened.

## Working Layer (Ralph's — fully owned by him)
[Empty at deposit. Ralph fills this via TDD bottom-up.]
His failing tests are his task list.
```

**`scope.toml` (machine-readable — new):**

```toml
[scope]
modify  = ["internal/agent/agent.go", "internal/config/config.go"]
create  = ["internal/agent/governance_test.go"]
avoid   = ["cmd/springfield/main.go", "pkg/logger/"]
hubs    = ["internal/agent/agent.go"]  # imported by 4 others — careful

[baseline]
test_count     = 47
test_command   = "just test"
known_failures = []
known_warnings = []

[guardrails]
script = "guardrails-td-a3f8.sh"
# Ralph runs this before writing the job ticket.
# Non-zero exit blocks the ticket. Zero exit = safe to signal.
```

**`guardrails-{id}.sh` (generated by Lisa):**

```bash
#!/bin/bash
# Assertions Ralph runs before signalling success.
# Non-zero exit = do not write job ticket.

# Scope: only allowed files were touched
TOUCHED=$(git diff --name-only origin/main)
echo "$TOUCHED" | grep -q "cmd/springfield/main.go" \
  && { echo "FAIL: main.go is out of scope"; exit 1; }

# Baseline: test count did not decrease
COUNT=$(just test 2>&1 | grep -oP '\d+ passed' | grep -oP '\d+' | tail -1)
[ "${COUNT:-0}" -ge 47 ] \
  || { echo "FAIL: test count dropped to $COUNT (baseline: 47)"; exit 1; }

# Commit format: conventional commits
git log -1 --format="%s" | grep -qE "^(feat|fix|refactor|test|docs)\(" \
  || { echo "FAIL: commit message not conventional"; exit 1; }

echo "GUARDRAILS PASSED"
```

**Ralph's context comparison:**

| | Current | With Planner |
|---|---|---|
| Conventions (prose) | ~500 tokens | 0 — in guardrails script |
| Task description (prose) | ~300 tokens | ~100 tokens (intent layer) |
| Approach (to derive) | implicit, ~200 tokens exploration | ~50 tokens (approach layer) |
| Where to put things (to derive) | ~300 tokens exploration | ~30 tokens (scope.toml) |
| What test to write (to derive) | ~400 tokens exploration | ~100 tokens (scaffold) |
| Baseline (unknown) | ~200 tokens investigation | 20 tokens (baseline section) |
| Mechanics (git, PR, signal) | ~200 tokens | ~50 tokens |
| **Total** | **~2200 tokens + exploration** | **~350 tokens** |

6× context reduction. More importantly: the exploration work is gone.
Ralph executes; he does not investigate.

---

## 9. Bart's Dual Output

Bart currently produces one output: `FEEDBACK.md` for Ralph.
Under the redesign he produces two: `FEEDBACK.md` for Ralph, and a
structured retrospective signal for Lisa via `td log --decision`.

The retrospective signal schema:

```markdown
## Retrospective Signal

### Approach Layer Verdict
verdict: confirmed | partially_confirmed | invalidated
notes: >
  "The workspace_id filter approach was confirmed. The StateVersionsAPI
  correctly avoided the round-trip when workspace_id was provided directly."

### ADR Verdict  
adr_id: ADR-008
verdict: confirmed | rejected | amended
notes: >
  "td as shared state store survived parallel worktree access. Confirmed."

### Constraints Revealed (for next epic's LRM)
- "paginate_with_meta mock gap: always mock in conftest — affects any epic
   that touches TFCClient"
- "worktree-locked gh pr merge: use --repo flag, omit --delete-branch"

### Tech Debt Signals
- severity: minor
  item: "no multi-page test for paginate() delegation path"
  action: backlog

### Next Epic Observations
- "F3 workspace create/delete will touch workspace_cmd.py — hub file.
   Lisa should scope carefully to avoid collision with C2."
```

This goes to `td log <epic-id> "retrospective" --decision`.

Lisa's session starts with:
```bash
td query "type = epic AND status = done AND closed >= <watermark>" \
  --output json | jq '.[] | .logs[] | select(.type == "decision")'
```

She reads typed signals, not prose she has to interpret.

---

## 10. The Sentinel vs. File-Drop Question

The original concept doc proposed workers write job ticket files to
`/tmp/jobs/`. Springfield uses `springfield signal --sentinel <uuid>`.

**Springfield's approach is better for three reasons:**

1. **Authorisation**: The sentinel proves the right agent is signalling.
   A file in `/tmp/jobs/` can be written by any process. An agent that
   completed successfully but whose sentinel was for a different task
   could corrupt state.

2. **Atomicity**: The signal command is a single operation. Writing a
   file, then having the supervisor detect it, is two operations with
   a window for partial state.

3. **Process termination**: Springfield intercepts the signal and
   terminates the agent process immediately. File-based approach
   requires the agent to exit cleanly after writing, which it may not.

**Adopt Springfield's sentinel protocol for agent-mux too.**

The job ticket file concept survives as the *payload* structure
(pr_number, worktree, baseline_failures), but the trigger is the
sentinel signal, not file presence.

Updated agent-mux exit contract:

```
Ralph's last steps:
  1. Run guardrails-{id}.sh — non-zero exit stops here
  2. gh pr create → capture PR number
  3. springfield signal \
       --sentinel {SENTINEL} \
       --status success \
       --reason "PR #{PR_NUMBER} opened" \
       --epic {TASK_ID}

  (The orchestrator reads the sentinel, gets PR number from reason
   field or from td, launches Bart, terminates Ralph.)

Failure path:
  springfield signal \
    --sentinel {SENTINEL} \
    --status failed \
    --reason "tests_failed: <last 3 lines of output>"
```

---

## 11. Open Questions (for next session)

**Q1: Lisa managing multiple option sets for parallel epics**

ADR-007 notes this is unresolved. If the orchestrator spawns multiple
Ralph worktrees in parallel (Manager-Worker), Lisa needs to maintain
separate ToT option sets per epic in flight. How does she track which
options belong to which epic? Proposal: each epic's approach layer in
`TODO-{id}.md` is the record. Lisa reads it back if she needs to
re-enter option evaluation for a blocked epic. `td` stores the
epic ID, which is the key into the right TODO file.

**Q2: Self-Consistency threshold**

How many passes before Lisa commits to an approach? ADR-007 says
"stable across passes" but doesn't define the count or what "stable"
means precisely. Proposal: 2 passes. If the conclusion differs, run
a third. If still differs, surface as a planning uncertainty and
consult Marge before proceeding. Document in the Epic Decomposition
Protocol.

**Q3: Marge's per-epic role**

We've established Marge fires per-release. But Lisa may encounter
scope ambiguity mid-epic (Ralph signals blocked, Lisa re-enters ToT,
but the re-planning reveals the Feature Brief was underspecified).
Does Marge get invoked to clarify? Proposal: yes, but as an exception
path, not a standard path. The orchestrator supports this: blocked
epic with `reason=feature_brief_ambiguous` → invoke Marge.

**Q4: Lovejoy's CHANGELOG source**

At per-release trigger, Lovejoy needs to produce release notes for
all epics in the milestone. Where does he read from?
Proposal: `td query "type = epic AND milestone = v0.7.0"` gives him
all epic titles and descriptions. `git log main --since=<last-release>`
gives him the commits. The Feature Briefs in PLAN.md give him the
user-facing intent. GECR: draft from these sources → critique against
the actual code merged → refine into coherent narrative.

**Q5: Guardrail script generation — LLM or deterministic?**

Lisa generates `guardrails-{id}.sh` during LRM. Is this:
(a) LLM-generated (Lisa writes it in her session) — flexible but
    potentially hallucinated; or
(b) Template-generated from scope.toml (deterministic) — reliable but
    less nuanced?
Proposal: template-generated for the structural checks (scope, test
count, commit format), LLM-generated for the semantic checks
(e.g., "verify the cloud block parser handles the `tags` variant").
Separate the two in the guardrail script with clear comments.

**Q6: Epic Decomposition Protocol**

ADR-007 defers the formal spec to `docs/standards/epic-decomposition-
protocol.md`. This document does not exist yet. It needs to specify:
- The precise three-layer TODO-{id}.md template
- The Self-Consistency sampling count and stability threshold
- Bart's retrospective signal schema (see Section 9)
- Lisa's option evaluation rubric (ADR constraints + Farley + Adzic)
- The formal definition of "epic closure" (what Bart must signal)

This is the most important missing document. Priority 1 and 2
implementations should be blocked until it exists (or written as
part of those epics).

---

## 12. Migration Path

Additive, not a rewrite. Each phase is independently deployable.

**Phase 1: Activate the Bart→Lisa feedback loop**
- Implement Priority 1 (Bart's retrospective signal)
- Write `docs/standards/epic-decomposition-protocol.md` (unblocks all else)
- Update `docs/adr/ADR-007` to `Accepted`
- No orchestrator changes required

**Phase 2: Activate Lisa's LRM model**
- Implement Priority 2 (Lisa's ToT + Self-Consistency)
- Implement Priority 3 (tighter handoff contract)
- Lisa now produces `scope.toml` + guardrails; Ralph consumes them
- Ralph's context shrinks; execution becomes more mechanical

**Phase 3: Correct trigger granularities**
- Implement Priority 6 (Lovejoy per-release trigger)
- Remove Marge's per-epic merge gate from orchestrator
- Add `StatusAllEpicsDone` check to orchestrator
- Per-epic merges become orchestrator-deterministic

**Phase 4: Ralph's OHECI escalation**
- Implement Priority 5
- Reduces context burn on unexpected failures
- Surfaces assumption breaks faster

**Phase 5: Marge's Discovery loops**
- Implement Priority 4
- Improves Feature Brief quality
- Reduces downstream rework from underspecified briefs

**Phase 6: agent-mux supervisor (if used independently of Springfield)**
- Implement Priority 7
- Decouple from tmux pane
- Concurrent Ralph + Bart slots

---

## 13. Summary

**What we're not changing:**
- Five agents, two loop levels — the distribution is correct at the right
  granularity. No agents eliminated.
- Sentinel protocol — adopt everywhere.
- `td` as state store — correct, keep.
- Go orchestrator — correct, keep.
- Ralph Wiggum stateless loop — correct, keep.
- Worktree provisioning before Ralph starts — correct, keep.

**What we are changing:**
- Marge and Lovejoy trigger granularity — per-release, not per-epic
- Lisa's prompt — ToT + Self-Consistency before LRM commit
- Lisa's handoff — machine-readable contract (scope.toml + guardrails),
  not just narrative prose
- Bart's output — dual: FEEDBACK.md (for Ralph) + retrospective signal
  (for Lisa via `td log --decision`)
- Ralph's prompt — OHECI escalation path for assumption breaks
- Per-epic merge — deterministic orchestrator, not Lovejoy

**The single highest-leverage change:**
Bart's structured retrospective signal (Priority 1). It activates
ADR-007's entire LRM model. Without it, every Lisa session starts cold
and the system cannot learn across iterations. Everything else is
improvement; this is the missing connection that makes the architecture
actually work as designed.

**Start here next session:**
1. Read `docs/adr/ADR-007` — it specifies exactly what's needed
2. Write `docs/standards/epic-decomposition-protocol.md` — formalises
   the schema that Priorities 1, 2, and 3 all depend on
3. Implement Priority 1 (Bart's retrospective signal) as a standalone epic
4. Implement Priority 2 (Lisa's ToT) as a follow-on epic that depends on 1
