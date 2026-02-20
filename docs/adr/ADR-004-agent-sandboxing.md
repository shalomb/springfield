# ADR-004: Agent Sandboxing Strategy

**Date:** 2026-02-18
**Status:** Proposed

## Context
Agents in the Springfield protocol (Ralph, Lisa, Bart, etc.) currently run with the full privileges of the host user. This poses significant risks:
1.  **Accidental Damage:** An agent might delete critical project files (e.g., `.git`, `scripts`).
2.  **Security:** Agents might access sensitive host files (e.g., SSH keys) or exfiltrate data.
3.  **Reproducibility:** Dependencies might be missing on the host environment.

We need a mechanism to isolate agent execution while preserving their ability to modify the project workspace.

## Decision
We will use **Podman** containers to sandbox agent execution.

### 1. Technology: Podman
- **Reasoning:** Podman is available in the `pi` environment and supports rootless containers, avoiding the need for `sudo` or a Docker daemon.
- **Image:** `docker.io/library/alpine:3.19` (pinned version for stability).

### 2. Isolation Configuration
- **Network:** Enabled by default (host/slirp4netns) to support LLMs, git, and web search.
- **Workspace:** The current working directory is mounted to `/workspace` inside the container.
- **Protection:**
    - `.git` directory is mounted as **read-only** to prevent history corruption.
    - `scripts` directory is mounted as **read-only** to prevent modification of the sandbox tool itself.

### 3. Execution Interface
- A wrapper script `scripts/sandbox.sh` handles the container lifecycle.
- Usage: `scripts/sandbox.sh "command to run"`
- The script automatically detects `.git` and `scripts` directories and applies protection.

## Consequences
### Positive
- **Safety:** Agents cannot accidentally destroy git history or modify tooling.
- **Consistency:** Execution happens in a clean Alpine environment.
- **Connectivity:** Agents have full network access for LLMs, git operations, and web search.

### Negative
- **Security:** Outbound network access is unrestricted (by design).
- **Performance:** Slight overhead for container startup (milliseconds to seconds).
- **Tooling:** Agents are limited to tools available in Alpine or those we install in the image.
- **Persistence:** Changes outside `/workspace` are lost after execution.

## Future Improvements
- Create a custom image with more developer tools (git, python, node, etc.) pre-installed.
