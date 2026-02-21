# ADR-007: Epic Refinement, Lisa's LRM Role, and the Planning Loop

**Date:** 2026-02-19
**Status:** Proposed

---

## Context

The current Springfield planning loop has an unowned gap. The flow is:

```
Marge (Feature Brief) → Lisa (PLAN.md + TODO.md) → Ralph (execute) → Bart (verify)
```

Three problems have been identified:

**1. The option-selection gap.** When Lisa converts an Epic into a TODO.md task
list, she is implicitly making a solution-space decision — choosing one
implementation approach over alternatives — with no formal process, no recorded
rationale, and no feedback from previous iterations informing the choice.
Functional requirement decisions (the "how" of behaviour, not the "why" of
architecture) have no home. They appear silently inside `.feature` files without
the alternatives being recorded.

**2. The premature decomposition problem.** Lisa currently produces a pre-sliced
TODO.md task list before Ralph begins work. This is top-down decomposition done
at the worst possible moment — before execution feedback exists — and it removes
autonomy from Ralph, who is best positioned to decompose his own work via TDD.

**3. The flat fidelity problem.** All Epics in PLAN.md are treated as equally
refined regardless of their proximity to execution. Far-term Epics receive the
same structural treatment as the next Epic to be implemented, which is either
wasted effort (over-refining things that will change) or false confidence
(treating low-information placeholders as actionable plans).

---

## Decision

We adopt a **Just-In-Time Epic Refinement** model with the following principles:

### 1. Fidelity is intentionally graduated

PLAN.md holds a **deliberately low-fidelity, low-granularity** view of
far-term Epics. Premature refinement of distant Epics is waste — the information
needed to refine them correctly does not yet exist.

The fidelity gradient is:

| Stage | Artifact | Fidelity | Owner |
|:---|:---|:---|:---|
| Far-term backlog | PLAN.md Epic (stub) | Low — intent + rough scope only | Lisa |
| Near-term candidate | PLAN.md Epic (enriched) | Medium — options generated, not yet decided | Lisa |
| Ready for execution | TODO.md (high-fidelity Epic) | High — option decided, constraints explicit, guardrails stated | Lisa → Ralph |
| In execution | TODO.md (working) | Full — Ralph's decomposition added | Ralph |

### 2. Lisa's role shifts from task-decomposer to option-generator and LRM decision-maker

Lisa does **not** pre-slice Epics into tasks. That is Ralph's responsibility,
executed bottom-up via TDD as he works.

Lisa's planning role becomes:

- **Option generation:** Use Tree of Thoughts (ToT) to generate 2–3 candidate
  implementation approaches for the near-term Epic. Each candidate is evaluated
  against: ADR constraints, Farley Atomic predictively, Adzic Focused
  predictively, tech debt register signals, and retrospective learnings from
  the previous iteration.

- **LRM decision:** Lisa commits to one option at the **Last Responsible
  Moment** — after Bart has closed the previous Epic and produced a
  retrospective signal. She does not decide earlier, because every completed
  iteration produces constraint information that should inform the choice.

- **Self-Consistency validation:** Before committing, Lisa runs her constraint
  evaluation multiple times on the chosen option. If the evaluation is not
  stable across passes, the option has hidden ambiguity and should not be
  selected.

- **Iteration control:** Lisa decides what Bart's FEEDBACK.md means for the
  current Epic — fix now (Ralph iterates), defer (moves to debt register with
  a linked future Epic), or learning (enriches next option evaluation). She is
  the single consumer of Bart's retrospective output.

### 3. Ralph receives a high-fidelity Epic, not a task list

The TODO.md handed to Ralph is the **Epic restated at full fidelity** — context-
dense, constraint-explicit, guardrail-linked. It contains:

