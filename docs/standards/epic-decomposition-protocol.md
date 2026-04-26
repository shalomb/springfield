# Epic Decomposition Protocol (EDP)

**Status:** Active  
**Version:** 1.1  
**Related ADR:** [ADR-007-epic-refinement-and-lisa-lrm.md](../adr/ADR-007-epic-refinement-and-lisa-lrm.md)

## 1. Overview
The Epic Decomposition Protocol governs the structural handoff between agents in the Inner Loop (Lisa → Ralph → Bart). It establishes the "Constraint-Driven" model: Lisa provides the boundaries and constraints, Ralph acts as the Archi-Engineer who autonomously designs and executes the solution within those boundaries, and Bart verifies both the code and the architectural hypothesis.

## 2. The High-Fidelity Epic: `TODO-{id}.md`
Lisa generates this artifact during her planning phase. It consists of three immutable layers and one mutable working layer for Ralph.

### Template Structure

```markdown
# TODO-td-{id}.md — {Epic Title}

## Intent Layer (Immutable)
**Origin:** Marge's Feature Brief
**User Need:** {Why are we doing this?}
**Acceptance Criteria (BDD):**
- {Gherkin scenarios to be satisfied}

## Context & Constraint Layer (Lisa's Guidance)
**Architectural Hypothesis:** {Proposed architectural approach, if applicable}
**Rationale:** {Why this approach was considered best}
**ADRs in Effect:**
- [ADR-XXX](...) - {Specific constraint to uphold}
**Tech Debt Landmines:**
- {Areas of the codebase to avoid worsening}
**Quality Boundaries:**
- Farley Index / Adzic Index targets

## Working Layer (Ralph's Execution & Architectural Reality)
*Ralph fills this layer bottom-up. Failing tests become the task list. If the proposed architectural hypothesis breaks, Ralph documents the pivot here.*
```

## 3. The Constraint-Driven Execution Model

The redesign intentionally removes brittle, LLM-generated bash scripts (`guardrails.sh`) and premature file-scoping (`scope.toml`).

### 3.1. Lisa (Context Provider)
Lisa's job is not to dictate step-by-step tasks or precise file edits. She provides the **Intent** and the **Constraints**. She drafts `Proposed` ADRs as hypotheses and identifies which existing boundaries Ralph must respect.

### 3.2. Ralph (Archi-Engineer)
Ralph is not a naive executor; he is a Staff-level Archi-Engineer. He receives the `TODO-{id}.md` and uses TDD to explore the codebase.
- He autonomously decides which files to modify or create.
- If Lisa's proposed architecture fails upon contact with the code (an Assumption Break), Ralph is empowered to pivot, design a better approach within the constraints, and document the change in his Working Layer.
- If a constraint makes the feature impossible to build, Ralph halts execution and signals `blocked` with the OHECI pattern.

### 3.3. Baseline Constraints (Orchestrator Enforced)
Instead of LLM-generated guardrail scripts, the Orchestrator deterministically enforces structural constraints when Ralph calls `springfield signal --status success`:
1. `just test` must exit 0 (BDD and TDD pass).
2. Commit messages must follow the Atomic Commit Protocol (ACP).

## 4. Lisa's Planning Loop

Before generating the `TODO-{id}.md`, Lisa executes a Tree of Thoughts (ToT) + Self-Consistency loop.

### 4.1. The Option Evaluation Rubric
For every epic, Lisa must generate at least 2 distinct implementation approaches. She evaluates each against:
1. **ADR Conformance:** Does it violate any established architectural decisions?
2. **Farley/Adzic Indices (Predictive):** Will this yield atomic, fast tests and declarative step definitions?
3. **Retrospective Signals:** Does it mitigate the failures logged in Bart's previous retrospectives?

### 4.2. Self-Consistency Threshold
Lisa runs the evaluation **twice** independently.
- **Stable (2 identical conclusions):** Lisa commits the approach as her hypothesis.
- **Unstable (Differing conclusions):** Lisa runs a 3rd tie-breaker pass. If still fundamentally unstable, Lisa signals `blocked` with `reason=feature_brief_ambiguous` and orchestrator invokes Marge.

## 5. Bart's Dual Output & Retrospective Signal

Bart acts as the **Systemic Quality Gate**. He verifies the code (tactical), ensures rules were followed (structural), and confirms if Lisa's initial hypotheses were correct (architectural).

He produces two artifacts:
1. `FEEDBACK.md`: Tactical code review for Ralph (if the PR is rejected).
2. **Retrospective Signal**: A structured payload sent to `td log --decision <epic-id>` for Lisa.

### 5.1. Definition of Epic Closure
An epic is only closed when Bart issues `springfield signal --status success`.

### 5.2. Retrospective Signal Schema
Bart must log this payload to the `td` database before signalling success or failure.

```markdown
## Retrospective Signal

### Approach Layer Verdict
verdict: confirmed | partially_confirmed | invalidated
notes: >
  {Did Lisa's architectural hypothesis survive contact with Ralph's code?}

### ADR Verdict  
adr_id: {ADR-XXX}
verdict: confirmed | rejected | amended
notes: >
  {If a Proposed ADR was tested, what is the empirical verdict?}

### Constraints Revealed
- {Constraint 1 that Lisa should know for the next epic}

### Tech Debt Signals
- severity: {minor|medium|critical}
  item: {Description of debt introduced or discovered}
```

If the `Approach Layer Verdict` is `invalidated`, Bart signals `blocked` (Option Viability Failure), skipping Ralph's retry loop and returning the Epic directly to Lisa.