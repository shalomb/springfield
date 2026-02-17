# Essential Documents Reference

Springfield Protocol uses 7 core documents as the shared state across all agents and skills.

## The 7 Core Documents

### 1. PLAN.md — Roadmap & Status

**Purpose:** Epic-level roadmap showing what we decided to build and current status

**Structure:**
```markdown
# PLAN.md - Feature Roadmap

## Feature: [Feature Name]

### Epic 1: [Epic Name]
- [ ] Task 1: [Task description]
  - Status: unstarted | in-progress | verified | failed
  - Assignee: [agent name]
  - Depends on: [other tasks]
  - ADR: ADR-XXX (decision made)
  - Feature Brief: Feature.md#anchor
  - BDD: scenarios.feature#anchor
```

**Holds:**
- Current state of all epics and tasks
- Task status (drives Ralph Wiggum Loop)
- Links to architectural decisions (ADRs)
- Links to feature requirements

**Updated by:** Planning skill (create), all skills (status updates)

**Read by:** All agents (to know current state), Ralph Wiggum Loop (to find next task)

---

### 2. TODO.md — Sprint Tasks

**Purpose:** Executable task list for current sprint

**Structure:**
```markdown
# TODO.md - Current Sprint

- [ ] Task 1: [Task description]
  - Assigned to: [agent]
  - Subtasks:
    - [ ] Subtask A
    - [ ] Subtask B
  - Blocked by: [other tasks]
  - Time estimate: [hours]
  - Started: [date]
  - Learning: "[What we discovered]"
```

**Holds:**
- Immediate work queue
- Task progress and subtasks
- Dependencies and blockers
- Learning captured during work

**Updated by:** Implementation skill, testing skill (during work)

**Read by:** Implementation agents (to know what to work on)

---

### 3. Feature.md — Feature Brief

**Purpose:** Feature requirements, constraints, and decision context

**Structure:**
```markdown
# Feature.md - [Feature Name]

## Problem
[Root cause analysis - why we're building this]

## Requirements
[User need statement]

## Acceptance Criteria
See: scenarios.feature#anchor

## Constraints & Unknowns
- [Hard constraint]
- **Unknown:** [Question] - See ADR-XXX for decision
- **Assumption:** [What we're betting on]

## Scope
✅ [In scope]
❌ [Out of scope - future]

## Success Criteria
- [Measurable success metric]
```

**Holds:**
- Problem statement (prevents building wrong feature)
- Requirements and user needs
- Constraints (from ADRs)
- Explicit unknowns (linked to ADRs where resolved)
- Explicit assumptions (what could break)
- Scope boundaries

**Updated by:** Discovery skill (initial), learning skill (when assumptions break)

**Read by:** All agents (to validate requirements), implementation agents (to know what to build)

---

### 4. ADRs — Architecture Decision Records

**Purpose:** Record architectural decisions with full rationale

**Location:** `docs/adr/ADR-XXX-title.md`

**Structure:**
```markdown
# ADR-XXX: [Decision Title]

**Date:** YYYY-MM-DD
**Status:** Proposed | Accepted | Superseded
**Related Feature:** Feature.md#anchor

## Problem
[What question did we need to answer?]

## Decision
[What did we decide?]

## Rationale
[Why this over alternatives?]

## Consequences
- [Benefit/tradeoff]
- [Implementation note]

## Alternatives Considered
1. [Alternative A] - Why not?
2. [Alternative B] - Why not?
3. [Chosen option] ✅
```

**Holds:**
- Architectural decisions with full thinking
- Unknowns from Feature.md that were resolved here
- Alternatives considered (captures thinking)
- Tradeoffs and consequences

**Updated by:** Architecture skill

**Read by:** Implementation skill (to follow ADR), review skill (to validate ADR compliance)

---

### 5. scenarios.feature — BDD Specs

**Purpose:** Executable acceptance criteria

**Location:** `features/[domain]/[feature].feature`

**Structure:**
```gherkin
# features/[domain]/[feature].feature

Feature: [User-facing feature name]
  Problem: [What problem does this solve]
  Feature: Feature.md#anchor

  Scenario: [Specific case]
    Given [precondition]
    When [action]
    Then [expected outcome]
```

