#!/bin/bash
set -euo pipefail
# test_final_verification.sh - Verify all files and links exist

FILES=(
    "AGENTS.md"
    "README.md"
    "QUICK_START.md"
    "docs/standards/atomic-commit-protocol.md"
    "docs/standards/coding-conventions.md"
    "docs/features/README.md"
    "docs/adr/ADR-000-compliance-and-safety.md"
)

for f in "${FILES[@]}"; do
    if [ ! -f "$f" ]; then
        echo "FAIL: $f does not exist"
        exit 1
    fi
done

echo "PASS: All required files exist"
exit 0
