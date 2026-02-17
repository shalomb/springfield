# Git Branching Standard (Trunk-Based Coordination)

## 1. Core Principles
- **Trunk-Based**: `main` is the single source of truth for code and the roadmap (`PLAN.md`).
- **Short-Lived Branches**: Implementation happens on feature branches that last no longer than one Epic.
- **Squash & Merge**: Keeps `main` history clean and atomic.
- **Shared State**: `PLAN.md` updates on `main` coordinate multiple agents.

## 2. Branch Naming Conventions
- `main`: The trunk.
- `feat/epic-[ID]-[short-name]`: For new feature work (e.g., `feat/epic-001-git-standard`).
- `fix/epic-[ID]-[short-name]`: For bug fixes.
- `chore/[short-name]`: For non-functional updates (e.g., `chore/update-glossary`).

## 3. Workflow Lifecycle

### Phase 1: Planning (Lisa)
1. Lisa updates `PLAN.md` on `main` to define a new Epic.
2. Lisa marks the Epic as `Ready`.
3. Lisa (or the harness) creates a new branch `feat/epic-ID-name` from `main`.
4. Lisa creates `TODO.md` on the new branch.

### Phase 2: Implementation (Ralph)
1. Ralph monitors `main`. If an Epic is `Ready`, he switches to the feature branch.
2. Ralph executes the **Ralph Wiggum Loop** on the feature branch.
3. Ralph periodically merges `main` into the feature branch to stay in sync with the roadmap.

### Phase 3: Review & Verification (Bart)
1. Bart reviews the code on the feature branch.
2. Bart runs BDD scenarios to verify the Epic's acceptance criteria.
3. Bart updates `FEEDBACK.md` and marks the Epic as `Verified` in the branch's `PLAN.md`.

### Phase 4: Release (Lovejoy)
1. Lovejoy performs a **Squash & Merge** of the feature branch into `main`.
2. `main` now reflects the new code and shows the Epic as `Done` in `PLAN.md`.
3. The feature branch is deleted.

## 4. Conflict Resolution
- If two epics touch the same code, the second one to merge must rebase/merge `main`.
- `PLAN.md` conflicts are resolved by prioritizing the `main` branch's list of Epics while preserving the feature branch's status for its specific Epic.
