#!/usr/bin/env bash
# scripts/tmux-orch.sh - Springfield Agent Orchestrator

set -euo pipefail

SESSION_NAME=$(basename "$PWD" | tr '.' '-')
AGENTS=("marge" "lisa" "ralph-1" "bart" "lovejoy")

log_event() {
    local message=$1
    local level=${2:-INFO}
    local agent=${3:-orchestrator}
    python3 scripts/logger.py "$message" --level "$level" --agent "$agent"
}

create_agent_window() {
    local agent=$1
    local session=$2
    local window_name=$agent
    
    # Create window if it doesn't exist
    if ! tmux list-windows -t "$session" -F "#{window_name}" | grep -q "^$window_name$"; then
        tmux new-window -t "$session" -n "$window_name"
        log_event "Created window for agent: $agent" "DEBUG"
    fi

    # Create log file
    touch "logs/${agent}.log"

    # Split window: Top 80% for execution, Bottom 20% for logs
    tmux send-keys -t "${session}:${window_name}.0" "clear" C-m
    
    # Check if we already have panes, if not split
    local pane_count
    pane_count=$(tmux list-panes -t "${session}:${window_name}" | wc -l)
    if [ "$pane_count" -lt 2 ]; then
        tmux split-window -v -p 20 -t "${session}:${window_name}"
    fi

    # Start log tail in bottom pane
    tmux send-keys -t "${session}:${window_name}.1" "tail -f logs/${agent}.log" C-m
    
    # Focus top pane
    tmux select-pane -t "${session}:${window_name}.0"
}

start_flow() {
    mkdir -p logs
    log_event "Starting orchestration flow" "INFO"
    
    if [ -n "${TMUX:-}" ]; then
        # Already inside a tmux session
        CURRENT_SESSION=$(tmux display-message -p '#S')
        log_event "Adopting existing tmux session: $CURRENT_SESSION" "INFO"
        for agent in "${AGENTS[@]}"; do
            create_agent_window "$agent" "$CURRENT_SESSION"
        done
    else
        # Not in tmux, create a new session
        if tmux has-session -t "$SESSION_NAME" 2>/dev/null; then
            log_event "Session $SESSION_NAME already exists. Attaching." "INFO"
            tmux attach-session -t "$SESSION_NAME"
        else
            log_event "Creating new tmux session: $SESSION_NAME" "INFO"
            tmux new-session -d -s "$SESSION_NAME" -n "marge"
            for agent in "${AGENTS[@]}"; do
                create_agent_window "$agent" "$SESSION_NAME"
            done
            tmux attach-session -t "$SESSION_NAME"
        fi
    fi
}

stop_flow() {
    log_event "Stopping orchestration flow" "INFO"
    if [ -n "${TMUX:-}" ]; then
        CURRENT_SESSION=$(tmux display-message -p '#S')
        for agent in "${AGENTS[@]}"; do
             tmux kill-window -t "${CURRENT_SESSION}:${agent}" 2>/dev/null || true
        done
    else
        tmux kill-session -t "$SESSION_NAME" 2>/dev/null || true
    fi
}

# Simple CLI
COMMAND=${1:-start}

case $COMMAND in
    start)
        start_flow
        ;;
    stop)
        stop_flow
        ;;
    attach)
        if [ -z "${TMUX:-}" ]; then
            tmux attach-session -t "$SESSION_NAME"
        else
            echo "Already inside a tmux session."
        fi
        ;;
    *)
        echo "Usage: $0 {start|stop|attach}"
        exit 1
        ;;
esac
