# How to Implement a Feature

Step-by-step guide for taking a feature from plan to release.

## The Workflow

```
Feature Brief (validated)
    ↓
@lisa: Break into tasks → PLAN.md + TODO.md
    ↓
@ralph: Implement task 1 (TDD loop)
    ↓
@bart: Review for security & quality
    ↓
@herb: Verify coverage (95%+)
    ↓
@marge: Check user alignment
    ↓
@lovejoy: Release
```

## Step 1: Plan the Work (@lisa)

**Who:** Lisa (Strategic Planner)

**What to do:**
- Review the Feature Brief
- Break into executable tasks
- Create PLAN.md with epics
- Create TODO.md with specific tasks
- Document dependencies

**Output:** 
- `PLAN.md` - Roadmap of work
- `TODO.md` - Tasks for Ralph to execute

**When you're done:** Feature Brief is ready for implementation.

---

## Step 2: Implement (@ralph)

**Who:** Ralph (TDD Executor)

**What to do:**
- Pick a task from TODO.md
- Run the Ralph Wiggum Loop:
  1. **Red** - Write failing test
  2. **Green** - Write minimal code to pass
  3. **Refactor** - Clean up implementation
  4. **Commit** - Small, atomic git commits
- Repeat until task is complete
- Update TODO.md

**Output:**
- Tested code
- git commits (one per task)

**When you're done:** Code is working and tested. Move to next task.

---

## Step 3: Adversarial Review (@bart)

**Who:** Bart (Adversarial Reviewer)

**What to do:**
- Review code for:
  - Security vulnerabilities (CSRF, injection, auth)
  - Edge cases (null, empty, overflow)
  - Error handling
  - Performance bottlenecks
  - API misuse
- Document findings
- Challenge assumptions

**Output:**
- List of issues to fix
- Questions for Ralph

**When you're done:** Ralph fixes issues and resubmits.

---

## Step 4: Verify Quality (@herb)

**Who:** Herb (Quality Engineer)

**What to do:**
- Check test coverage (must be 95%+)
- Run quality checks:
  - Linting
  - Type checking
  - Static analysis
- Verify no regressions
- Sign off on quality

**Output:**
- Coverage report
- Quality checklist (passed/failed)

**When you're done:** Code meets quality standards.

---

## Step 5: User Alignment (@marge)

**Who:** Marge (Empathy & Guardrails)

**What to do:**
- Review feature against original needs
- Conduct user feedback session (if possible)
- Verify scope didn't drift
- Check roadmap alignment
- Document any changes needed

**Output:**
- User feedback
- Go/no-go recommendation

**When you're done:** Feature is aligned with user needs.

---

## Step 6: Release (@lovejoy)

**Who:** Lovejoy (Release Master)

**What to do:**
- Bump version (semantic versioning)
- Create CHANGELOG entry
- Tag release in git
- Publish to registry
- Announce release

**Output:**
- git tag (v1.2.3)
- CHANGELOG.md entry
- Published package

**When you're done:** Feature is live!

---

## Troubleshooting

**Coverage is low**
- Herb reports gaps
- Ralph adds tests to missing areas
- Loop back to step 2

**Bart finds security issues**
- Ralph fixes issues
- Loop back to step 3

**User feedback is negative**
- Marge documents concerns
- Loop back to step 2 (redesign)

**Release fails**
- Lovejoy investigates
- Fix in code, loop back to step 2

---

## Files Involved

- `PLAN.md` - Epic roadmap
- `TODO.md` - Task list
- `src/` - Implementation
- `tests/` - Test cases
- `CHANGELOG.md` - Release notes

## See Also

- **LOOP_CATALOG.md** - Full Ralph Wiggum Loop specification
- **CHARACTER_SKILLS.md** - Detailed agent profiles
- **QUICK_START.md** - Quick reference
