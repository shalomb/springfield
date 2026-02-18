# Springfield Protocol Justfile

# Launch the agent orchestration mesh
flow:
    @mkdir -p logs
    @bash scripts/tmux-orch.sh

# Attach to the current orchestration session
attach:
    @bash scripts/tmux-orch.sh attach

# Clean up agents and tmux windows
stop:
    @bash scripts/tmux-orch.sh stop
