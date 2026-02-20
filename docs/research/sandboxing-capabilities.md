# Research: Agent Sandboxing Capabilities

**Date:** 2026-02-18
**Author:** @Ralph
**Context:** EPIC-004: Agent Sandboxing

## Goal
Investigate the available tools and capabilities in the `pi` environment to support agent sandboxing and execution isolation.

## Findings

### 1. Docker
- **Availability:** ❌ Not installed.
- **Command:** `docker: command not found`
- **Socket:** `/var/run/docker.sock` does not exist.
- **Permissions:** User `unop` is in the `docker` group, but the command is missing.

### 2. Podman
- **Availability:** ✅ Installed.
- **Version:** 5.4.2
- **Functionality:** 
  - Capable of running rootless containers.
  - Successfully pulled and ran `docker.io/library/alpine:latest`.
  - Requires fully qualified image names (e.g., `docker.io/library/alpine` instead of `alpine`) due to missing registry aliases configuration.
- **Permissions:** Works for current user without `sudo`.

### 3. Nsenter
- **Availability:** ✅ Installed.
- **Version:** From `util-linux 2.41`.
- **Functionality:** Can enter namespaces, useful for debugging or advanced isolation techniques if containers are insufficient.

### 4. User Permissions
- **User:** `unop` (uid 1001)
- **Groups:** `docker`, `sudo`, `ollama`, `parallels`, etc.
- **Rootless Containers:** Validated with `podman`.

## Conclusion
The environment supports containerized isolation via **Podman**. Docker is not available.

## Recommendation
- Use **Podman** as the primary engine for agent sandboxing.
- Ensure all image references are fully qualified (e.g., prefix with `docker.io/library/`).
- If Docker compatibility is required for scripts, consider aliasing `docker` to `podman` or updating scripts to use `podman` directly.