**Holds:**
- Acceptance criteria (executable)
- Test cases for all scenarios
- Constraints documented as scenarios (rate limiting, timeouts, etc.)

**Updated by:** Feature specification + implementation

**Read by:** Testing skill (to execute tests), implementation skill (to know what to code)

---

### 6. FEEDBACK.md — Review & Gate Results

**Purpose:** Code review feedback, questions, and gate results

**Structure:**
```markdown
# FEEDBACK.md - Review Results

## Task X: [Task description]
**Reviewed:** YYYY-MM-DD HH:MM UTC
**Reviewer:** [skill/agent]

### Issues Found
- [x] Fixed: [Issue]
- [ ] Open: [Issue]
  - Risk: [Level]
  - Suggested fix: [Recommendation]

### Questions
- Why [decision]?
  - Answer: [Rationale, links to ADR/Feature.md]

### Gate Results
- [x] Tests pass: XXX/XXX ✓
- [x] Coverage: XX% (> 95%) ✓
- [x] Security: [Result]
- [x] Performance: [Result]

**Status:** READY TO MERGE | NEEDS FIXES | BLOCKED
```

**Holds:**
- Code review issues and questions
- Gate check results (coverage, security, performance)
- Blocking issues and approvals
- Learning captured from review

**Updated by:** Review skill, verification skill

**Read by:** Implementation agents (to fix issues), teams (to understand gate status)

---

### 7. CHANGELOG.md — Release History

**Purpose:** Record what changed and learning from each release

**Structure:**
```markdown
# CHANGELOG.md

## [Unreleased]

### Added
- [Feature] (see Feature.md#anchor)
  - [Capability]
  - [Capability]

### Fixed
- [Bug fix] (ADR-XXX rationale)

### Learning
- [What we learned]
  - [Discovery during implementation]
  - [Updated Feature.md#constraints]

### Security
- [Security improvement]

---

## [1.0.0] - YYYY-MM-DD

### Added
- [Feature]

...
```

**Holds:**
- What changed in each release
- Features added with links to Feature.md
- Bug fixes and improvements
- Learning captured (what surprised us, what changed)
- Security improvements

**Updated by:** Release skill (structure), learning skill (learning notes)

