#!/usr/bin/env bash
# scripts/tmux-orch.sh - Springfield Agent Orchestrator

set -euo pipefail

SESSION_NAME=$(basename "$PWD" | tr '.' '-')
AGENTS=("marge" "lisa" "ralph-1" "bart" "lovejoy")

create_agent_window() {
    local agent=$1
    local session=$2
    local window_name=$agent
    
    # Create window if it doesn't exist
    if ! tmux list-windows -t "$session" -F "#{window_name}" | grep -q "^$window_name$"; then
        tmux new-window -t "$session" -n "$window_name"
    fi

    # Create log file
    touch "logs/${agent}.log"

    # Split window: Top 80% for execution, Bottom 20% for logs
    # We clear the window first to ensure a clean state
    tmux send-keys -t "${session}:${window_name}.0" "clear" C-m
    
    # Check if we already have panes, if not split
    local pane_count=$(tmux list-panes -t "${session}:${window_name}" | wc -l)
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
    
    if [ -n "${TMUX:-}" ]; then
        # Already inside a tmux session
        CURRENT_SESSION=$(tmux display-message -p '#S')
        echo "Adopting existing tmux session: $CURRENT_SESSION"
        for agent in "${AGENTS[@]}"; do
            create_agent_window "$agent" "$CURRENT_SESSION"
        done
    else
        # Not in tmux, create a new session
        if tmux has-session -t "$SESSION_NAME" 2>/dev/null; then
            echo "Session $SESSION_NAME already exists. Attaching..."
            tmux attach-session -t "$SESSION_NAME"
        else
            echo "Creating new tmux session: $SESSION_NAME"
            tmux new-session -d -s "$SESSION_NAME" -n "marge"
            # marge is window 0, rename it if necessary or just start from 1
            # For simplicity, we create the others
            for agent in "${AGENTS[@]}"; do
                create_agent_window "$agent" "$SESSION_NAME"
            done
            # Remove the default window if it was created and we didn't use it
            # Actually marge is window 1 in our create_agent_window logic usually
            tmux attach-session -t "$SESSION_NAME"
        fi
    fi
}

stop_flow() {
    if [ -n "${TMUX:-}" ]; then
        CURRENT_SESSION=$(tmux display-message -p '#S')
        echo "Cleaning up springfield windows in $CURRENT_SESSION..."
        for agent in "${AGENTS[@]}"; do
             tmux kill-window -t "${CURRENT_SESSION}:${agent}" 2>/dev/null || true
        done
    else
        echo "Killing session $SESSION_NAME..."
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
