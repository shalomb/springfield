# Task Decomposition Standard

**Owner:** Ralph (Build Agent)
**Applies to:** All task creation within the working layer of a high-fidelity
Epic handoff (see ADR-007, ADR-008).
**Referenced by:** `.github/agents/ralph.md`, `.pi/agent/skills/ralph/SKILL.md`

---

## Purpose

This standard defines what a well-decomposed Epic task looks like and provides
Ralph with concrete strategies for deriving his working layer from Lisa's
high-fidelity handoff document (`TODO-{td-id}.md`).

It is not a protocol — Ralph owns decomposition and exercises judgement. This
standard gives him a shared vocabulary, a bias toward good decomposition, and
a continuity convention so that subsequent Ralph sessions inherit the
decomposition intent of prior sessions without re-doing the work.

External parties reading this document will understand how Springfield
translates Epic acceptance criteria into executable, testable units of work.

---

## The INVEST Properties

Every `td` task Ralph creates must satisfy INVEST. These properties are
evaluated at the moment of task creation — before the first line of test code
is written.

### Independent
The task can be started and completed without depending on another incomplete
task. If two tasks are entangled — if you can't write the failing test for one
without the other being done first — there is a hidden shared dependency. That
dependency is its own task. Extract it and make it explicit in td:

```bash
td dep add td-a3f8.2 td-a3f8.1   # task 2 depends on task 1
```

### Negotiable
Within the working layer, Ralph can reshape tasks as he learns. The intent and
approach layers from Lisa (`TODO-{td-id}.md`) are fixed. The working layer task
list is Ralph's to adjust as the implementation reveals new information. INVEST
Negotiable is the named permission for this — it is not scope creep, it is
responsive decomposition.

When Ralph reshapes a task, he logs the reasoning:

```bash
td log td-a3f8.2 --decision "split into read/write paths — different failure modes"
```

### Valuable
Every task must be traceable to an acceptance criterion in the intent layer.
A task that serves Ralph's implementation convenience but maps to no
user-observable behaviour is waste. If Ralph cannot point to the acceptance
criterion a task satisfies, the task should not exist.

This is the Farley **Necessary** property applied at the task level rather than
the test level.

### Estimable
If Ralph cannot articulate what the failing test for a task looks like, the
task is too vague to start. The test definition *is* the estimate. You do not
estimate in story points — you estimate by being able to describe the Red state
in one sentence:

> *"A test that calls `BudgetEnforcer.Check()` with a session that has exceeded
> its token limit and asserts it returns `ErrBudgetExceeded`."*

If that sentence cannot be written, split or clarify the task first.

### Small
One Red-Green-Refactor cycle. One ACP commit. If a task would require more
than one cycle to complete, it needs splitting.

The measurable proxy: if Ralph anticipates a Farley **Atomic** setup-to-assertion
ratio above 3:1 for the test this task demands, the task is too broad. Either
the task needs splitting, or the production code design needs rethinking before
implementation begins.

### Testable
Every task has exactly one failing test that defines it. The test is written
before any implementation. If a task cannot be expressed as a single failing
test, it is either too broad (split it) or not a unit of work at all
(it may be a constraint or a design decision — record it in td as a
`--decision` log, not a task).

---

## Decomposition Strategies

When an acceptance criterion is too large for one Red-Green-Refactor cycle,
Ralph uses one of these named strategies to find the cut. The strategy chosen
is logged as a `--decision` in td so subsequent sessions know why the tasks
are shaped the way they are.

### By Workflow Step
Decompose along the sequence of a user journey. Each step in the flow is a
task. The test for each step asserts the behaviour at that boundary, not the
end-to-end outcome.

**When to use:** Features with a sequential flow (request → validate → process
→ respond). Each step has distinct inputs, outputs, and failure modes.

**Example:** Budget enforcement has three steps — receive request, check
remaining budget, reject or permit. Three tasks.

### By Business Rule
Each distinct rule is a task. Rules that share a code path are still separate
tasks if they have different test conditions.

**When to use:** Acceptance criteria that list multiple "when X, then Y"
conditions. Each condition is a task.

**Example:** "Budget enforcement cuts off at the session limit" and "Budget
enforcement cuts off at the daily limit" are two rules, two tasks — even if
they share an implementation function.

### By Data Variation (Happy Path First)
Start with the simplest valid input. Make it work. Then add boundary cases
and edge cases as separate subsequent tasks. Never generalise before you have
one concrete case working.

**When to use:** Any feature involving input validation, data processing, or
calculation. The happy path is always the first task.

**Ordering:**
1. Happy path (valid, typical input)
2. Boundary values (minimum, maximum, exactly at limit)
3. Invalid input (missing fields, wrong types)
4. Edge cases (empty, nil, zero, overflow)

### By Population / Persona
When behaviour differs by user type, each population is a task. Do not test
multiple personas in one scenario — that violates Adzic **Focused**.

**When to use:** Features with role-based behaviour (admin vs. user vs. guest,
authenticated vs. anonymous).

**Example:** Budget enforcement for a session owner vs. budget enforcement for
a delegated agent are two tasks — same code path, different authorisation
context.

