#!/bin/bash
# scripts/sandbox_run.sh - Execute commands in a sandboxed Podman container

set -e

WORKSPACE=$(pwd)
IMAGE="alpine:latest"
MEMORY_LIMIT="512m"
CPU_LIMIT="1"

# Check if podman is available
if ! command -v podman > /dev/null 2>&1; then
    echo "Error: podman is not installed. Sandboxing requires podman." >&2
    exit 1
fi

# Run the command in the sandbox
# --rm: Remove container after exit
# --cap-drop=ALL: Remove all capabilities
# --security-opt=no-new-privileges: Prevent privilege escalation
# --memory: Limit memory usage
# --cpus: Limit CPU usage
# -v: Mount workspace
# -w: Set working directory
# :Z flag ensures correct SELinux context if enabled
podman run --rm \
    --cap-drop=ALL \
    --security-opt=no-new-privileges \
    --memory="$MEMORY_LIMIT" \
    --cpus="$CPU_LIMIT" \
    -v "$WORKSPACE:/workspace:Z" \
    -w /workspace \
    "$IMAGE" \
    "$@"
