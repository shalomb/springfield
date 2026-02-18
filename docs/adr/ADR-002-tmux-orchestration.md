# ADR-002: Tmux Agent Orchestration Strategy

**Date:** 2026-02-17
**Status:** Proposed
**Epic:** EPIC-002

## Problem
Running multiple specialized agents (Marge, Lisa, Ralph, Bart, Lovejoy) concurrently requires a way to manage their lifecycles, view their progress, and access their logs without cluttering the developer's terminal.

## Decision
We will use **Tmux** as the orchestration layer for concurrent agent execution.

## Proposed Architecture

### 1. Smart Session Adoption & Window Structure
- **Session Reuse:** The orchestrator will detect if it's already running inside an active Tmux session (via the `$TMUX` environment variable).
    - **In-Session:** If already in Tmux, it will simply create new windows/panes within the *current* session.
    - **Out-of-Session:** If not in Tmux, it will create a new session (named after the current project directory) and attach to it.
- **Project Isolation:** By basing session/window context on the current directory, we allow multiple Springfield projects to run in parallel in separate sessions or window groups.
- **Window Naming:** Each window is titled with the agent's name (e.g., `marge`, `lisa`).
- **Ralph Suffix:** As multiple Ralph agents can run concurrently for different epics, they will use a numbered suffix (e.g., `ralph-1`, `ralph-2`).

### 2. Window Layout
Each agent window will be split into two panes:
- **Top Pane (80%):** Active agent execution/TUI.
- **Bottom Pane (20%):** Continuous tail of the agent-specific log file (`logs/[agent-name].log`).

### 3. Orchestration Script (`scripts/tmux-orch.sh`)
- **Start:** Checks if session `springfield` exists. If not, creates it and initializes windows/panes.
- **Stop:** Kills the `springfield` session.
- **Attach:** Attaches the user to the existing session.
- **Spawn:** (Future) Dynamically adds a new window for a specific worktree/epic.

### 4. Integration
- `just flow`: The primary entry point. It runs the orchestrator in the background and attaches the UI.

## Rationale
- **Isolation:** Each agent runs in its own shell environment within a tmux window.
- **Persistence:** Agents can continue working if the user disconnects from the terminal.
- **Visibility:** Easy switching between agent contexts using standard tmux keys (`Ctrl-b n/p`).
- **Low Overhead:** Tmux is lightweight and available in most CLI environments (including Pi).

## Consequences
- **Tooling Dependency:** Requires `tmux` to be installed on the host.
- **Context Management:** The orchestrator must ensure agents are aware of their specific `TODO.md` and worktree.
- **Log Volume:** Agent logs will grow; requires a rotation strategy (EPIC-003).

## Alternatives Considered
- **Screen:** Less powerful window/pane management than Tmux.
- **Docker Compose:** Good for containerized agents, but higher overhead for simple script execution; harder to "attach" to multiple TUIs interactively.
- **Custom TUI:** Over-engineering; Tmux provides the "terminal multiplexing" for free.
