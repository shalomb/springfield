# ADR-006: Adoption of the Atomic Commit Protocol (ACP)

**Date:** 2026-02-18
**Status:** Accepted

## Context
As the Springfield project grows and involves multiple autonomous agents (Ralph, Lisa, Bart, etc.) alongside human contributors, the clarity and integrity of the project history become paramount. 

Current commit practices follow Conventional Commits but lack a formal enforcement of "atomicity" that includes specifications, tests, and documentation within a single indivisible unit.

## Decision
We will formally adopt **The Atomic Commit Protocol (ACP) 1.0.0** as the standard for all contributions to the Springfield project.

### Key Requirements:
1.  **Indivisible Units:** Every functional commit **MUST** contain:
    *   One or more formal BDD specifications.
    *   A set of passing TDD tests (verifying a prior RED state).
    *   The minimal implementation required.
    *   Updated end-user documentation.
2.  **Commit Message Structure:**
    *   Conform to Conventional Commits 1.0.0.
    *   Header: Capitalized, imperative, no period, <= 50 chars.
    *   Body: Explains the "why", wrapped at 72 chars.
3.  **Strict Separation:** Refactoring and reformatting **MUST** be committed separately from functional changes.
4.  **Agent Adherence:** Agents (Ralph, Lisa) **MUST** be updated to strictly follow this protocol in their workflows.

## Consequences

### Positive
*   **Auditability:** Every change is self-documenting and self-verifying.
*   **Resilience:** project history becomes a "navigable database" of intent.
*   **Safety:** Forces agents to externalize intent (BDD) and reasoning (commit body).
*   **Reversibility:** Clean rollbacks with zero side effects on other features.

### Negative
*   **Overhead:** Requires more discipline and potentially more commits.
*   **Friction:** Agents may take slightly longer to complete tasks as they must ensure all ACP components are present.

## Implementation Plan
1.  Create `docs/standards/atomic-commit-protocol.md` containing the spec.
2.  Update `Justfile` recipes (e.g., `ralph`) to include ACP instructions in the system prompt.
3.  Update Agent definitions in `.pi/agents/`.