### By Read vs. Write
Separate query behaviour from mutation behaviour. They have different failure
modes, different concurrency concerns, and different test surfaces.

**When to use:** Any data-centric feature. Reading the budget state is a task.
Updating the budget state is a task. Enforcing the budget state is a task.

### By Happy Path vs. Error Path
Never mix success and failure behaviour in one task. The happy path is one
task. Each named error condition is its own task.

**When to use:** Always. This is the default decomposition cut when no other
strategy is more specific.

**Example:**
- Task 1: Budget check passes — request proceeds
- Task 2: Budget check fails — request rejected with `ErrBudgetExceeded`
- Task 3: Budget check fails — graceful degradation to fallback model

---

## Sequencing: Where to Start

Ralph does not start with the acceptance criterion that is easiest to
implement. He starts with the one whose failing test forces the most
foundational design decision — the one that, once answered, shapes everything
else.

**The sequencing heuristic:**

1. **Start with the boundary that defines the core abstraction.** If the Epic
   requires a new interface or type, the first task is the failing test that
   demands it. Everything else is built on top.

2. **Happy path before error paths.** Within a strategy, always establish the
   working case before handling failure.

3. **Independent before dependent.** Tasks with no blockers in the td
   dependency graph come before tasks that depend on them. `td ready` gives
   Ralph the unblocked queue.

4. **Acceptance criteria order is not implementation order.** Lisa writes
   acceptance criteria in user-facing terms. Ralph sequences tasks in
   implementation terms. These are not the same sequence.

---

## The Continuity Convention

Ralph's decomposition must survive context resets. A new Ralph session must be
able to pick up the decomposition exactly where the previous session left it —
not re-derive it from scratch, not re-shape it differently.

**The convention:** When Ralph ends a session, he records the decomposition
reasoning in td, not just the task status.

```bash
td handoff td-a3f8.1 \
  --done "happy path budget enforcement — BudgetEnforcer.Check() passes" \
  --remaining "error path: ErrBudgetExceeded, daily limit variant" \
  --decision "decomposed by error path per task-decomposition.md — each \
              named error is a separate task" \
  --uncertain "grace period not specified in acceptance criteria — needs Lisa"
```

A new Ralph session runs `td usage` at cold-start and reads:
- What was done
- What shape the remaining tasks have
- *Why* they are shaped that way
- What is genuinely unclear and needs escalation

The `--decision` field is the critical one. Without it, the next Ralph sees
tasks but not the reasoning. With it, the decomposition is self-documenting
across sessions.

**The bias this creates:** Subsequent Ralph sessions are biased to *continue*
the decomposition strategy already chosen, not to start fresh. They may add
tasks, but they do not reshape the existing ones without logging a reason.
The commit graph remains coherent. Bart can read the git history and the td
handoff log together and reconstruct the full decomposition reasoning.

---

## td Usage Conventions for Ralph

```bash
# At session start — load full context
td usage                          # live task state, current focus, decisions

# Creating tasks from acceptance criteria
td create "budget check: happy path" --priority P1 --parent td-a3f8
td create "budget check: ErrBudgetExceeded" --priority P1 --parent td-a3f8
td dep add td-a3f8.2 td-a3f8.1   # error path depends on happy path

# Starting a task (atomic claim — safe under parallel Ralph instances)
td start td-a3f8.1

# Logging as you work
td log td-a3f8.1 "Red: BudgetEnforcer.Check() test written"
td log td-a3f8.1 "Green: minimal implementation passes"
td log td-a3f8.1 --decision "using middleware pattern — see ADR-008"

# Ending a session
td handoff td-a3f8.1 \
  --done "happy path implemented and committed" \
  --remaining "error paths: ErrBudgetExceeded, daily limit" \
  --decision "decomposed by error path — task-decomposition.md §By Happy Path vs Error Path" \
  --uncertain "grace period behaviour unclear"

# Checking what's unblocked next
td ready                          # tasks with no open blockers
td next                           # highest priority unblocked task
```

---

## Size Reference Card

| Signal | Meaning | Action |
|:---|:---|:---|
| Farley Atomic ratio > 3:1 anticipated | Task too broad | Split using a decomposition strategy |
| Cannot write the failing test in one sentence | Task too vague | Clarify or split |
| Two acceptance criteria in one task | Task not focused | One criterion per task |
| Task depends on incomplete task | Not independent | Extract the dependency as its own task |
| No acceptance criterion traceable | Not valuable | Delete or convert to `--decision` log |
| More than one ACP commit anticipated | Not small | Split |

---

## References

- `docs/reference/farley-index.md` — Test quality rubric; Atomic and Necessary
  properties apply at task level
- `docs/reference/adzic-index.md` — Specification quality; Focused property
  applies to task scope
- `docs/standards/atomic-commit-protocol.md` — One task = one ACP commit
- `docs/adr/ADR-007-epic-refinement-and-lisa-lrm.md` — Working layer ownership
- `docs/adr/ADR-008-planning-state-td-springfield-orchestrator.md` — td as task
  store; TODO-{td-id}.md handoff contract
- INVEST criteria: Bill Wake, *"INVEST in Good Stories, and SMART Tasks"* (2003)
