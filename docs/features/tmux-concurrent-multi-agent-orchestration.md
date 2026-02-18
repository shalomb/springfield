# Feature.md - Tmux Concurrent Multi-Agent Orchestration

## Problem
Running multiple agents concurrently can be chaotic and difficult to manage. There's no centralized way to start, stop, monitor, or manage the lifecycle of concurrent agent sessions.

## Requirements
- Use `tmux` to manage multiple terminal sessions for concurrent agent execution.
- Create scripts to easily launch and attach to agent sessions.
- Provide a clean interface for orchestrating multiple agents working on different tasks.

## Acceptance Criteria
- [ ] Script to launch a new tmux session with named windows for different agents.
- [ ] Integration with `just` commands (e.g., `just tmux-start`, `just tmux-attach`).
- [ ] Ability to run agents in detached mode and view logs later.

## Constraints & Unknowns
- **Constraint:** Must work within the existing pi environment which likely has tmux available.
- **Unknown:** Handling of agent failures within tmux sessions.

## Options Considered
- [ ] Tmux: Lightweight, standard, scriptable.
- [ ] Screen: Older, less feature-rich.
- [ ] Docker Compose: Good for containerized orchestration, but maybe overkill for just script execution.

## Scope
✅ Tmux configuration and scripts
✅ `just` command integration
❌ Full-blown agent monitoring UI (separate feature)

## Success Criteria
- Ability to run 3+ agents concurrently without manual window management.
- Easy access to agent logs and output.
