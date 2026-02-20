# TODO: Epic-009 Fixes & Reliability

Address critical feedback from Bart regarding the Springfield Binary Orchestrator.

## 🛠️ Tasks

### 1. Robust Worktree Verification
- [ ] **Task:** Update `internal/orchestrator/worktree.go` to verify worktree validity.
- [ ] **Details:** 
    - In `EnsureWorktree`, if directory exists, run `git worktree list --porcelain` and verify the directory is a managed worktree.
    - If it's a stale directory (not a worktree), remove it and recreate.
    - Verify `.git` file presence as a secondary check.
- [ ] **Success Criteria:** Orchestrator correctly handles pre-existing non-worktree directories in the worktrees path.
- [ ] **ACP:** `feat(orchestrator): robust worktree verification`

### 2. Implement Circuit Breaker / Invoke Lisa on Bart Rejection
- [ ] **Task:** Update `internal/orchestrator/orchestrator.go` to invoke Lisa when Bart rejects an implementation.
- [ ] **Details:**
    - Change `StatusImplemented` -> `bart_fail_implementation` logic.
    - Instead of transitioning to `in_progress` and running `ralph`, transition to `blocked` (or a new `revising` state) and invoke `lisa`.
    - Lisa will then read `FEEDBACK.md` and update `TODO.md` before transitioning back to `ready`.
- [ ] **Success Criteria:** `internal/orchestrator/orchestrator_test.go` updated and passing with Lisa invocation on Bart failure.
- [ ] **ACP:** `fix(orchestrator): invoke lisa on bart rejection to avoid infinite loops`

### 3. Strict Handoff Deposit
- [ ] **Task:** Make `DepositHandoff` failures hard errors in `processEpic`.
- [ ] **Details:**
    - In `orchestrator.go`, check the error returned by `o.Worktree.DepositHandoff(id)`.
    - If error != nil, return the error instead of just logging a warning.
- [ ] **Success Criteria:** Orchestrator stops and reports error if `TODO-{id}.md` is missing when transitioning to `in_progress`.
- [ ] **ACP:** `fix(orchestrator): treat missing handoff as hard error`

### 4. Robust TD Output Parsing
- [ ] **Task:** Refactor `internal/orchestrator/td.go` to handle inconsistent JSON output from `td`.
- [ ] **Details:**
    - Ensure `GetEpic` and `QueryIDs` (or `QueryEpics`) use robust unmarshaling.
    - Check if `td show --json` always returns an object or sometimes an array.
- [ ] **Success Criteria:** Parsing logic handles both single object and array-of-one if applicable.
- [ ] **ACP:** `refactor(orchestrator): robust td json parsing`

### 5. Integration Test Fixes
- [ ] **Task:** Incorporate Bart's fixes into `tests/integration/orchestrator_runner_test.go`.
- [ ] **Details:**
    - Ensure `runner.Run` is called with the correct number of arguments (3).
- [ ] **Success Criteria:** `go test ./tests/integration/...` passes.
- [ ] **ACP:** `test(orchestrator): sync integration test with runner interface`
