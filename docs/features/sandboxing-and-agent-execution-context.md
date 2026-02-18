# Feature.md - Sandboxing and Agent Execution Context

## Problem
Agents currently run with full permissions in the host environment. This poses security risks (e.g., unintended file deletion, command execution) and creates a fragile execution environment.

## Requirements
- Isolate agent execution environments from the host system.
- Provide a clean, consistent execution context for each agent task.
- Ensure agents cannot access or modify unauthorized files or resources.

## Acceptance Criteria
- [ ] Define a sandbox environment (e.g., Docker container, chroot, or restricted user).
- [ ] Script to setup and teardown sandbox environments.
- [ ] Mechanism to mount specific directories or files into the sandbox.

## Constraints & Unknowns
- **Constraint:** Must work within the existing pi environment which may limit container capabilities.
- **Unknown:** Performance overhead of sandboxing.

## Options Considered
- [ ] Docker: Standard containerization, strong isolation, but might require privileges.
- [ ] chroot/jail: Lighter weight, but less isolation.
- [ ] User-level sandboxing (e.g., restricted shell): Minimal overhead, but weakest isolation.

## Scope
✅ Sandbox definition (Docker or similar)
✅ Context setup script
✅ Mount/volume management for shared workspace
❌ Network isolation (future)

## Success Criteria
- Agents run in isolated environments without full host access.
- Execution environment is reproducible and clean.
