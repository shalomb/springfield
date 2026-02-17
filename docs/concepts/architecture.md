# Springfield Protocol: Architecture Overview

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                   SPRINGFIELD PROTOCOL v0.2                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  AGENTS (2 types)              SKILLS (9 capabilities)         │
│  ├─ Human agents               ├─ discovery-skill              │
│  └─ Computational agents       ├─ architecture-skill           │
│                                ├─ planning-skill               │
│                                ├─ implementation-skill         │
│                                ├─ testing-skill                │
│                                ├─ review-skill                 │
│                                ├─ verification-skill           │
│                                ├─ release-skill                │
│                                └─ learning-skill               │
│                                                                 │
│  EXECUTION ENGINE (Ralph Wiggum Loop)                          │
│  ├─ Monitor PLAN.md                                            │
│  ├─ Spawn agent with clean context                            │
│  ├─ Agent exercises skills                                    │
│  ├─ Verify results                                            │
│  ├─ Update documents                                          │
│  └─ Loop until verified                                       │
│                                                                 │
│  FLOW PATTERNS (2 diamonds)                                    │
│  ├─ Discovery Diamond (Design Thinking)                       │
│  │  └─ Investigate → Validate → Feature Brief                │
│  └─ Delivery Diamond (Agile)                                  │
│     └─ Plan & Build → Verify → Release                        │
│                                                                 │
│  SHARED STATE (7 documents)                                    │
│  ├─ PLAN.md (epic roadmap)                                    │
│  ├─ TODO.md (sprint tasks)                                    │
│  ├─ Feature.md (requirements)                                 │
│  ├─ ADRs (decisions)                                          │
│  ├─ scenarios.feature (acceptance criteria)                   │
│  ├─ FEEDBACK.md (review results)                              │
│  └─ CHANGELOG.md (release history)                            │
│                                                                 │
│  INTERFACE (just CLI)                                          │
│  ├─ just impersonate {agent}                                 │
│  ├─ just utilize {skill}                                     │
│  ├─ just flow {phase}                                        │
│  ├─ just gate {checkpoint}                                   │
│  └─ just loop                                                │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Data Flow: Issue to Release

```
GitHub Issue / Request
        ↓
    [Discovery Diamond]
        ↓
   Human Agent exercises:
   - discovery-skill (investigate)
   - architecture-skill (validate)
   ↓ produces
   Feature.md + ADRs
   ↓
   🚪 Gate: Problem clear?
   ↓
   [Delivery Diamond]
   ↓
   Planning Agent exercises:
   - planning-skill
   ↓ produces
   PLAN.md + TODO.md
   ↓
   Ralph Wiggum Loop:
   ├─ Spawn Implementation Agent
   │  ├─ planning-skill (understand)
   │  ├─ implementation-skill (build)
   │  ├─ testing-skill (test)
   │  ├─ review-skill (review)
   │  ↓ produces
   │  Code + tests
   │  ↓ updates
   │  TODO.md + FEEDBACK.md
   │
   ├─ Spawn Verification Agent
   │  ├─ verification-skill (gate check)
   │  ↓ updates
   │  PLAN.md (status → verified)
   │
   └─ [Loop until all tasks verified]
   ↓
   🚪 Gate: Ready to release?
   ↓
   Release Agent exercises:
   - release-skill
   ↓ produces
   CHANGELOG.md + version tag
   ↓
   Learning Agent exercises:
   - learning-skill (capture insights)
   ↓ updates
   Feature.md + CHANGELOG.md
   ↓
   Released ✓
```

---

## Document State Machine

