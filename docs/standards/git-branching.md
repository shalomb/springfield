# Git Branching Standard (Trunk-Based Coordination)

## 1. Core Principles
- **Trunk-Based**: `main` is the single source of truth for code and the roadmap (`PLAN.md`).
- **Short-Lived Branches**: Implementation happens on feature branches that last no longer than one Epic.
- **Squash & Merge**: Keeps `main` history clean and atomic.
- **Shared State**: `PLAN.md` updates coordinate multiple agents. In protected environments, a dedicated `coordination` branch or "Planning PRs" are used to bypass `main` restrictions for the roadmap.
- **Concurrent Execution**: Use **Git Worktrees** to allow multiple agents to work on different epics/branches in parallel on the same host.
- **Merge Queues**: High-concurrency merges are managed via a merge queue to ensure integration tests pass against the "eventual main" state.

## 2. Branch Naming Conventions
- `main`: The trunk.
- `feat/epic-[ID]-[short-name]`: For new feature work (e.g., `feat/epic-001-git-standard`).
- `fix/epic-[ID]-[short-name]`: For bug fixes.
- `chore/[short-name]`: For non-functional updates (e.g., `chore/update-glossary`).

## 3. Workflow Lifecycle

### Phase 1: Planning (Lisa)
1. Lisa updates `PLAN.md` to define a new Epic and marks it `Ready`.
2. **Protected Environments**: Lisa commits to a `coordination` branch or opens a "Planning PR" for the roadmap update.
3. Once the Plan is "Ready", Lisa (or the harness) creates a new branch `feat/epic-ID-name` from `main`.
4. Lisa creates `TODO.md` on the feature branch.

### Phase 2: Implementation (Ralph)
1. Ralph monitors the coordination source (`main` or `coordination`).
2. If an Epic is `Ready`, Ralph spawns a **Git Worktree** for the specific feature branch.
3. Ralph executes the **Ralph Wiggum Loop** inside that isolated worktree.
4. Ralph periodically merges the coordination source into his worktree to stay in sync.

### Phase 3: Review & Verification (Bart)
1. Bart reviews the code on the feature branch.
2. Bart runs BDD scenarios to verify the Epic's acceptance criteria.
3. Bart updates `FEEDBACK.md` and marks the Epic as `Verified` in the branch's `PLAN.md`.

### Phase 4: Release (Lovejoy)
1. Lovejoy submits the feature branch to a **Merge Queue** (e.g., GitHub Merge Queue).
2. The queue validates the branch against the *resulting* state of `main`.
3. Upon success, the queue performs a **Squash & Merge** into `main`.
4. `main` now reflects the new code and shows the Epic as `Done` in `PLAN.md`.
5. The feature branch and its git worktree are deleted.

## 4. Conflict Resolution
- If two epics touch the same code, the second one to merge must rebase/merge `main`.
- `PLAN.md` conflicts are resolved by prioritizing the `main` branch's list of Epics while preserving the feature branch's status for its specific Epic.
