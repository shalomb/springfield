# The Ralph Wiggum Loop: Stateless Resampling

Understanding the core execution engine of the Springfield Protocol.

## What Is It?

A **stateless resampling loop** that ensures quality through persistent iteration rather than one-shot perfection.

**The idea:** Each iteration starts from a clean context. Failures inform strategy but don't accumulate. Keep iterating until the plan is satisfied.

---

## The Loop in 6 Steps

```
┌─────────────────────────────────────┐
│  Monitor PLAN.json for tasks        │ (Source of truth)
└──────────────┬──────────────────────┘
               │
        ┌──────▼──────┐
        │ Spawn agent │ (ephemeral context)
        │ in clean    │ (git worktree or
        │ worktree    │  isolated env)
        └──────┬──────┘
               │
        ┌──────▼──────────┐
        │ Execute task    │ (TDD: Red-Green-Refactor)
        │ (Strict TDD)    │ (Honest effort)
        └──────┬──────────┘
               │
        ┌──────▼──────────┐
        │ Verify results  │ (Coverage, security,
        │ (Herb/Bart)     │  edge cases, perf)
        └──────┬──────────┘
               │
        ┌──────▼──────────┐
        │ Update PLAN.md  │ (Mark task as
        │                 │  verified/failed)
        └──────┬──────────┘
               │
        [Loop back for next task]
```

### 1. Monitor PLAN.md
- PLAN.md is the source of truth
- Tracks epics and their validation state
- Scheduler checks: are there failed or unstarted tasks?

### 2. Spawn Ephemeral Agent
- Create a clean environment for this iteration
- **Git worktree:** Isolated filesystem for this task
- **Or:** Container/sandbox/fresh clone
- **Why clean context?**
  - Prevents accumulation of errors
  - No context rot or hallucination drift
  - Each iteration starts with fresh reasoning

### 3. Execute Task (Strict TDD)
Ralph implements using the TDD cycle:
- **Red:** Write a failing test
- **Green:** Write minimal code to pass
- **Refactor:** Clean up (extract, simplify, optimize)
- **Commit:** Small, atomic git commit
- Repeat until task is complete

### 4. Verify Results
Post-execution verification by Herb (or Bart for security):
- **Coverage:** 95%+ code coverage minimum
- **Security:** Adversarial review for vulnerabilities
- **Edge cases:** Does it handle null, empty, errors?
- **Performance:** Is it acceptable?
- **Regressions:** Did we break anything?

If verification **fails:**
- Task marked as "failed" in PLAN.md
- Loop back to step 1
- Ralph gets another attempt with fresh context

### 5. Update PLAN.md
- Task marked as "verified" (passed all checks)
- Or task marked as "failed" with reason
- For failures: root cause documented
- Strategy for next iteration noted

### 6. Loop Back
- Scheduler picks next unstarted task
- Repeat until PLAN.md shows all tasks verified

---

## Why This Design?

### Problem It Solves: Context Rot

In traditional approaches:
```
Iteration 1: Implement feature A
  ✓ Passes tests
  ✓ But slight design issue introduced

Iteration 2: Implement feature B  
  (context still has design issue from A)
  ✓ Works, but compounds issue
  ✗ Errors accumulate

Iteration N: Systems failing
  (Why? Context rot. Hard to untangle.)
```

**With Ralph Wiggum Loop:**
```
Iteration 1: Implement feature A
  (clean context)
  ✓ Passes tests
  ✗ Design issue caught by Herb
  → Loop back with fresh context

Iteration 2: Implement feature A (again)
  (clean context, informed by failure)
  ✓ Tests pass
  ✓ Design is solid
  ✓ Verified
  → Move to next task

Iteration 3: Implement feature B
  (clean context, no accumulated issues)
  ✓ Works cleanly
```

### Benefits

**Quality Through Persistence**
- Don't aim for perfection on first try
- Fail fast, learn, iterate
- Quality emerges through repeated attempts

**No Hallucination Accumulation**
- Fresh context each iteration
- Agent can't contradict itself from earlier in the context
- Reasoning starts from facts, not previous mistakes

**Visible Progress**
- PLAN.md shows steady march toward completion
- Each task either passes or fails (not vague)
- Easy to see where you are

**Graceful Degradation**
- One failed iteration doesn't cascade
- Next iteration learns from failure without inheriting errors
- Dead-end explorations get pruned