```
DISCOVERY PHASE
───────────────

Feature.md (state: draft)
  ├─ problem: [blank]
  ├─ requirements: [blank]
  ├─ constraints: [discovered]
  ├─ unknowns: [identified]
  └─ assumptions: [list]
       ↓
  Discovery Skill exercises
       ↓
Feature.md (state: proposed)
  ├─ problem: ✓
  ├─ requirements: ✓
  ├─ constraints: ✓
  └─ unknowns: [list of ADRs needed]
       ↓
  Architecture Skill exercises
       ↓
ADRs (state: created)
  ├─ ADR-001: [decision]
  ├─ ADR-002: [decision]
  └─ ...
       ↓
Feature.md (state: complete)
  ├─ unknowns: [all linked to ADRs]
  └─ assumptions: [explicit list]


DELIVERY PHASE
──────────────

PLAN.md (state: created)
  ├─ Epic 1: unstarted
  ├─ Epic 2: unstarted
  └─ ...
       ↓
TODO.md (state: populated)
  ├─ Task 1: unstarted
  ├─ Task 2: unstarted
  └─ ...
       ↓
Ralph Wiggum Loop iteration 1
       ↓
TODO.md + FEEDBACK.md (state: in-progress)
  ├─ Task 1: in-progress
  ├─ Task 1 issues found
  └─ ...
       ↓
Ralph Wiggum Loop iteration 2
       ↓
PLAN.md (state: verified)
  ├─ Task 1: verified ✓
  ├─ Task 2: in-progress
  └─ ...
       ↓
[Loop continues...]
       ↓
PLAN.md (state: complete)
  └─ All tasks: verified ✓


RELEASE PHASE
─────────────

CHANGELOG.md (state: created)
  ├─ [Unreleased]
  │  ├─ Added: [features from PLAN.md]
  │  ├─ Fixed: [from FEEDBACK.md]
  │  └─ Learning: [what surprised us]
  └─ ...
       ↓
Release Skill exercises
       ↓
CHANGELOG.md (state: released)
  ├─ [1.0.0] - [date]
  │  ├─ Added: [features]
  │  └─ Learning: [captured]
  └─ ...
```

---

## Skill Execution Model

```
Agent exercises a skill:

┌─────────────────────────────────────────┐
│ Skill Definition (.github/skills/*)    │
│                                         │
│  SKILL.md                              │
│  ├─ Purpose                            │
│  ├─ Inputs (documents to read)         │
│  ├─ Procedure (steps)                  │
│  ├─ Outputs (documents to produce)     │
│  └─ Examples                           │
└─────────────────────────────────────────┘
        ↓
┌─────────────────────────────────────────┐
│ Agent Context                           │
│                                         │
│  ├─ Available documents                │
│  ├─ Current task                       │
│  ├─ Previous results                   │
│  └─ Learning from prior iterations     │
└─────────────────────────────────────────┘
        ↓
┌─────────────────────────────────────────┐
│ Skill Exercise                          │
│                                         │
│  1. Read inputs (PLAN.md, Feature.md)  │
│  2. Execute procedure                  │
│  3. Update outputs (TODO.md, FEEDBACK) │
│  4. Capture learning                   │
│  5. Report results                     │
└─────────────────────────────────────────┘
        ↓
┌─────────────────────────────────────────┐
│ Results                                 │
│                                         │
│  ├─ New/updated documents              │
│  ├─ Learning captured                  │
│  └─ Next action (proceed/loop back)    │
└─────────────────────────────────────────┘
```

---

## Ralph Wiggum Loop (Detailed)

```
┌─────────────────────────────────────────────────────────────┐
│  Monitor PLAN.md                                            │
│  for unstarted/failed tasks                                │
└──────────────────────┬──────────────────────────────────────┘
                       │
                ┌──────▼──────┐
                │ Find next   │
                │ unstarted   │
                │ task        │
                └──────┬──────┘
                       │
            ┌──────────▼──────────┐
            │ All tasks verified? │
            └──────┬──────────────┘
                   │
             ┌─────▼─────┐
             │   YES     │ NO
             │ (exit)    │
             └─────┬─────┘
                   │
          ┌────────▼──────────┐
          │ Create ephemeral  │
          │ context           │
          │ (clean slate)     │
          └────────┬──────────┘
                   │
          ┌────────▼──────────┐
          │ Spawn agent       │
          │ (ask which)       │
          └────────┬──────────┘
                   │
          ┌────────▼──────────┐
          │ Load task from    │
          │ TODO.md           │
          └────────┬──────────┘
                   │
       ┌───────────▼────────────┐
       │ Agent exercises skills:│
       │ - planning-skill       │
       │ - implementation-skill │
       │ - testing-skill        │
       │ - review-skill         │
       │                        │
       │ Updates:               │
       │ - TODO.md (progress)   │
       │ - FEEDBACK.md (results)│
       └───────────┬────────────┘
                   │
       ┌───────────▼────────────┐
       │ Verification agent     │
       │ exercises:             │
       │ - verification-skill   │
       │                        │
       │ Checks:                │
       │ - Coverage > 95%?      │
       │ - All tests pass?      │
       │ - Security OK?         │
       │ - Performance OK?      │
       └───────────┬────────────┘
                   │
         ┌─────────▼─────────┐
         │ Results OK?       │
         └────┬──────────┬───┘
              │ PASS     │ FAIL
              │          │
         ┌────▼────┐  ┌──▼──────────┐
         │ Mark    │  │ Capture why │
         │ task    │  │ Update TODO │
         │ verified│  │ Loop back   │
         │ in      │  │ with fresh  │
         │ PLAN.md │  │ context     │
         └────┬────┘  └──┬──────────┘
              │          │
              └──────┬───┘
                     │
                [Loop back to top]
```

