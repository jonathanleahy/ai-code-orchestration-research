#!/usr/bin/env bash
# validate-gate.sh — Run a validation gate command in a working directory
#
# Usage: bash validate-gate.sh <gate_command> <workdir>
# Output: JSON to stdout: {"pass":true/false,"stdout":"...","stderr":"...","exit_code":N}

set -uo pipefail

GATE_CMD="${1:?Usage: validate-gate.sh <command> <workdir>}"
WORKDIR="${2:?Usage: validate-gate.sh <command> <workdir>}"

stdout_file=$(mktemp)
stderr_file=$(mktemp)
trap 'rm -f "$stdout_file" "$stderr_file"' EXIT

exit_code=0
(cd "$WORKDIR" && eval "$GATE_CMD") > "$stdout_file" 2> "$stderr_file" || exit_code=$?

stdout=$(head -c 2000 "$stdout_file" | python3 -c "import sys,json; print(json.dumps(sys.stdin.read()))" 2>/dev/null || echo '""')
stderr=$(head -c 2000 "$stderr_file" | python3 -c "import sys,json; print(json.dumps(sys.stdin.read()))" 2>/dev/null || echo '""')

pass="false"
[ "$exit_code" -eq 0 ] && pass="true"

echo "{\"pass\":$pass,\"exit_code\":$exit_code,\"stdout\":$stdout,\"stderr\":$stderr}"
