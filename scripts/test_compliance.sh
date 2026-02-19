#!/bin/bash
set -euo pipefail
# test_compliance.sh - Verify ADR-000 existence and sections

FILE="docs/adr/ADR-000-compliance-and-safety.md"

if [ ! -f "$FILE" ]; then
    echo "FAIL: $FILE does not exist"
    exit 1
fi

grep -q "ADR-000: Enterprise Compliance & Safety Standards" "$FILE" || { echo "FAIL: Missing Title"; exit 1; }
grep -q "Building Blocks" "$FILE" || { echo "FAIL: Missing Building Blocks reference"; exit 1; }
grep -q "RBAC" "$FILE" || { echo "FAIL: Missing RBAC reference"; exit 1; }
grep -q "Audit Logging" "$FILE" || { echo "FAIL: Missing Audit Logging reference"; exit 1; }

echo "PASS: $FILE meets standards"
exit 0
