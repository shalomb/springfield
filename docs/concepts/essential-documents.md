# The 7 Essential Documents

Springfield Protocol uses 7 core documents as the **shared state** that agents and skills read from and write to.

---

## Document Stack Overview

```
DISCOVERY PHASE              DELIVERY PHASE              RELEASE PHASE
──────────────────          ──────────────────          ─────────────────

Feature.md                  PLAN.md + TODO.md           CHANGELOG.md
ADRs                        FEEDBACK.md                 (+ learning capture)
```

All documents work together. **Documents ARE the shared state.** No separate database, no message queue—just markdown files that agents and skills read and update.

---

## 1. PLAN.md — Epic Roadmap

**Purpose:** High-level roadmap showing what we decided to build and current status

**Owns:**
- Epics and their status
- Task assignments and dependencies
- Links to requirements (Feature.md) and decisions (ADRs)
- Current progress

**Updated by:** Planning skill (creates), all skills (update status)

**Read by:** Ralph Wiggum Loop (to find next task), all agents (to understand current work)

**Example:**
```markdown
# PLAN.md - Feature Roadmap

## Feature: User Authentication

### Epic 1: Login System
- [ ] Task 1: Create login endpoint
  - Status: unstarted | in-progress | verified | failed
  - Assignee: implementation-agent
  - ADR: ADR-001 (session strategy)
  - Feature Brief: Feature.md#login
  - BDD: scenarios.feature#login

- [x] Task 2: Password hashing
  - Status: verified ✓
  - Coverage: 97%
```

---

## 2. TODO.md — Sprint Tasks

**Purpose:** Executable task list for current sprint

**Owns:**
- Immediate work queue
- Task progress and subtasks
- Dependencies and blockers
- Learning captured during work

**Updated by:** Implementation skill, testing skill (during work)

**Read by:** Implementation agents (what to work on next)

**Example:**
```markdown
# TODO.md - Current Sprint

- [ ] Task 1: Create login endpoint
  - Assigned to: implementation-agent
  - Subtasks:
    - [ ] Write failing test for GET /login
    - [ ] Implement endpoint
    - [ ] Add rate limiting
  - Blocked by: None
  - Time estimate: 4 hours
  - Started: 2026-02-17
  - Learning: "Password validation more complex than expected"
```

---

## 3. Feature.md — Feature Brief

**Purpose:** Feature requirements, constraints, and decision context

