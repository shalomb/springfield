# Glossary

Key terms used in the Springfield Protocol.

## Core Concepts

**Ralph Wiggum Loop**
: A stateless resampling engine that ensures quality through persistent iteration. Each iteration starts from a clean context to prevent hallucination and context rot.

**Ephemeral Context**
: A clean, isolated environment for each iteration (e.g., git worktree). Prevents accumulated errors across iterations.

**Stateless Resampling**
: Iteration without memory of previous failures. Each attempt is independent; failures inform strategy but don't accumulate.

**PLAN.md**
: The source of truth. Tracks epics, tasks, and their validation state. Updated after each iteration.

**TODO.md**
: Executable task list derived from PLAN.md. What agents work from.

**Definition of Ready (DoR)**
: Checklist ensuring an issue is clear enough to start work. Enforced by Wiggum.

**Feature Brief**
: Validated problem statement. Output of discovery track; input to delivery track.

**Discovery Track**
: Design thinking phase. Understand user needs, validate assumptions, document unknowns.

**Delivery Track**
: Agile execution phase. Plan, implement, verify, release.

## Loops & Patterns

**Sense-Plan-Act**
: Observe → Think → Act. The foundational agent loop.

**ReAct**
: Reason + Act. Verbalize thinking before taking actions.

**Tree of Thoughts (ToT)**
: Explore multiple reasoning paths, evaluate, and prune low-scoring options.

**Plan-and-Execute**
: Break into tasks upfront, then execute sequentially.

**GECR**
: Generate → Evaluate → Critique → Refine. Iterative polishing loop.

**TALAR**
: Test → Analyze → Learn → Adjust → Retest. Experiment-driven optimization.

**Manager-Worker**
: Orchestrator (manager) coordinates specialized workers in parallel.

**Dialogue Loop**
: Two agents iterating (propose → critique → refine).

## Roles & Responsibilities

**Gate**
: A decision point where a character blocks or approves progress.

**Skill**
: A character's specialized capability (e.g., "TDD Execution").

**Ephemeral Agent**
: An agent spawned for a single task with clean context, then discarded.

## File Organization

**PLAN.md**
: Roadmap of epics and their state (planned, in-progress, verified, released).

**TODO.md**
: Task list for immediate execution. Derived from PLAN.md.

**ADR (Architecture Decision Record)**
: Document recording architectural decisions, rationale, and implications.

**KEDB (Known Error Database)**
: Crowdsourced troubleshooting knowledge base for common issues.

## Workflow Concepts

**Five Whys**
: Questioning technique to uncover root causes (used in discovery).

**Gemba Walk**
: Investigating actual systems/docs to gather context (used in discovery).

**Unknowns Map**
: Explicit capture of what we don't know and its impact.

**Zero-Change Import**
: Brownfield import of existing infrastructure with no changes to the imported state.

**Mock-First Testing**
: Write mocks before implementation to clarify dependencies.

**95%+ Coverage**
: Herb's quality standard. Code coverage threshold for release.

---

For more details, see the root documentation:
- [`../START_HERE.md`](../START_HERE.md)
- [`../QUICK_START.md`](../QUICK_START.md)
- [`../core-principles.md`](../core-principles.md)
