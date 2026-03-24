#!/usr/bin/env bash
# call-claude-p.sh — Call claude -p with imperative prompts (subscription = free)
#
# Usage: bash call-claude-p.sh <prompt_file> <workdir> [model] [budget]
# Output: JSON metrics to stdout (same format as call-model.py)

set -uo pipefail

PROMPT_FILE="${1:?Usage: call-claude-p.sh <prompt_file> <workdir> [model] [budget]}"
WORKDIR="${2:?Usage: call-claude-p.sh <prompt_file> <workdir> [model] [budget]}"
MODEL="${3:-sonnet}"
BUDGET="${4:-0.30}"

PROMPT=$(cat "$PROMPT_FILE")
LOG_FILE="$WORKDIR/claude-p.log"

start_ts=$(date +%s)

# Run claude -p from the workdir with permissions
cd "$WORKDIR" && claude -p \
    --dangerously-skip-permissions \
    --model "$MODEL" \
    --max-budget-usd "$BUDGET" \
    --verbose \
    --output-format stream-json \
    "$PROMPT" > "$LOG_FILE" 2>&1 || true

cd - > /dev/null 2>&1

elapsed=$(( $(date +%s) - start_ts ))

# Extract cost
cost=$(grep -o '"total_cost_usd":[0-9.]*' "$LOG_FILE" 2>/dev/null | tail -1 | cut -d: -f2 || echo "0")

# Count files created (exclude the log)
files_created=0
while IFS= read -r f; do
    files_created=$((files_created + 1))
done < <(find "$WORKDIR" -name "*.go" -o -name "*.graphql" -o -name "*.ts" -o -name "*.svelte" -o -name "*.json" 2>/dev/null | grep -v "claude-p.log\|response.json\|prompt.txt")

# Output metrics (same format as call-model.py)
echo "{\"model\":\"claude-${MODEL}\",\"cost_usd\":${cost:-0},\"tokens_in\":0,\"tokens_out\":0,\"files_created\":${files_created},\"file_paths\":[],\"time_s\":${elapsed},\"error\":\"\"}"