- **Intent layer** (immutable — from Marge's Feature Brief):
  The user need, the acceptance criteria, the definition of done.

- **Approach layer** (decided by Lisa at LRM — fixed for this iteration):
  The chosen option, the rationale, the alternatives that were considered and
  rejected.

- **Constraint layer** (inherited — not negotiable):
  Relevant ADR links synthesised for this Epic, Farley/Adzic thresholds,
  ACP requirements, tech debt items that must not be worsened.

- **Working layer** (Ralph's — fully owned by him):
  Ralph's own task decomposition, derived bottom-up via TDD. His failing tests
  *are* his task list. This layer is explicitly his to reshape.

Ralph can reshape the working layer freely. He cannot renegotiate the intent
or approach layers — those decisions are closed. If he discovers the approach
layer is wrong, he escalates via a distinct signal (see below).

### 4. Bart produces a dual output

Bart's existing FEEDBACK.md covers code-level review. Under this model Bart
produces a second output — a **retrospective signal** explicitly structured
for Lisa's consumption:

- What did this iteration reveal about the next Epic?
- What constraints proved wrong or incomplete?
- What did the Farley/Adzic scores reveal about the Epic's decomposability?
- Were any assumptions in the current Epic's approach layer invalidated?

This signal is the primary input to Lisa's next LRM decision. Without it,
Lisa is deciding with last-iteration's information rather than this-iteration's.

### 5. Option viability failure has an explicit escalation path

Not all failures are Ralph-level implementation problems. Some indicate that
Lisa's chosen option is fundamentally unviable — a constraint nobody foresaw,
an assumption that broke under execution.

Bart is responsible for distinguishing:
- **Implementation failure:** Ralph's code is wrong. Ralph iterates.
- **Option viability failure:** Lisa's chosen approach is wrong. Escalates to Lisa.

When Lisa receives an option viability failure:
1. She re-enters option evaluation with the new constraint information.
2. She consults her preserved ToT candidate list — if another candidate
   survives the new constraint, she selects it.
3. If no candidate survives, she re-runs ToT with the new constraint included.
4. Ralph pauses on approach-dependent work; continues on approach-independent
   work (tests, interfaces, scaffolding) while Lisa re-decides.
5. The failed option and the new constraint are recorded in the Epic's approach
   layer as a learning — not discarded.

---

## Consequences

### Positive

- **Lisa is fully utilised.** She moves from sophisticated ticket-writer to
  genuine architect of the delivery phase — holding cross-Epic context,
  generating options, deciding at the right moment, absorbing learnings.

- **Ralph is fully utilised.** He receives rich context and owns his own
  decomposition. His TDD discipline is the task-generation mechanism, not
  Lisa's upfront planning.

- **Bart's output gains architectural value.** His retrospective signal to
  Lisa makes every iteration an input to the next planning decision, not just
  a quality gate on the current one.

- **The fidelity gradient prevents waste.** Far-term Epics stay deliberately
  coarse. Refinement effort is concentrated where it has value — the next
  Epic to execute, informed by the most recent retrospective.

- **Functional decisions get recorded.** The approach layer of the high-fidelity
  Epic captures what was considered and why the chosen option won. This fills
  the gap between ADRs (architectural/NFR decisions) and `.feature` files
  (behavioural outcomes) — the reasoning behind functional choices now has a home.

- **No new agent required.** The restructuring is achieved by redistributing
  responsibilities within the existing five-agent model. Adding an agent
  upstream of Lisa would add latency at both ends of the slowest part of the
  cycle and hollow out Lisa's role.

### Negative

- **Lisa's planning pass takes longer.** ToT exploration and Self-Consistency
  validation before LRM commitment is more expensive than writing a task list.
  This is intentional — the cost is front-loaded into planning rather than
  back-loaded into rework.

- **The LRM wait introduces a sequencing constraint.** Lisa cannot refine
  the next Epic until Bart closes the current one. This implies sequential
  Epic execution by default. Parallel workstreams (Manager-Worker loop) need
  explicit handling — Lisa may need to maintain multiple option sets in flight.

- **Ralph's autonomy requires trust.** If Ralph self-decomposes poorly, there
  is no pre-written task list to fall back on. The Farley per-test checklist
  and ACP discipline are the guardrails — they must be internalised, not
  enforced by Lisa's upfront planning.

- **EPIC-005 (Agent Governance) is a prerequisite for ToT and Self-Consistency
  at scale.** Multiple LLM passes during Lisa's option evaluation are expensive.
  Without a budget enforcer, Lisa could be the most expensive agent in the
  system.

---

## Amendment A: ADR Lifecycle (2026-02-19)

The original decision implied Bart reviews ADRs before Ralph executes. This is
revised.

**ADRs drafted by Lisa at LRM are hypotheses, not ratified decisions.**

The revised lifecycle is:

| Stage | Actor | ADR Action |
|:---|:---|:---|
| LRM refinement | Lisa | Drafts ADR. Status: `Proposed`. Records the hypothesis and the alternatives considered. |
| Execution | Ralph | Implements against the ADR. Surfaces assumption breaks in his working layer — not judgement, just evidence. |
| Review | Bart | Reviews Ralph's code *and* the ADR together. One invocation, dual output. |
| Bart verdict: Confirmed | Bart → Lisa | ADR status moves to `Accepted`. Becomes a constraint for future Epics. |
| Bart verdict: Partially confirmed | Bart → Lisa | ADR updated with caveats before acceptance. |
| Bart verdict: Invalidated | Bart → Lisa | ADR status moves to `Rejected`. Option viability failure signal raised. Lisa re-enters option evaluation. The failed option and new constraint are preserved as a learning — not discarded. |

**Why this ordering:**

Bart reviewing the ADR before Ralph executes means Bart is critiquing Lisa's
reasoning in the abstract — he has no data. Bart reviewing after Ralph executes
means he has empirical evidence. His ADR verdict is zero additional latency —
it happens in the same pass as his code review.

Ralph's existing responsibility — *"Flag Surprises: when assumptions don't hold
in practice"* — extends explicitly to ADR assumptions. He surfaces tensions; he
does not protect the ADR at the expense of honest implementation.

---

## Amendment B: State Machine and Structured Format (2026-02-19)

The planning loop described in this ADR requires a **mechanically enforceable
state machine**. The current harness (Justfile) infers state from file presence
(`TODO.md` exists?) and keyword grep on prose (`FEEDBACK.md` contains
"critical"?). This is fragile for the same reasons EPIC-007 identified with
string matching — it is pattern matching on prose, not typed state.

**Two decisions are deferred to ADR-008:**

**1. Epic state as a structured, typed attribute.**

PLAN.md markdown is human-readable but machine-hostile. The orchestrator needs
to read Epic status as a typed value from a structured format, not infer it from
emoji and heading conventions. The valid state enumeration (from the git
branching standard) is:

```
planned → ready → in_progress → implemented → verified → done
                                                        → deferred
                                                        → blocked
```

Whether this is YAML, TOML, or a hybrid (PLAN.md for human narrative +
PLAN.yaml for machine state) is the core decision for ADR-008.

**2. The Springfield binary as the orchestrator.**

The Justfile cannot be tested. The orchestration logic — which agent to invoke,
under what Epic state, with what signal — cannot be unit tested as shell
conditionals. The Springfield binary (Go) is the correct home for this logic
because:

- State transitions become pure functions, coverable by table-driven tests
- Epic states become typed enums; invalid transitions are compile errors
- The Justfile becomes a thin invocation layer (`just do` calls the binary;
  the binary decides what happens next)
- The full delivery loop is backed by the same test discipline Ralph applies
  to production code

The signal type taxonomy Bart produces (implementation failure vs. option
viability failure vs. ADR verdict) must also be a typed enumeration readable
by the binary — not keyword-matched from free prose in FEEDBACK.md.

**A third open question flagged for evaluation:**

`td(1)` (https://github.com/marcus/td) is a structured CLI todo list manager.
It may offer a foundation for Ralph's working layer task management within the
high-fidelity Epic — replacing the freeform TODO.md working layer with a
machine-readable, queryable task store. To be evaluated before ADR-008 is
finalised.

---

## What This Does Not Decide

The following are explicitly deferred to a subsequent standard:

- The precise structure of the high-fidelity Epic / TODO.md template
  (the three-layer format described above is directional, not final).
- The exact format of Bart's retrospective signal to Lisa.
- How Lisa manages multiple option sets during parallel Epic execution.
- The Self-Consistency sampling count and stability threshold.
- The formal trigger definition for "Bart has closed an Epic"
  (what exactly constitutes closure beyond FEEDBACK.md being clear).

These are implementation details of the **Epic Decomposition Protocol**,
which this ADR motivates and which will be specified in
`docs/standards/epic-decomposition-protocol.md`.

---

## References

- `docs/concepts/model.md` — The Discovery/Delivery Diamond model
- `docs/reference/loop-catalog.md` — Tree of Thoughts, Self-Consistency,
  Plan-and-Execute, Ralph Wiggum Variant
- `docs/reference/farley-index.md` — TDD test quality rubric (used predictively
  by Lisa during option evaluation)
- `docs/reference/adzic-index.md` — BDD specification quality rubric (used
  predictively by Lisa during option evaluation)
- `docs/standards/atomic-commit-protocol.md` — Commit atomicity standard
  (constraint on Ralph's working layer)
- `.pi/agents/lisa.md` — Lisa's agent definition (to be updated)
- `.pi/agents/ralph.md` — Ralph's agent definition (to be updated)
- `.pi/agents/bart.md` — Bart's agent definition (to be updated)
- Adzic, G. — *Specification by Example*
- Farley, D. — *Modern Software Engineering*
- Last Responsible Moment:
  https://levelup.gitconnected.com/the-last-responsible-moment-an-architects-perspective-4d94d2276066
- `td(1)` — structured CLI todo list manager (under evaluation):
  https://github.com/marcus/td
- ADR-008 — Planning State Boundary, td(1) Adoption, and Springfield Binary
  Orchestration (implements the mechanical substrate this ADR requires)
