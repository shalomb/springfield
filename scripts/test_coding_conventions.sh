#!/bin/bash
set -euo pipefail
# test_coding_conventions.sh - Verify Coding Conventions document structure

FILE="docs/standards/coding-conventions.md"

if [ ! -f "$FILE" ]; then
    echo "FAIL: $FILE does not exist"
    exit 1
fi

grep -q "# Coding Conventions" "$FILE" || { echo "FAIL: Missing Title"; exit 1; }
grep -q "## Go Standards" "$FILE" || { echo "FAIL: Missing Go Standards"; exit 1; }
grep -q "## Python Standards" "$FILE" || { echo "FAIL: Missing Python Standards"; exit 1; }
grep -q "Atomic Commit Protocol" "$FILE" || { echo "FAIL: Missing Reference to ACP"; exit 1; }

echo "PASS: $FILE meets standards"
exit 0
