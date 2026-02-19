#!/bin/bash
set -euo pipefail
# test_agents_map.sh - Verify AGENTS.md exists and has required links

FILE="AGENTS.md"

if [ ! -f "$FILE" ]; then
    echo "FAIL: $FILE does not exist"
    exit 1
fi

grep -q "# Springfield Agent Command Center" "$FILE" || { echo "FAIL: Missing Title"; exit 1; }
grep -q "README.md" "$FILE" || { echo "FAIL: Missing link to README.md"; exit 1; }
grep -q "QUICK_START.md" "$FILE" || { echo "FAIL: Missing link to QUICK_START.md"; exit 1; }
grep -q "docs/" "$FILE" || { echo "FAIL: Missing link to docs/"; exit 1; }
grep -q "LLM Guidance" "$FILE" || { echo "FAIL: Missing LLM Guidance section"; exit 1; }

echo "PASS: $FILE meets standards"
exit 0