---

## The Human Interpretation

Ralph's character trait: **"I'm not the smartest, but I'll keep trying."**

The loop embodies this:
- Ralph isn't trying to be perfect
- He's trying to be persistent
- Failures are just opportunities to try again
- Eventually, through iteration, quality emerges

**Famous Ralph quote:** "I'm learnding!" 

That's the spirit of this loop.

---

## Implementation Notes

### PLAN.md Format (Example)

```markdown
# PLAN.md

## Feature: User Authentication

### Epic 1: Login System
- [ ] Task 1: Create login endpoint
  - Status: Verified ✓
  - Coverage: 97%
  
- [ ] Task 2: Password hashing
  - Status: Failed (no error handling)
  - Next iteration: Add validation tests
  
- [ ] Task 3: Session management
  - Status: Started
  - Current iteration: Implementing...

### Epic 2: Registration
- [ ] Task 4: Create registration endpoint
  - Status: Unstarted
  - Awaiting: Epic 1 completion
```

### Worktree Example (Git)

```bash
# Main branch has PLAN.md
main: "Task 1: Create login endpoint - UNSTARTED"

# For each task, create a worktree
git worktree add login-endpoint main
cd login-endpoint

# Ralph works here (clean context)
# ... TDD implementation ...
# ... commit ...

# Back to main
cd ..
git worktree remove login-endpoint

# Update PLAN.md
main: "Task 1: Create login endpoint - VERIFIED ✓"

# Create new worktree for next task
git worktree add password-hashing main
# ... repeat ...
```

### When Verification Fails

```bash
# Task 1 failed verification (coverage too low)
main: "Task 1: Create login endpoint - FAILED (78% coverage)"

# Create new worktree for retry
git worktree add login-endpoint-retry main
cd login-endpoint-retry

# Ralph adds missing tests
# ... coverage now 95% ...
# ... commit ...

cd ..
# Update PLAN.md
main: "Task 1: Create login endpoint - VERIFIED ✓"
```

---

## Real-World Example: REST API Feature

### PLAN.md (Initial)
```
## Feature: Search API

- [ ] Task 1: Create /search endpoint
- [ ] Task 2: Query parsing
- [ ] Task 3: Result formatting
- [ ] Task 4: Error handling
- [ ] Task 5: Performance optimization
```

### Iteration 1: Task 1
```
Ralph: Creates /search endpoint
- Writes test: GET /search with query param
- Implements endpoint
- Test passes ✓

Herb: Verifies
- Coverage: 95% ✓
- Edge cases: Missing params handled? ✓
- Perf: Acceptable? ✓

PLAN.md: Task 1 → VERIFIED ✓
```

### Iteration 2: Task 2
```
Ralph: Query parsing
- Test: Parse "foo AND bar" query
- Implement parser
- Tests pass ✓

Herb: Verifies
- Coverage: 93% ✓
- Edge cases: Empty query? ✗ (not handled)

PLAN.md: Task 2 → FAILED (missing edge case)
```

### Iteration 3: Task 2 (Retry)
```
Ralph: Query parsing (fresh context)
- Previous attempt failed on empty query
- Adds test for empty query
- Adds handling in parser
- All tests pass ✓

Herb: Verifies
- Coverage: 98% ✓
- Edge cases: All handled ✓

PLAN.md: Task 2 → VERIFIED ✓
```

### And so on...

Each task iterates until verified. The loop is relentless and fair: keep trying until it works.

---

## Contrast with Other Approaches

| Aspect | Waterfall | Agile (Sprint) | Ralph Wiggum Loop |
|--------|-----------|----------------|-------------------|
| **Context** | Persistent | Persistent (sprint) | Ephemeral |
| **Iteration** | Linear (phases) | Timeboxed | Task-based |
| **Failure Mode** | Whole phase fails | Sprint fails | Task retried |
| **Plan Updates** | At phase gates | Sprint review | After each task |
| **Context Rot** | Severe | Moderate | Minimal |
| **Learning** | Post-phase | Post-sprint | Per-task |

---

## See Also

- **principles.md** - "Iteration over Perfection"
- **LOOP_CATALOG.md** (root) - Section 4.2 for full technical spec
- **QUICK_START.md** (root) - Quick reference for the loop
- **how-to/implement-feature.md** - Practical guide using the loop
