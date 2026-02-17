# Repository Protection Rules (Simulated)

To ensure the integrity of the `main` branch and the enforcement of the Springfield Protocol, the following protection rules are established. In a live GitHub/GitLab environment, these would be configured in the repository settings.

## 1. Branch: `main`

### Protections
- **Require a pull request before merging**: All changes to `main` must be made via a feature branch and a PR.
- **Require approvals**:
    - Minimum of **1 approval** required.
    - Approval must come from the **Quality Agent (Bart)** or a designated human reviewer.
- **Dismiss stale pull request approvals when new commits are pushed**: Ensures that the final code state is what was actually reviewed.
- **Require status checks to pass before merging**:
    - `test-structure`: Validates repository layout.
    - `test-lint`: Static analysis and formatting.
    - `test-fast`: Unit and integration tests.
    - `test-bdd`: Behavior-driven scenarios (Acceptance Criteria).
    - `coverage`: Must be **95% or higher** (verified by Bart).
- **Require conversation resolution before merging**: All comments in `FEEDBACK.md` or on the PR must be resolved.
- **Lock branch**: No direct pushes allowed.
- **Restrict deletions**: `main` cannot be deleted.

### Merge Strategy
- **Allow squash merging**: **Mandatory**. Keeps the trunk history clean and atomic.
- **Restrict to Squash & Merge**: Disable standard Merge Commits and Rebase Merging to maintain a flat history on `main`.

## 2. Branch: `coordination`

### Protections
- **Allow direct pushes from Planning Agent (Lisa)**: To facilitate rapid roadmap updates.
- **Automated Sync**: Any merge to `main` must trigger a downstream sync (merge/rebase) into `coordination`.

## 3. Enforcement (Simulated)

The `just gate` command will be developed (in a future Epic) to simulate these checks locally for agents before they attempt to merge.

| Check | Tool | Requirement |
| :--- | :--- | :--- |
| Approval | `Bart` | `Verified` status in `PLAN.md` |
| BDD | `pytest-bdd` | 100% Pass |
| Coverage | `pytest-cov` | >= 95% |
| Lint | `ruff` | 0 errors |
