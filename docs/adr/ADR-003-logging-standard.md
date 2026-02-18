# ADR-003: Structured Logging Standard

**Date:** 2026-02-18
**Status:** Proposed
**Epic:** EPIC-003

## Problem
Currently, agent output is mostly unstructured text or captured via Tmux panes. This makes automated analysis, debugging, and audit trails difficult. We need a machine-readable, consistent logging format across all agents and the orchestrator.

## Decision
We will implement a **JSON-based structured logging standard**. All Springfield agents and scripts must emit logs in a consistent JSON format.

## Specification

### 1. Log Format (JSON)
Every log line must be a valid JSON object.

**Required Fields:**
- `timestamp`: ISO 8601 format (e.g., `2026-02-18T10:00:00Z`).
- `level`: Log level (e.g., `INFO`, `DEBUG`, `WARN`, `ERROR`, `FATAL`).
- `agent`: The name of the agent or component (e.g., `ralph-1`, `lisa`, `orchestrator`).
- `message`: A clear, human-readable description of the event.

**Optional/Contextual Fields:**
- `epic`: The ID of the current Epic (e.g., `EPIC-003`).
- `task`: The description or ID of the current Task.
- `error`: Detailed error message or stack trace (only for `ERROR` and `FATAL`).
- `data`: A nested object for event-specific metadata.

### 2. Log Levels
- `DEBUG`: Verbose technical details for troubleshooting.
- `INFO`: Significant operational events (e.g., Task started, Task complete).
- `WARN`: Non-blocking issues or unexpected behavior.
- `ERROR`: A task or command failed, but the agent/orchestrator can continue.
- `FATAL`: Unrecoverable failure that stops execution.

### 3. Log Storage
- **Per-Agent Logs:** `logs/[agent_name].log` (e.g., `logs/ralph-1.log`).
- **Orchestration Logs:** `logs/orchestrator.log` for session management and window lifecycle events.
- **Combined View:** `logs/all.log` (Future: a background process to merge/stream all logs for easier tailing).

## Rationale
- **Machine Readable:** JSON allows tools like `jq`, `grep`, and future observability agents to parse and filter logs easily.
- **Traceability:** Including the agent and epic ID allows us to reconstruct the flow of a feature across multiple parallel loops.
- **Pi Compatibility:** Simple file-based logging works natively in the `pi` environment without external service dependencies.

## Consequences
- **Code Change:** Existing scripts (like `tmux-orch.sh`) must be updated to use the new format.
- **Storage Management:** Log files will grow over time; we will eventually need rotation (EPIC-003 follow-up).
- **Overhead:** Small performance hit for JSON serialization, negligible for our use cases.

## Alternatives Considered
- **Plain Text (Syslog style):** Harder to parse nested metadata.
- **External Services (Loki/CloudWatch):** Over-engineered for local bootstrap; adds external dependencies.
