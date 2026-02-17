# Standard Operating Workflows

This guide provides step-by-step procedures for the most common tasks in the Springfield Protocol.

---

## 1. Discovery Diamond: From Idea to Brief

The Discovery Diamond ensures we solve the right problem before we start building.

1.  **DIVERGE: Orthogonal Investigation**
    -   **Product Agent (Marge)** generates options using **Empathy bias** (How does this impact the user?).
    -   **Planning Agent (Lisa)** adds **Structural bias** (How does this fit our architecture?).
    -   **Result:** A wide "Tree of Thoughts" of potential problem definitions and solution directions.
2.  **POOL:** Consolidate all biased options into a candidate list.
3.  **CONVERGE: Validate & Shortlist**
    -   **Planning Agent** filters by architectural fit and cost.
    -   **Product Agent** filters by user alignment.
    -   **Output:** A validated **Feature.md** (Brief) and necessary **ADRs**.

---

## 2. Delivery Diamond: From Brief to Release

---

## 2. Debugging an Issue

Standard procedure for resolving bugs.

1.  **Product Agent:** Triages issue → Enforces "Definition of Ready."
2.  **KEDB Check:** Search the Known Error Database for existing solutions.
3.  **Build Agent:** Uses **ReAct Loop** to reproduce error via failing test → Implements fix.
4.  **Quality Agent:** Verifies regression tests pass.
5.  **Release Agent:** Deploys hotfix → Captures learning in KEDB if it was a new error.

---

## 3. Designing Architecture

How to make and record technical decisions.

1.  **Planning Agent (Lisa):** Proposes approach based on existing patterns → Drafts ADR.
2.  **Quality Agent (Bart):** Adversarial review of ADR ("Poke holes" in the design).
3.  **Planning Agent (Lisa):** Refines ADR based on feedback.
4.  **Acceptance:** Mark ADR as `Accepted`.
5.  **Build Agent:** Implements feature using the ADR as a guardrail.

---

## 4. Releasing a Version

The ceremonial path to shipping code.

1.  **Verification:** Ensure all tasks in `PLAN.md` are marked `verified`.
2.  **Release Agent (Lovejoy):**
    -   Determine next version (SemVer).
    -   Gather changes from `TODO.md` and `FEEDBACK.md`.
    -   Generate `CHANGELOG.md` entry.
    -   Git tag and commit version bump.
    -   Deploy to production.
3.  **Learning:** Conduct a brief retrospective → Update `Feature.md` assumptions if any were broken.

---

## Workflow Summary Table

| Task | Primary Agent | Key Output |
| :--- | :--- | :--- |
| **New Feature** | Product / Planning | `Feature.md`, `PLAN.md` |
| **Bug Fix** | Build | Failing reproduction test |
| **Tech Decision** | Planning | ADR |
| **Shipping** | Release | git tag, `CHANGELOG.md` |
