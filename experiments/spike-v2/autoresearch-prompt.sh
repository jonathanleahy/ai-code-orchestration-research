#!/usr/bin/env bash
# autoresearch-prompt.sh — Karpathy autoresearch loop for prompt optimization
#
# Goal: Find the prompt that makes Qwen3-30B ($0.0005/call) produce correct
# file blocks 100% of the time. If we get 100%, the cost per app drops to $0.017.
#
# Method: Try prompt variations, test each on 3 sub-tasks, record results.
# Keep variations that improve, discard those that don't.
#
# Usage: source ~/.config/gyrum/control-plane.env && bash scripts/dev-spike-v2/autoresearch-prompt.sh

set -euo pipefail

SPIKE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MODEL="qwen/qwen3-coder-30b-a3b-instruct"
RESULTS="$SPIKE_DIR/prompt-autoresearch-results.tsv"
ARCH=$(cat "$SPIKE_DIR/architecture.md")

[ -z "${OPENROUTER_API_KEY:-}" ] && { echo "Set OPENROUTER_API_KEY"; exit 1; }

log() { echo "[$(date '+%H:%M:%S')] $*"; }

# Initialize results
[ ! -f "$RESULTS" ] && echo -e "variation\tpass_rate\tcost\tformat_ok\tfiles_created\tnotes\ttimestamp" > "$RESULTS"

# Test task: write lib/validator.cjs
TEST_SUBTASK='{
  "id": "ST-02",
  "title": "Implement lib/validator.cjs",
  "files_to_create": ["lib/validator.cjs"],
  "gate_command": "node -e \"const v=require('"'"'./lib/validator.cjs'"'"'); if(typeof v.isValidSemver!=='"'"'function'"'"') process.exit(1)\"",
  "description": "Semver and SPDX validation module"
}'

# Architecture excerpt for validator
ARCH_EXCERPT=$(echo "$ARCH" | sed -n '/### lib\/validator/,/### lib\/parser/p' | head -30)

