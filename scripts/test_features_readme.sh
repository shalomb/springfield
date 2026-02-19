#!/bin/bash
# test_features_readme.sh - Verify Features README structure

FILE="docs/features/README.md"

if [ ! -f "$FILE" ]; then
    echo "FAIL: $FILE does not exist"
    exit 1
fi

grep -q "# Springfield Features & BDD" "$FILE" || { echo "FAIL: Missing Title"; exit 1; }
grep -q "## BDD Workflow" "$FILE" || { echo "FAIL: Missing BDD Workflow section"; exit 1; }
grep -q "## Role of Agents" "$FILE" || { echo "FAIL: Missing Role of Agents section"; exit 1; }
grep -q "## Running Tests" "$FILE" || { echo "FAIL: Missing Running Tests section"; exit 1; }
grep -q "just" "$FILE" || { echo "FAIL: Missing reference to just command"; exit 1; }

echo "PASS: $FILE meets standards"
exit 0
