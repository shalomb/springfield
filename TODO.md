# TODO - EPIC-004: Agent Sandboxing

> **Epic:** EPIC-004: Agent Sandboxing
> **Context:** Isolate agent execution environments from the host system to ensure security and stability.
> **Status:** 🔍 Discovery

## Tasks

- [x] **Task 1: Research `pi` environment capabilities** @Ralph
    - [x] Check if `docker` command is available and functional.
    - [x] Check if `podman` or `nsenter` are available.
    - [x] Determine if current user has permissions to run containers.
    - [x] Document findings in `docs/research/sandboxing-capabilities.md`.
    - **Acceptance Criteria:** A clear report on what isolation technologies are available in the current environment.

- [ ] **Task 2: Draft ADR-004: Agent Sandboxing Strategy** @Lisa
    - [ ] Review research from Task 1.
    - [ ] Propose isolation strategy (e.g., Docker, restricted user, etc.).
    - [ ] Document resource constraints and workspace mounting strategy.
    - **Acceptance Criteria:** `docs/adr/ADR-004-agent-sandboxing.md` exists in "Proposed" state.

- [x] **Task 3: Create BDD Scenarios** @Ralph
    - [x] Create `features/sandboxing.feature`.
    - [x] Define scenarios for:
        - Successful execution in sandbox.
        - Prevention of host file access (outside workspace).
        - Preservation of workspace state.
    - **Acceptance Criteria:** `features/sandboxing.feature` exists and reflects the requirements.

- [x] **Task 4: Prototype isolation script** @Ralph
    - [x] Based on ADR-004, create a minimal `scripts/sandbox.sh`.
    - [x] Script should be able to run a simple command (e.g., `ls`) in the sandbox.
    - [x] Attempt to verify isolation (e.g., try to touch a file in `/root`).
    - **Acceptance Criteria:** Prototype script successfully demonstrates basic isolation.

## Feedback Iteration (EPIC-004)

- [ ] **Task 5: Harden Sandbox Isolation** @Ralph
    - [ ] Pin container image to `docker.io/library/alpine:3.19` (avoid `latest`).
    - [ ] Disable network access (`--network none`) by default.
    - [ ] Mount `.git` and `scripts` directories as **read-only** to prevent accidental damage.
    - [ ] Improve command injection safety in `scripts/sandbox.sh` (use arrays/exec if possible).
    - **Acceptance Criteria:** `sandbox.sh` runs with `--network none` and read-only system mounts.

- [ ] **Task 6: Verify Git Safety (BDD)** @Ralph
    - [ ] Add scenario to `features/sandboxing.feature`: "Prevent modification of .git directory".
    - [ ] Verify that attempting to `rm -rf .git` or `touch .git/config` from within the sandbox fails.
    - **Acceptance Criteria:** Test suite passes and confirms `.git` protection.

## Notes & Blockers
- **Resolved Blocker:** The `pi` harness supports `podman` for container execution. `docker` is not available.
- **Reference:** See `docs/features/sandboxing-and-agent-execution-context.md` for initial requirements.
