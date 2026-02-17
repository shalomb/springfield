# Springfield Protocol v0.2: The Complete Model

## Overview

Springfield Protocol is a **document-driven, agent-based framework** for coordinating work across design thinking, agile delivery, and DevOps.

**Core Components:**
- 5 focused agents (single-pizza team)
- 9 standardized skills
- 2 flow diamonds (design thinking + agile)
- 1 iteration engine (Ralph Wiggum Loop)
- 7 core documents
- 1 CLI UX (just commands)

---

## The 5-Agent Team

We use a "Single Pizza Team" model to keep agent context windows focused and expedient.

### 1. Product Agent (The "What & Why")
- **Focus:** Discovery Diamond (Diverge)
- **Skills:** `discovery-skill` (investigate), `wiggum` (triage)
- **Role:** Understand user needs, define problems, create Feature Briefs.
- **Context:** Loaded with user research, gemba walk, and problem definition prompts.

### 2. Planning Agent (The "How & Structure")
- **Focus:** Discovery Converge & Delivery Planning
- **Skills:** `planning-skill` (breakdown), `architecture-skill` (patterns)
- **Role:** Turn briefs into executable plans, validate architecture, create ADRs.
- **Context:** Loaded with architectural patterns, dependency graph logic, and breakdown strategies.

### 3. Build Agent (The "Doer")
- **Focus:** Delivery Diamond (Diverge/Build)
- **Skills:** `implementation-skill` (code), `testing-skill` (TDD)
- **Role:** Write code, write tests, build infrastructure. Optimistic mindset.
- **Context:** Loaded with TDD rules, language syntax, clean code guidelines.

### 4. Quality Agent (The "Critic")
- **Focus:** Delivery Diamond (Converge/Verify)
- **Skills:** `review-skill` (adversarial), `verification-skill` (gates)
- **Role:** Challenge assumptions, find bugs, verify gates. Pessimistic mindset.
- **Context:** Loaded with OWASP security checklists, edge case heuristics, and style guides.

### 5. Release Agent (The "Shipper")
- **Focus:** Release & Loop Feedback
- **Skills:** `release-skill` (ceremony), `learning-skill` (insights)
- **Role:** Manage releases, update changelogs, capture learning.
- **Context:** Loaded with semver rules, changelog formats, and git tagging logic.

---

## The Nine Skills

Skills are the capabilities agents exercise.

1. **discovery-skill** (Investigate, interview, Five Whys) -> *Product Agent*
2. **architecture-skill** (Validate fit, design decisions) -> *Planning Agent*
3. **planning-skill** (Break features into tasks) -> *Planning Agent*
4. **implementation-skill** (Write code, TDD) -> *Build Agent*
5. **testing-skill** (Write tests, verify coverage) -> *Build Agent*
6. **review-skill** (Adversarial review, security) -> *Quality Agent*
7. **verification-skill** (Check quality gates) -> *Quality Agent*
8. **release-skill** (Semantic versioning, publish) -> *Release Agent*
9. **learning-skill** (Capture unknowns, insights) -> *Release/All Agents*

---

## The Two Diamonds

### Discovery Diamond (Design Thinking)

```
Problem/Request
      ↓
   DIVERGE: Investigate
   (Product Agent)
      ↓
🚪 Gate: Problem Clear?
      ↓
   CONVERGE: Validate
   (Planning Agent)
      ↓
Feature Brief (Feature.md + ADRs)
```

### Delivery Diamond (Agile)

```
Feature Brief
      ↓
   DIVERGE: Plan & Build
   (Planning Agent -> Build Agent)
      ↓
🚪 Gate: Quality Ready?
      ↓
   CONVERGE: Verify
   (Quality Agent)
      ↓
Release (CHANGELOG.md)
   (Release Agent)
```

---

## The Ralph Wiggum Loop (Execution Engine)

The **stateless resampling** loop that runs the framework:

```
Monitor PLAN.md for unstarted/failed tasks
           ↓
Spawn Agent (Build or Quality)
with clean context (ephemeral environment)
           ↓
Agent exercises skills
           ↓
Verification (Quality Agent)
           ↓
If passed: Mark verified
If failed: Loop back with fresh context
```

---

## The Seven Core Documents

**Shared State:**
1. **PLAN.md** (Roadmap)
2. **TODO.md** (Tasks)
3. **Feature.md** (Requirements)
4. **ADRs** (Decisions)
5. **scenarios.feature** (BDD Specs)
6. **FEEDBACK.md** (Review results)
7. **CHANGELOG.md** (History)

---

## Core Principles

1. **Plan Before You Build** (Discovery first)
2. **Steer As You Go** (Adjust based on learning)
3. **Iteration Over Perfection** (Ralph Wiggum Loop)
4. **Explicit Uncertainty** (Document unknowns)
5. **Documents Are Shared State** (Markdown is truth)
6. **Avoid the Distracted Agent** (Focused contexts)

---

## Why This Model Works

✅ **Focused Contexts** - 5 specialized agents prevent context window overload.
✅ **Clear Separation** - Builders build, Critics critique. No hallucinations of quality.
✅ **Scalable** - Works for 1 person (wearing 5 hats) to large teams.
✅ **Coherent** - Aligns with industry standard Design Thinking + Agile flows.
