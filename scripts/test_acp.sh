#!/bin/bash
# test_acp.sh - Verify Atomic Commit Protocol document exists and has required sections

FILE="docs/standards/atomic-commit-protocol.md"

if [ ! -f "$FILE" ]; then
    echo "FAIL: $FILE does not exist"
    exit 1
fi

grep -q "# Atomic Commit Protocol" "$FILE" || { echo "FAIL: Missing Title"; exit 1; }
grep -q "## Commit Message Standard" "$FILE" || { echo "FAIL: Missing Commit Message Standard"; exit 1; }
grep -q "## Commit Scope" "$FILE" || { echo "FAIL: Missing Commit Scope"; exit 1; }
grep -q "## Examples" "$FILE" || { echo "FAIL: Missing Examples"; exit 1; }

echo "PASS: $FILE meets standards"
exit 0