**Read by:** Teams (release notes), stakeholders (what's in this version)

---

## Document Flow Through Phases

### Discovery Phase
```
Discovery Skill
  ↓ produces
Feature.md (problem + requirements)
  ↓ populates unknowns
Feature.md#unknowns
  ↓ triggers
Architecture Skill
  ↓ creates/updates
ADRs (decisions for unknowns)
  ↓ links back
Feature.md → ADR-XXX
```

### Delivery Phase
```
Planning Skill
  ↓ consumes Feature.md + ADRs
  ↓ produces
PLAN.md + TODO.md
  ↓ agents work from
Implementation/Testing Skill
  ↓ execute against
scenarios.feature
  ↓ update
TODO.md (progress + learning)
  ↓ review code against
Review Skill
  ↓ produces
FEEDBACK.md (issues + questions)
  ↓ checks
Verification Skill
  ↓ updates
PLAN.md (task status → verified)
```

### Release Phase
```
Release Skill
  ↓ consumes Feature.md + FEEDBACK.md + changes
  ↓ produces
CHANGELOG.md (what changed)
  ↓ captures learning
Learning Skill
  ↓ updates
CHANGELOG.md (what surprised us)
  ↓ backfeeds
Feature.md (if assumptions broken)
```

---

## What Each Document Holds

| Document | Purpose | Key Content |
|----------|---------|-------------|
| **PLAN.md** | Epic roadmap | Status, decisions (via links), current work |
| **TODO.md** | Sprint tasks | What agents work on, progress, learning |
| **Feature.md** | Requirements | Problem, constraints, assumptions, unknowns |
| **ADRs** | Decisions | Unknowns resolved → decisions + rationale |
| **scenarios.feature** | Acceptance criteria | Executable tests, constraints as scenarios |
| **FEEDBACK.md** | Review results | Issues, questions, gate status |
| **CHANGELOG.md** | Release history | What changed, learning, versions |

---

## Where Key Information Lives

### Problem Statement
- **Lives in:** Feature.md#problem
- **Captured by:** Discovery skill
- **Updated by:** Learning skill (if root cause discovered)

### Unknowns & Questions
- **Lives in:** Feature.md#unknowns
- **Linked to:** ADRs (where resolved)
- **Captured by:** Discovery + architecture skills
- **Resolved by:** Architecture skill → ADR

### Assumptions
- **Lives in:** Feature.md#assumptions
- **Updated by:** Learning skill (when broken during work)
- **Learning:** TODO.md → FEEDBACK.md → CHANGELOG.md

### Decisions Made
- **Lives in:** ADRs (full reasoning)
- **Referenced in:** Feature.md (links), PLAN.md (task links)
- **Validated in:** Review skill (FEEDBACK.md checks compliance)

### Learning Captured
- **During work:** TODO.md (learning notes)
- **During review:** FEEDBACK.md (questions, surprises)
- **During release:** CHANGELOG.md (what changed, what we learned)
- **Backfeed:** Updates Feature.md if assumptions invalidated

---

## Agent/Skill Awareness

### Every agent should have access to:
1. **PLAN.md** — Current roadmap and task status
2. **TODO.md** — Current sprint tasks
3. **Feature.md** — Requirements and constraints
4. **ADRs** — Architectural guardrails
5. **scenarios.feature** — Acceptance criteria
6. **FEEDBACK.md** — Review results and gate status

### Every skill should produce/update:
| Skill | Produces | Updates |
|-------|----------|---------|
| **Discovery** | Feature.md | - |
| **Architecture** | ADRs | Feature.md (ADR links) |
| **Planning** | PLAN.md, TODO.md | - |
| **Implementation** | Code | TODO.md (progress, learning) |
| **Testing** | Tests, coverage | TODO.md (test results, learning) |
| **Review** | FEEDBACK.md | - |
| **Verification** | Gate sign-off | PLAN.md (status), FEEDBACK.md |
| **Release** | CHANGELOG.md, version tag | - |
| **Learning** | - | Feature.md, TODO.md, FEEDBACK.md, CHANGELOG.md |

---

## Implementation Tips

### Using with `just` commands
```bash
just utilize discovery-skill      # Produces/updates Feature.md
just utilize architecture-skill   # Produces/updates ADRs, Feature.md links
just flow delivery && just loop   # Agents work from PLAN.md/TODO.md
just gate quality                 # Checks FEEDBACK.md gates
```

### Linking documents
- Feature.md links to ADRs: `See ADR-001 for decision`
- PLAN.md links to Feature.md: `Feature Brief: Feature.md#anchor`
- FEEDBACK.md links to ADRs: `Why? See ADR-001 for rationale`
- CHANGELOG.md links to Feature.md: `See Feature.md#anchor for details`

### Document locations
```
.
├── PLAN.md                    # Epic roadmap
├── TODO.md                    # Sprint tasks
├── CHANGELOG.md               # Release history
├── FEEDBACK.md                # Review results
├── docs/
│   ├── features/
│   │   └── Feature.md         # Feature briefs
│   ├── adr/
│   │   └── ADR-XXX.md         # Architecture decisions
│   └── bdd/
│       └── scenarios.feature  # Acceptance criteria
```

---

## Document Lifecycle

```
1. Issue/Request arrives
   ↓
2. Discovery Skill creates Feature.md
   ↓
3. Architecture Skill creates ADRs (for unknowns)
   ↓
4. Planning Skill creates PLAN.md + TODO.md
   ↓
5. Implementation/Testing Skills execute, update TODO.md + FEEDBACK.md
   ↓
6. Verification Skill checks FEEDBACK.md gates, updates PLAN.md status
   ↓
7. Release Skill creates CHANGELOG.md entry, version tag
   ↓
8. Learning Skill backfeeds to Feature.md if assumptions broken
   ↓
9. [Loop back for next feature/task]
```

This is the document-driven foundation of Springfield Protocol.
