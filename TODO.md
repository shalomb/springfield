# TODO: EPIC-007 - Autonomous Development Loop ("just do")

## Context
Implementing the end-to-end autonomous loop to allow Springfield to self-evolve with minimal human intervention. This includes automated planning, execution, and verification.

## Tasks

### 🔧 Phase 0: Stabilization & Debt (Reflect & Learn)
- [x] **Task 0.1: Robustness for scripts (set -euo pipefail)**
  - Success Criteria: All `.sh` scripts in `scripts/` have `set -euo pipefail`.
  - ACP: `fix(scripts): add robust error handling to bash scripts`
- [x] **Task 0.2: Justfile DRYness for Agents**
  - Success Criteria: Repetitive agent calls in Justfile are refactored into a shared variable or template.
  - ACP: `refactor(build): DRY up agent targets in Justfile`
- [x] **Task 0.3: Complete .gitignore**
  - Success Criteria: Standard exclusions for Python and Go artifacts are added to avoid Ralph loops.
  - ACP: `fix(infra): update .gitignore for Python and Go environments`

### 🏗️ Phase 1: Planning Automation (Lisa)
- [ ] **Task 1.1: Implement `just plan` command**
  - Success Criteria: Command that runs Lisa specifically to update `TODO.md` from `PLAN.md`.
  - ACP: `feat(protocol): add just plan command for Lisa`
- [ ] **Task 1.2: Lisa branch management logic**
  - Success Criteria: Lisa can detect current branch and create `feat/` branches if on `main`.
  - ACP: `feat(lisa): implement automated branch management logic`

### 🔨 Phase 2: Execution & Loop (Ralph & Orchestrator)
- [ ] **Task 2.1: Implement `just do` entry point**
  - Success Criteria: A `just do` command in `Justfile` that chains Lisa, then Ralph.
  - ACP: `feat(protocol): implement initial just do orchestrator`
- [ ] **Task 2.2: Continuous Ralph Execution**
  - Success Criteria: Ralph continues working as long as tasks exist in `TODO.md`. (Partially implemented, needs verification).
  - ACP: `feat(ralph): ensure Ralph loop handles task exhaustion correctly`

### 🔍 Phase 3: Quality Gate (Bart & Herb)
- [ ] **Task 3.1: Automated Verification loop**
  - Success Criteria: `just do` includes Bart/Herb and fails back to Lisa if `FEEDBACK.md` has blockers.
  - ACP: `feat(protocol): implement feedback loop between Bart and Lisa`

---
**Status:** 🏗️ In Progress
**Current Task:** Task 0.1
