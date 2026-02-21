# Feature Brief: Justfile Agent Command Migration to Springfield Binary

**Prepared for:** Marge Simpson (Product Agent)  
**Date:** 2026-02-20  
**Status:** Ready for Discovery Gate Review  
**Analyst:** Claude Code

---

## Problem Statement

### The User's Problem
The Springfield team needs to **retire shell-based agent invocation** from the Justfile and move to the **type-safe Go binary**. The current system has two parallel orchestration patterns:

1. **Shell pattern (Justfile):** `just ralph`, `just lisa`, `just bart`, `just lovejoy`
   - Invokes agents via `npm exec @mariozechner/pi-coding-agent`
   - Loop logic is hardcoded in bash
   - Relies on grep to parse FEEDBACK.md and TODO.md
   - Cannot be unit tested
   - Fragile to LLM output format changes

2. **Binary pattern (Go):** `springfield orchestrate` (EPIC-009, just completed)
   - Type-safe state machine
   - Uses `td(1)` as source of truth
   - Fully tested and approved by Bart
   - Delegates agent execution to `springfield agent` command (not yet fully wired)

**The Pain:** Developers must mentally context-switch between two orchestration systems. EPIC-009 proved the Go binary works, but we haven't connected the "last mile" — the individual agent commands (ralph, lisa, bart, lovejoy) still need to move from shell to binary.

### Root Cause
When EPIC-009 was scoped, the decision was made to migrate the high-level `just do` orchestration loop first, leaving the individual agent commands for a follow-up epic. Now that EPIC-009 is complete and approved, the **technical precondition for migration is satisfied** — there's no blocker preventing Ralph from moving the agent commands to Go.

### Why Now?
- ✅ EPIC-009 passed final review (FEEDBACK.md: "Ready for merge/release")
- ✅ No active epic depends on the shell versions of these commands
- ✅ All prompts are stable (no changes in last 5+ commits)
- ✅ `config.toml` system is in place to manage prompts
- ⚠️ **Risk:** Keeping two orchestration patterns creates ongoing maintenance burden and confuses new contributors

---

## User Needs & Acceptance Criteria

### Who Benefits?
1. **Ralph (Build Agent)** — Can test loop logic in unit tests, doesn't rely on grep or string matching
2. **Lisa (Planning Agent)** — Can use td-based state instead of file grepping
3. **Bart (Quality Agent)** — Can log decisions to `td log --decision` instead of parsing prose FEEDBACK.md
4. **Developers** — Single orchestration pattern, no context switching, clear source of truth (td)
5. **Future Contributors** — Can understand orchestration logic by reading Go code, not bash

### Business Acceptance Criteria
- [ ] **Single Orchestration Source:** No developer ever asks "do I use Justfile or binary?" — the answer is always "binary"
- [ ] **Testability:** All loop logic is unit-testable (migration succeeds if tests cover >90% of new code)
- [ ] **No Regression:** `just ralph`, `just lisa`, `just bart`, `just lovejoy` still work exactly as before (to end user)
- [ ] **Documentation Clarity:** AGENTS.md and README explain the architecture without mentioning Justfile recipes
- [ ] **Cost Control:** No additional LLM calls, no new build steps, no npm overhead (binary is smaller/faster)

---

## Scope Definition

### In Scope (What We're Building)
This epic consists of **4 phases**, each 1 sprint:

**Phase 1: Ralph Command Migration**
- Move Ralph's multi-iteration loop logic from Justfile to `internal/agent/ralph_runner.go`
- Implement TODO.md parsing and iteration control in Go
- Add unit tests for loop logic, git status checking, graceful termination
- Update Justfile recipe to thin wrapper: `just ralph` → `springfield agent ralph "work on <task>"`

**Phase 2: Lisa Command Migration**
- Move Lisa's planning prompt and FEEDBACK.md integration to `internal/agent/lisa_runner.go`
- Replace file grepping with td-based state updates
- Add tests for PLAN.md integration and feedback parsing

**Phase 3: Bart & Lovejoy**
- Migrate Bart's verdict checking to `td log --decision`
- Move Lovejoy's merge logic into orchestrator (not individual command)

**Phase 4: Cleanup**
- Remove npm dependency from Justfile header
- Delete all shell recipe versions
- Update AGENTS.md to document new architecture
- Write migration guide for documentation

### Out of Scope (What We're NOT Building)
- ❌ Changing how agents work (still invoke pi-coding-agent, same prompts)
- ❌ Changing PLAN.md, TODO.md, FEEDBACK.md formats (they stay as markdown)
- ❌ Changing orchestrator state machine (EPIC-009 already proved this works)
- ❌ Rewriting td integration (already in place)
- ❌ Changing the quality gates or review process

---

## WSJF Scoring

- **Value:** High (eliminates ongoing maintenance debt, closes EPIC-009)
- **Urgency:** Medium (no active epic blocked, but we should do it soon to avoid "stale shell code")
- **Complexity:** Medium (phased approach, no unknown unknowns)
- **Size:** 5-6 tasks spread over 4 sprints (moderate work)

**Estimated WSJF Score: 3.5** (CoD: 21 / Size: 6)
- **Comparison:** EPIC-005 (Governance) = 3.25, EPIC-006 (Compatibility) = 2.0
- **Recommendation:** Schedule for **next sprint** (after current epic completes)

---

## Success Metrics & Definition of Done

### How We Measure Success
1. **Test Coverage:** All new code has ≥90% coverage
2. **Behavior Parity:** `just ralph` produces identical output before/after migration
3. **Reliability:** Zero regressions when Ralph uses the binary version
4. **Documentation:** AGENTS.md is updated, no dangling references to Justfile recipes
5. **Git Quality:** All commits follow ACP (Atomic Commit Protocol)

---

**Prepared by:** Claude Code  
**Date:** 2026-02-20  
**Status:** Awaiting Marge's Discovery Gate Decision
