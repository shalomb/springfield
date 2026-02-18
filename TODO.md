# TODO - EPIC-004: Agent Sandboxing

> **Epic:** EPIC-004: Agent Sandboxing
> **Context:** Isolate agent execution environments from the host system to ensure security and stability.
> **Status:** 🔍 Discovery

## Tasks

- [ ] **Task 1: Research `pi` environment capabilities** @Ralph
    - [ ] Check if `docker` command is available and functional.
    - [ ] Check if `podman` or `nsenter` are available.
    - [ ] Determine if current user has permissions to run containers.
    - [ ] Document findings in `docs/research/sandboxing-capabilities.md`.
    - **Acceptance Criteria:** A clear report on what isolation technologies are available in the current environment.

- [ ] **Task 2: Draft ADR-004: Agent Sandboxing Strategy** @Lisa
    - [ ] Review research from Task 1.
    - [ ] Propose isolation strategy (e.g., Docker, restricted user, etc.).
    - [ ] Document resource constraints and workspace mounting strategy.
    - **Acceptance Criteria:** `docs/adr/ADR-004-agent-sandboxing.md` exists in "Proposed" state.

- [ ] **Task 3: Create BDD Scenarios** @Ralph
    - [ ] Create `features/sandboxing.feature`.
    - [ ] Define scenarios for:
        - Successful execution in sandbox.
        - Prevention of host file access (outside workspace).
        - Preservation of workspace state.
    - **Acceptance Criteria:** `features/sandboxing.feature` exists and reflects the requirements.

- [ ] **Task 4: Prototype isolation script** @Ralph
    - [ ] Based on ADR-004, create a minimal `scripts/sandbox.sh`.
    - [ ] Script should be able to run a simple command (e.g., `ls`) in the sandbox.
    - [ ] Attempt to verify isolation (e.g., try to touch a file in `/root`).
    - **Acceptance Criteria:** Prototype script successfully demonstrates basic isolation.

## Notes & Blockers
- **Blocker:** We don't know yet if the `pi` harness allows running Docker or other containerization tools. Task 1 is critical.
- **Reference:** See `docs/features/sandboxing-and-agent-execution-context.md` for initial requirements.
