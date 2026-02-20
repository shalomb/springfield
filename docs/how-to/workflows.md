# Standard Workflows

This is how we get things done.

---

## 1. The New Feature Workflow

**Goal:** Turn an idea into shipped code.
**Agents:** Marge, Lisa, Ralph, Bart, Lovejoy.

### Phase 1: Discovery (Marge & Lisa)
1.  **Start:** `just start-feature <name>`
2.  **Define:** Create `Feature.md`. What is the problem? Who cares?
3.  **Plan:** Lisa reads `Feature.md` and populates `PLAN.md` with Epics.
4.  **Architect:** If it's complex, Lisa writes an ADR to decide *how* we build it.

### Phase 2: Delivery (Ralph & Bart)
1.  **Tasking:** Lisa breaks Epics into atomic tasks in `TODO.md`.
2.  **The Loop:**
    *   **Ralph** picks a task, writes a failing test, makes it pass, and commits.
    *   **Bart** performs comprehensive quality verification:
        - Static analysis (code review, ACP adherence, best practices)
        - Dynamic verification (test execution, coverage, BDD validation)
        - Adversarial testing (edge cases, security, performance)
3.  **Repeat:** Keep looping until `TODO.md` is empty.

### Phase 3: Release (Lovejoy)
1.  **Verify:** Ensure `FEEDBACK.md` is empty of blockers.
2.  **Ship:** Lovejoy merges the branch and updates the `CHANGELOG.md`.

---

## 2. The Bug Fix Workflow

**Goal:** Fix a bug without creating two new ones.
**Agents:** Ralph, Bart.

1.  **Reproduce:** Ralph writes a test case that replicates the bug. **If you can't reproduce it with a test, you can't fix it.**
2.  **Fix:** Ralph writes the code to make the test pass.
3.  **Verify:** Bart runs the regression suite to ensure nothing else broke.
4.  **Ship:** Merge and deploy.

---

## 3. The Architecture Workflow

**Goal:** Make a hard technical decision.
**Agents:** Lisa, Bart, Ralph.

1.  **Propose:** Lisa drafts an ADR (Architectural Decision Record) in `docs/adr/`.
2.  **Critique:** Bart reviews the ADR for logical soundness, adherence to principles, and potential edge cases.
3.  **Validate:** Ralph implements a prototype to verify the decision is viable.
4.  **Refine:** Lisa updates the ADR based on feedback.
5.  **Decide:** Mark the ADR as `Accepted` or `Rejected`.
6.  **Comply:** Ralph must now follow this rule forever.

---

## 4. The Release Workflow

**Goal:** Ship code to production.
**Agents:** Lovejoy.

1.  **Check:** Are all tasks in `PLAN.md` done?
2.  **Clean:** Is `TODO.md` empty?
3.  **Review:** Are there any critical issues in `FEEDBACK.md`?
4.  **Ceremony:** Lovejoy bumps the version, updates `CHANGELOG.md`, and creates a git tag.

---

## Workflow Cheat Sheet

| I want to... | Run this... |
| :--- | :--- |
| **Start a feature** | `just start-feature <name>` |
| **Plan the work** | `just lisa` |
| **Build the code** | `just ralph` |
| **Verify quality** | `just bart` |
| **Ship the code** | `just lovejoy` |