---

## Team Scaling

```
SOLO DEVELOPER
├─ 1 Human Agent
└─ Exercises all 9 skills as needed


SMALL TEAM (3-5)
├─ Multiple Human Agents
│  ├─ Frontend dev (impl, testing, review)
│  ├─ Backend dev (impl, testing, arch)
│  └─ Lead (discovery, planning, verification)
└─ Computational Agents
   ├─ testing-agent (automated tests)
   └─ release-agent (automated releases)


MEDIUM TEAM (10-20)
├─ Human Agents (organized by domain)
│  ├─ Discovery Team
│  │  └─ discovery, architecture, learning
│  ├─ Development Team
│  │  └─ planning, implementation, testing, review
│  └─ QA Team
│     └─ testing, verification, learning
└─ Computational Agents
   ├─ security-agent
   ├─ performance-agent
   └─ release-agent


ENTERPRISE (50+)
├─ Human Agents (teams per product)
│  ├─ Each team exercises same 9 skills
│  ├─ Specialists for each skill
│  └─ Coordinated via shared documents
└─ Computational Agents (specialized)
   ├─ security-agent (cross-team)
   ├─ performance-agent (cross-team)
   ├─ devops-agent (infra)
   └─ monitoring-agent (prod)

KEY: Same skills, different team structures, different agent combinations.
```

---

## Technology Stack

```
Documents (Shared State)
├─ PLAN.md ......................... Git + Markdown
├─ TODO.md ......................... Git + Markdown
├─ Feature.md ...................... Git + Markdown
├─ ADRs ............................ Git + Markdown
├─ scenarios.feature ............... Git + Gherkin
├─ FEEDBACK.md ..................... Git + Markdown
└─ CHANGELOG.md .................... Git + Markdown

Skills (.github/skills/)
├─ SKILL.md ........................ Markdown + instructions
├─ examples/ ....................... Markdown examples
├─ tools/ .......................... Bash/Python scripts
└─ exercise.sh ..................... Executable procedure

Agents (.github/agents/)
├─ agent.md ........................ Markdown definition
└─ config.yaml ..................... Configuration (optional)

CLI Interface
└─ justfile ........................ Just commands

Workflows (optional automation)
└─ .github/workflows/ .............. GitHub Actions
   ├─ tests.yml .................... Run tests
   ├─ coverage.yml ................. Check coverage
   ├─ security.yml ................. Security scan
   └─ release.yml .................. Auto-release
```

---

## Integration Points

```
GitHub
├─ Issues → Discovery Phase
├─ Pull Requests ← FEEDBACK.md
├─ Workflows ← .github/workflows/
└─ Releases ← Release Skill

Git
├─ Commits ← Implementation Skill
├─ Tags ← Release Skill
└─ Branches ← Ephemeral contexts

Local Development
├─ justfile → CLI interface
├─ .github/skills/ → Skill definitions
└─ .github/agents/ → Agent definitions

CI/CD Pipeline
├─ Tests → testing-skill
├─ Coverage → verification-skill
├─ Security → review-skill
└─ Deployment → devops-agent
```

This is the complete architecture of Springfield Protocol v0.2.
