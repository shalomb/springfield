# TODO: Complete Autonomous Development Loop (EPIC-007 Recovery)

This plan addresses the critical failures identified in `FEEDBACK.md` by Bart Simpson. We must establish a working Go environment and fix the orchestrator recipes.

## 🏁 Goal
Stabilize the Springfield infrastructure so the autonomous loop (Lisa -> Ralph -> Bart -> Herb) can execute and verify work successfully.

## 🛠 Tasks

### Phase 1: Go Foundation (ACP-1)
- [x] **Task 1: Initialize Go Module**
  - **Action:** Run `go mod init github.com/shalomb/springfield` and `go mod tidy`.
  - **Success Criteria:** `go.mod` and `go.sum` exist.
- [x] **Task 2: Create Minimal Entrypoint**
  - **Action:** Create `cmd/springfield/main.go` with a basic CLI structure.
  - **Success Criteria:** `just build` succeeds and `./bin/springfield --help` runs.

### Phase 2: Justfile Harmonization (ACP-2)
- [x] **Task 3: Add Missing Lifecycle Recipes**
  - **Action:** Implement `start-feature` and `start-fix` in `Justfile`.
  - **Success Criteria:** `just list` shows the new recipes.
- [x] **Task 4: Fix Test Runner Mismatch & Logger Stability**
  - **Action:** Update `test-integration` to correctly handle existing Python/BDD tests. Implement basic file locking or sequential logging to prevent race conditions in `scripts/logger.py`.
  - **Success Criteria:** `just test` passes Phase 1 and Phase 2. Logs remain valid JSON under concurrent load.

### Phase 3: Loop Verification (ACP-3)
- [ ] **Task 5: Verify Ralph's Loop with a dummy task**
  - **Action:** Create a small TODO task and run `just ralph`.
  - **Success Criteria:** Ralph completes the task, commits, and removes the TODO.
- [ ] **Task 6: Final Integration Check**
  - **Action:** Run `just do` for a minor documentation fix.
  - **Success Criteria:** Full loop executes without error and updates `FEEDBACK.md` to ✅.

### Phase 4: Feedback Loop Intelligence (ACP-4)
- [ ] **Task 7: Implement Lisa's Triage Logic**
  - **Action:** Update `just lisa` prompt to parse `FEEDBACK.md` and make Go/No-Go decisions.
  - **Success Criteria:** Lisa can read a mock `FEEDBACK.md` and either create a `TODO.md` task or update `PLAN.md`.
- [ ] **Task 8: Implement Loop Limiter (Circuit Breaker)**
  - **Action:** Add a counter (e.g., `git rev-list --count HEAD` or a temporary file) to `just do` to prevent infinite loops. Max retries = 2.
  - **Success Criteria:** `just do` aborts if `FEEDBACK.md` remains non-empty after 2 cycles.

## ⚖️ Standards Check
- [x] ADR-000 (Building Blocks) - N/A for this infrastructure fix.
- [x] ADR-001 (Git Branching) - Work performed on `feat/` branch.
- [x] ADR-005 (Atomic Commit Protocol) - Tasks mapped to ACP units.