**Owns:**
- Problem statement (why we're building this)
- Requirements (what users need)
- Constraints (hard limits)
- Explicit unknowns (linked to ADRs where resolved)
- Explicit assumptions (what could break)
- Scope boundaries

**Updated by:** Discovery skill (initial), learning skill (when assumptions break)

**Read by:** All agents (validate requirements), implementation agents (know what to build)

**Example:**
```markdown
# Feature.md - User Authentication

## Problem
Users cannot log in. Needed for enterprise customers requiring account isolation.
Root cause: Authentication system doesn't exist.

## Requirements
As a user, I want to log in with email/password so I can access my account.

## Acceptance Criteria
See: scenarios.feature#login

## Constraints & Unknowns
- Must use bcrypt for password hashing
- Session timeout: 24 hours
- Rate limit: 5 failed attempts per IP
- **Unknown:** Session storage (Database vs Redis) - See ADR-001 for decision
- **Assumption:** Users have email addresses (from user profile feature)
- **Assumption:** HTTPS is always enabled

## Scope
✅ Email/password login
✅ Session management
❌ OAuth/SSO (future)
❌ Password reset (future)

## Success Criteria
- Users can create account and log in
- Sessions persist across page reloads
- Auto-logout after 24 hours
```

**Key insight:** Problems, unknowns, and assumptions are all here, not in separate documents.

---

## 4. ADRs — Architecture Decision Records

**Purpose:** Record architectural decisions with full thinking

**Location:** `docs/adr/ADR-XXX-title.md`

**Owns:**
- Architectural decisions with full rationale
- Unknowns from Feature.md that were resolved here
- Alternatives considered (captures the thinking)
- Tradeoffs and consequences

**Updated by:** Architecture skill

**Read by:** Implementation skill (follow ADR), review skill (validate ADR compliance)

**Example:**
```markdown
# ADR-001: Session Storage Strategy

**Date:** 2026-02-17
**Status:** Accepted
**Related Feature:** Feature.md#login
**Related Unknown:** Feature.md#unknowns

## Problem
How should we store user sessions?

## Decision
Use Redis for session storage with database fallback.

## Rationale
- Redis: Fast, temporary data, scales well
- Fallback: Database for persistence across restarts
- Tradeoff: Extra complexity vs better performance

## Consequences
- Sessions don't survive app restart without database
- Need Redis client library (ioredis)
- Must implement cleanup cron for expired sessions

## Alternatives Considered
1. Database only (simpler, slower)
2. In-memory only (fast, lost on restart)
3. Redis + database ✅ (chosen)

## Implementation Notes
- Store session ID in Redis
- Store user ID → session ID in database
- TTL: 24 hours on Redis keys
```

---

## 5. scenarios.feature — BDD Specs

**Purpose:** Executable acceptance criteria

**Location:** `features/[domain]/[feature].feature`

**Owns:**
- Acceptance criteria (executable as tests)
- Test cases for all scenarios
- Constraints documented as scenarios (rate limiting, timeouts, edge cases)

**Updated by:** Feature specification + implementation

**Read by:** Testing skill (execute tests), implementation skill (know what to code)

**Example:**
```gherkin
# features/authentication/login.feature

Feature: User Login
  Problem: Users need to authenticate
  Feature: Feature.md#login
  
  Scenario: Valid login
    Given I have an account with email "user@example.com"
    When I navigate to login and enter credentials
    Then I should be logged in
    And I should see my dashboard

  Scenario: Rate limiting
    Given I have an account
    When I attempt login 6 times with wrong password
    Then I should see error "Too many attempts"
    And my IP should be rate limited for 15 minutes

  Scenario: Session timeout
    Given I am logged in
    When 24 hours pass
    And I refresh the page
    Then I should be logged out
    And redirected to login
```

---

## 6. FEEDBACK.md — Review & Gate Results

**Purpose:** Code review feedback, questions, and gate status

**Owns:**
- Code review issues and questions
- Gate check results (coverage, security, performance)
- Blocking vs. ready-to-merge status
- Learning captured from review

**Updated by:** Review skill, verification skill

**Read by:** Implementation agents (fix issues), teams (understand gate status)

**Example:**
```markdown
# FEEDBACK.md - Review Results

## Task 1: Create login endpoint
**Reviewed:** 2026-02-17 10:30 UTC
**Reviewer:** review-skill

### Issues Found
- [x] Fixed: Password validation doesn't check common passwords
- [ ] Open: Missing error handler for database timeout
  - Risk: Medium
  - Suggested fix: Implement retry logic

### Questions
- Why Redis fallback instead of just database?
  - Answer: See ADR-001 for rationale (performance vs consistency tradeoff)

### Security Review
- [x] No SQL injection (parameterized queries)
- [x] No XSS (output encoding)
- [x] HTTPS enforced
- [x] Session tokens secure

### Quality Gate Results
- [x] Tests pass: 342/342 ✓
- [x] Coverage: 96.6% (> 95%) ✓
- [x] Security scan: 0 critical issues ✓
- [x] Performance: < 100ms response time ✓

**Status:** READY TO MERGE (1 minor issue)
```

---

## 7. CHANGELOG.md — Release History

**Purpose:** Record what changed and learning from each release

**Owns:**
- What changed in each version
- Features added (linked to Feature.md)
- Bug fixes and improvements
- Learning captured ("what surprised us")
- Security improvements

**Updated by:** Release skill (structure + version), learning skill (learning notes)

**Read by:** Teams (release notes), stakeholders (what's in this version)

**Example:**
```markdown
# CHANGELOG.md

## [Unreleased]

### Added
- User login feature (Feature.md#login)
  - Email/password authentication
  - Session management with 24h timeout
  - Rate limiting on failed attempts

### Fixed
- Session cookie security flags (ADR-001)

### Security
- Password hashing with bcrypt (ADR-001)
- HTTPS enforcement
- Rate limiting on login attempts

### Learning
- Password validation was more complex than expected
  - Need to check against common passwords
  - Need to check against previous passwords
  - Updated Feature.md#constraints
- User ID column needed for audit trail
  - From TODO.md Task 0 learning
  - Added to database schema

---

## [1.0.0] - 2026-02-16

### Added
- Initial release
- Basic user management
```

---

## How They Flow Together

### Discovery Phase
```
Issue arrives
  ↓
Discovery Skill exercises
  ↓ produces
Feature.md (problem + requirements + unknowns)
  ↓
Architecture Skill exercises (on unknowns)
  ↓ produces
ADRs (decisions)
  ↓
Feature.md updated (links to ADRs)
  ↓
Feature Brief ready
```

### Delivery Phase
```
Feature Brief ready
  ↓
Planning Skill exercises
  ↓ produces
PLAN.md + TODO.md
  ↓
Ralph Wiggum Loop iteration:
  ├─ Implementation Skill exercises (builds code)
  ├─ Testing Skill exercises (writes tests)
  ├─ Review Skill exercises (finds issues)
  ├─ Verification Skill exercises (gates check)
  └─ Updates: TODO.md, FEEDBACK.md, PLAN.md
  ↓
Loop until all tasks verified
```

### Release Phase
```
Verified code ready
  ↓
Release Skill exercises
  ↓ produces
CHANGELOG.md + version tag
  ↓
Learning Skill exercises
  ↓ captures
CHANGELOG.md (what surprised us)
  ↓ backfeeds
Feature.md (if assumptions broken)
  ↓
Released
```

---

## Where Everything Lives

### Problem/Root Cause
- **Lives in:** Feature.md#problem
- **Captured by:** Discovery skill
- **Updated by:** Learning skill (if discovered during build)

### Unknowns & Questions
- **Lives in:** Feature.md#unknowns
- **Linked to:** ADRs (where resolved)
- **Captured by:** Discovery + architecture skills
- **Resolved by:** Architecture skill → ADR

### Assumptions
- **Lives in:** Feature.md#assumptions
- **Updated by:** Learning skill (when broken during work)
- **Learning captured:** TODO.md → FEEDBACK.md → CHANGELOG.md

### Decisions Made
- **Lives in:** ADRs (full reasoning)
- **Referenced in:** Feature.md (links), PLAN.md (task links)
- **Validated in:** Review skill (FEEDBACK.md checks compliance)

### Learning Captured
- **During work:** TODO.md (learning notes)
- **During review:** FEEDBACK.md (questions, surprises)
- **During release:** CHANGELOG.md (what changed, what surprised us)
- **Backfeed:** Updates Feature.md if assumptions invalidated

---

## The Minimal Set

If you had to start with just 3 documents:

1. **PLAN.md** — Can't loop without it (task status drives the loop)
2. **Feature.md** — Prevents building wrong thing (requirements + constraints)
3. **ADRs** — Prevents rework (architectural guardrails)

Everything else is derivative or supportive.

But we'd argue **FEEDBACK.md** and **CHANGELOG.md** are essential because they're where learning happens and gets captured.

---

## Agent/Skill Awareness

### Every agent should have access to:
1. **PLAN.md** — Current roadmap
2. **TODO.md** — Current tasks
3. **Feature.md** — Requirements + constraints
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
| **Release** | CHANGELOG.md, tag | - |
| **Learning** | - | Feature.md, TODO.md, FEEDBACK.md, CHANGELOG.md |

---

## Key Principle

**Documents ARE the shared state.**

Not a database. Not a message queue. Just markdown files.

- Agents read documents to understand context
- Skills write documents to capture work
- Documents link to each other (Feature.md → ADRs, PLAN.md → Feature.md)
- Learning is woven into documents, not captured separately

This is the foundation of Springfield Protocol v0.2.
