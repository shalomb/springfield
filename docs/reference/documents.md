# Essential Documents Reference

Springfield Protocol uses 7 core documents as the shared state across all agents and skills.

## The 7 Core Documents

### 1. PLAN.md — Epic Backlog (Value Roadmap)

**Purpose:** Defines the *Epics* (value increments) and their justification.

**Structure (The "Three Pillars" Model):**
```markdown
# PLAN.md - Epic Backlog

## EPIC-XXX: [Epic Name]
**Value Statement:** For [User], who [Problem], the [Epic] is a [Solution] that [Benefit].

**The "Why":** [Brief summary of Discovery findings].
**Scope:**
- ✅ [Included]
- ❌ [Excluded]

**Acceptance Criteria:**
- [ ] [Observable outcome]

**Attributes:**
- **Status:** Ready | In Progress | Blocked | Done
- **Complexity:** Low/Med/High
- **Urgency:** Low/Med/High
- **Dependencies:** [Other Epics]
```

**Holds:**
- **Desirability:** Why users want it (Value Statement).
- **Viability:** Why we should build it (The "Why").
- **Feasibility:** What exactly is being built (Scope & Criteria).
- **Prioritization:** Status, complexity, urgency.

**Updated by:** Planning Agent (Flow Planner).
**Read by:** Implementation Agent (Ralph) to create `TODO.md`.

---

### 2. TODO.md — Sprint Tasks (Execution Plan)

**Purpose:** Decomposed, executable tasks for a *single* active Epic.

**Structure:**
```markdown
# TODO.md - Sprint for EPIC-XXX

## Context
[Link to PLAN.md#EPIC-XXX]

## Tasks
- [ ] Task 1: [Atomic implementation step]
  - Assigned to: [agent]
  - Subtasks: ...
```

**Holds:**
- Immediate, granular work items (10m - 2h).
- Implementation details.
- Progress tracking.

**Lifecycle:**
- Created by **Planning Agent** when an Epic starts.
- Updated by **Implementation Agent** as work progresses.
- **Deleted** by Implementation Agent when Epic is complete.

---

### 3. Feature.md — Feature Brief

**Purpose:** Detailed problem definition and requirements.

**Structure:**
```markdown
# Feature.md - [Feature Name]

## Problem
[Root cause analysis]

## Requirements
[User need statement]

## Constraints & Unknowns
...
```

**Holds:**
- The source of truth for "What" and "Why".
- Detailed requirements that feed into `PLAN.md` Epics.

---

### 4. ADRs — Architecture Decision Records

**Purpose:** Technical decisions and rationale.

**Holds:**
- Decisions on "How" to solve technical constraints.
- Alternatives considered.

---

### 5. scenarios.feature — BDD Specs

**Purpose:** Executable tests defining acceptance criteria.

**Holds:**
- Gherkin syntax scenarios.
- Used to verify Epics are complete.

---

### 6. FEEDBACK.md — Review Results

**Purpose:** Code review and quality gate status.

**Holds:**
- Issues found during implementation.
- Approval status.

---

### 7. CHANGELOG.md — Release History

**Purpose:** Record of shipped value.

**Holds:**
- Version history.
- Features added/fixed.