# Prompt variations to test
declare -a VARIATIONS=(
# V1: Original (baseline)
"You are implementing lib/validator.cjs for the dep-doctor CLI.

$ARCH_EXCERPT

Output the complete file:

--- FILE: lib/validator.cjs ---
[your code]
--- END FILE ---

Rules: .cjs extension, 'use strict', zero npm deps, module.exports."

# V2: No dashes, use equals signs
"Write lib/validator.cjs.

$ARCH_EXCERPT

=== FILE: lib/validator.cjs ===
[complete code here]
=== END FILE ===

Output only the file block above filled in. No explanation."

# V3: JSON wrapper
"Output a JSON object with one key 'files' containing an array of {path, content} objects.

Write lib/validator.cjs based on:
$ARCH_EXCERPT

Output: {\"files\": [{\"path\": \"lib/validator.cjs\", \"content\": \"...\"}]}"

# V4: Very explicit format with example
"Write lib/validator.cjs for dep-doctor.

$ARCH_EXCERPT

YOUR OUTPUT MUST LOOK EXACTLY LIKE THIS (replace the ... with real code):

--- FILE: lib/validator.cjs ---
'use strict';

const SEMVER_RANGE = /^[\\^~>=<]*\\d+(\\.\\d+){0,2}([-.][a-zA-Z0-9]+)*\$/;
...more code...
module.exports = { isValidSemver, isValidSpdx, validateDependency };
--- END FILE ---

IMPORTANT: Start with --- FILE: and end with --- END FILE ---
Do NOT wrap in \`\`\`javascript fences. Output ONLY the file block."

# V5: Minimal — just the function signatures
"Create lib/validator.cjs exporting:
- isValidSemver(range) → boolean
- isValidSpdx(license) → boolean
- validateDependency(name, range) → [{type, message}]

Zero npm deps. Node.js built-ins only. CommonJS.

--- FILE: lib/validator.cjs ---
[code]
--- END FILE ---"

# V6: Roleplay as file system
"You are a file system. When asked to write a file, output its contents between markers.

Write lib/validator.cjs:
$ARCH_EXCERPT

BEGIN_FILE lib/validator.cjs
[output the code]
END_FILE"

# V7: Think step by step then output
"Think about what lib/validator.cjs needs:
$ARCH_EXCERPT

After thinking, output the file:
--- FILE: lib/validator.cjs ---
[code]
--- END FILE ---

Think briefly, then output the file block."

# V8: Repeat format 3 times
"OUTPUT FORMAT (read this 3 times):
--- FILE: lib/validator.cjs ---
code here
--- END FILE ---

--- FILE: lib/validator.cjs ---
code here
--- END FILE ---

--- FILE: lib/validator.cjs ---
code here
--- END FILE ---

Now actually write the code. $ARCH_EXCERPT

Output using the format above (--- FILE: ... --- END FILE ---)."
)

log "=== Autoresearch: Prompt Optimization for $MODEL ==="
log "Testing ${#VARIATIONS[@]} variations, 3 runs each"
log ""

best_rate=0
best_var=0

for i in "${!VARIATIONS[@]}"; do
    var_num=$((i + 1))
    prompt="${VARIATIONS[$i]}"
    pass=0
    total=3
    format_ok=0
    files=0
    cost_total=0

    log "--- Variation $var_num ---"

    for run in 1 2 3; do
        workdir=$(mktemp -d)
        mkdir -p "$workdir/lib"

        # Write prompt
        echo "$prompt" > "$workdir/prompt.txt"

        # Call model
        metrics=$(python3 "$SPIKE_DIR/call-model.py" \
            --model "$MODEL" \
            --prompt-file "$workdir/prompt.txt" \
            --workdir "$workdir" 2>/dev/null || echo '{"error":"failed"}')

        cost=$(echo "$metrics" | python3 -c "import sys,json; print(json.load(sys.stdin).get('cost_usd',0))" 2>/dev/null || echo 0)
        cost_total=$(echo "$cost_total + $cost" | bc 2>/dev/null || echo "$cost_total")
        fc=$(echo "$metrics" | python3 -c "import sys,json; print(json.load(sys.stdin).get('files_created',0))" 2>/dev/null || echo 0)
        files=$((files + fc))

        # Check format: does output have --- FILE: marker?
        if [ -f "$workdir/response.json" ]; then
            has_marker=$(python3 -c "
import json
data = json.load(open('$workdir/response.json'))
content = data.get('choices',[{}])[0].get('message',{}).get('content','')
print('1' if '--- FILE:' in content or '=== FILE:' in content or 'BEGIN_FILE' in content else '0')
" 2>/dev/null || echo "0")
            [ "$has_marker" = "1" ] && format_ok=$((format_ok + 1))
        fi

        # Check gate
        if [ -f "$workdir/lib/validator.cjs" ]; then
            gate_result=$(cd "$workdir" && node -e "const v=require('./lib/validator.cjs'); if(typeof v.isValidSemver!=='function') process.exit(1)" 2>&1 && echo "PASS" || echo "FAIL")
            [ "$gate_result" = "PASS" ] && pass=$((pass + 1))
        fi

        rm -rf "$workdir"
    done

    rate=$((pass * 100 / total))
    log "  V$var_num: pass=$pass/$total (${rate}%) format_ok=$format_ok files=$files cost=\$$cost_total"

    echo -e "V${var_num}\t${rate}%\t${cost_total}\t${format_ok}/3\t${files}\t\t$(date -u '+%Y-%m-%dT%H:%M:%SZ')" >> "$RESULTS"

    if [ "$rate" -gt "$best_rate" ]; then
        best_rate=$rate
        best_var=$var_num
    fi
done

log ""
log "=== RESULTS ==="
log "Best variation: V$best_var (${best_rate}% pass rate)"
log ""
cat "$RESULTS"
