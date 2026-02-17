# How to Debug an Issue

Step-by-step guide for investigating and fixing a problem.

## The Workflow

```
Issue reported
    ↓
@wiggum: Triage & categorize
    ↓
Search KEDB (Known Error Database)
    ↓
IF found → Document & close
IF not found → Investigate with @bart or @ralph
    ↓
Root cause identified
    ↓
Fix implemented
    ↓
@herb: Verify fix
    ↓
@lovejoy: Release hotfix
```

## Step 1: Triage (@wiggum)

**Who:** Wiggum (Triage Officer)

**What to do:**
- Assess severity (critical, high, medium, low)
- Assess type (bug, performance, security, documentation)
- Check Definition of Ready:
  - Is the issue clear?
  - Can it be reproduced?
  - Do we have enough context?
- Assign priority

**Output:**
- Triaged issue with labels
- Priority assigned

**When you're done:** Issue is ready for investigation.

---

## Step 2: Check KEDB

**Who:** Anyone

**What to do:**
- Search the Known Error Database for similar issues
- Look for:
  - Error messages
  - Stack traces
  - Related symptoms
  - Previous solutions

**If found:**
- Review the documented solution
- Apply it
- Test fix
- Close issue with reference to KEDB entry

**If not found:**
- Continue to Step 3
- Later: Document this as a new KEDB entry

---

## Step 3: Investigate

**Who:** Bart (Adversarial Reviewer) or Ralph (TDD Executor)

**What to do:**

### If it's a subtle bug:
- Use **ReAct Loop**:
  1. Reason: What do we know?
  2. Act: Write a test that reproduces the issue
  3. Observe: Confirm test fails
  4. Repeat until root cause is clear

### If it's a complex issue:
- Use **Tree of Thoughts Loop**:
  1. Generate multiple hypotheses
  2. Evaluate evidence for each
  3. Prune unlikely ones
  4. Explore most promising path

**Output:**
- Reproduction test case
- Root cause identified
- Fix strategy documented

**When you're done:** Root cause is clear and testable.

---

## Step 4: Fix

**Who:** Ralph (TDD Executor)

**What to do:**
- Write a test that captures the fix
- Implement minimal fix
- Verify test passes
- Check for regressions
- Commit with clear message

**Output:**
- git commit with fix
- Updated tests

**When you're done:** Issue is fixed and tested.

---

## Step 5: Verify (@herb)

**Who:** Herb (Quality Engineer)

**What to do:**
- Verify fix works
- Check test coverage for fix area
- Run regression tests
- Confirm no new issues introduced

**Output:**
- Regression test results
- Quality sign-off

**When you're done:** Fix is verified.

---

## Step 6: Release (@lovejoy)

**Who:** Lovejoy (Release Master)

**What to do:**
- Tag as hotfix version (e.g., v1.2.1)
- Create CHANGELOG entry
- Publish immediately
- Notify users

**Output:**
- Hotfix released
- Users notified

**When you're done:** Issue is fixed in production.

---

## Step 7: Document (Optional but Recommended)

**Who:** Anyone

**What to do:**
- Add entry to KEDB (Known Error Database)
- Document:
  - Error message/symptoms
  - Root cause
  - Solution
  - How to prevent

**Output:**
- KEDB entry
- Future teams can find this faster

---

## Troubleshooting

**Can't reproduce the issue**
- Request more details from reporter
- Try different environments
- Consider race conditions or timing issues

**Root cause is unclear**
- Switch from ReAct to Tree of Thoughts
- Generate multiple hypotheses
- Test each systematically

**Fix causes regressions**
- Roll back
- Loop back to Step 3 with new information
- Adjust strategy

---

## Files Involved

- Issue tracker (GitHub issues, etc.)
- KEDB (Known Error Database)
- `tests/` - Reproduction & fix tests
- Source code
- `CHANGELOG.md` - Hotfix notes

## See Also

- **LOOP_CATALOG.md** - ReAct and Tree of Thoughts specifications
- **QUICK_START.md** - Quick reference for loops
- **docs/reference/glossary.md** - KEDB definition
