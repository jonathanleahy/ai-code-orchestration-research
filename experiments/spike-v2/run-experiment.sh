#!/usr/bin/env bash
# run-experiment.sh — Run a single experiment config
#
# Usage:
#   bash run-experiment.sh --approach A --config A1
#   bash run-experiment.sh --approach A --config A2
#   bash run-experiment.sh --approach B --config B1
#
# Approach A: 4-layer quality-first (plan → execute → review)
# Approach B: generate-and-filter (5 plans → tests first → 10 candidates)
# Approach C: LLM council (3 models generate → peer review → chairman picks)
# Approach D: evolutionary (generate → test → select → mutate → repeat)

set -euo pipefail

SPIKE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ARCH_FILE="$SPIKE_DIR/architecture.md"
GOLDEN_DIR="$SPIKE_DIR/golden-master/dep-doctor"
RESULTS_FILE="$SPIKE_DIR/results.tsv"

log() { echo "[$(date '+%H:%M:%S')] $*"; }

# Parse args
APPROACH=""
CONFIG=""
while [ $# -gt 0 ]; do
    case "$1" in
        --approach) APPROACH="$2"; shift 2 ;;
        --config) CONFIG="$2"; shift 2 ;;
        *) echo "Unknown arg: $1"; exit 1 ;;
    esac
done
[ -z "$APPROACH" ] || [ -z "$CONFIG" ] && { echo "Usage: --approach A|B|C|D --config A1|B1|..."; exit 1; }

# Initialize results
[ ! -f "$RESULTS_FILE" ] && echo -e "approach\tconfig\tattempt\tlayer\tsub_task\tmodel\tcost_usd\ttime_s\ttokens_in\ttokens_out\tgate_pass\treview_verdict\tquality\tretry\tfiles\terror\ttimestamp" > "$RESULTS_FILE"

# Find next attempt number
ATTEMPT=$(printf "%03d" $(grep -c "^${APPROACH}	${CONFIG}" "$RESULTS_FILE" 2>/dev/null || echo 0))
ATTEMPT_DIR="$SPIKE_DIR/attempts/${CONFIG}-${ATTEMPT}"
mkdir -p "$ATTEMPT_DIR"

log "=== Experiment: Approach $APPROACH, Config $CONFIG, Attempt $ATTEMPT ==="

# Model configs
declare -A MODELS
case "$CONFIG" in
    A1) MODELS=([planner]="anthropic/claude-sonnet-4" [executor]="anthropic/claude-sonnet-4" [reviewer]="anthropic/claude-sonnet-4") ;;
    A2) MODELS=([planner]="google/gemini-2.5-flash" [executor]="qwen/qwen3-coder" [reviewer]="google/gemini-2.5-flash") ;;
    A3) MODELS=([planner]="google/gemini-2.5-flash" [executor]="minimax/minimax-m2.7" [reviewer]="deepseek/deepseek-v3.2") ;;
    A4) MODELS=([planner]="anthropic/claude-sonnet-4" [executor]="qwen/qwen3-coder-30b-a3b-instruct" [reviewer]="anthropic/claude-sonnet-4") ;;
    B1) MODELS=([planner]="google/gemini-2.5-flash" [test_writer]="google/gemini-2.5-flash" [executor]="qwen/qwen3-coder-30b-a3b-instruct") ;;
    B2) MODELS=([planner]="claude-sonnet" [test_writer]="claude-sonnet" [executor]="qwen/qwen3-coder-30b-a3b-instruct") ;;
    C1) MODELS=([planner]="google/gemini-2.5-flash" [council_a]="qwen/qwen3-coder" [council_b]="minimax/minimax-m2.7" [council_c]="deepseek/deepseek-v3.2" [chairman]="google/gemini-2.5-flash") ;;
    C2) MODELS=([planner]="claude-sonnet" [council_a]="qwen/qwen3-coder" [council_b]="minimax/minimax-m2.7" [council_c]="google/gemini-2.5-flash" [chairman]="claude-sonnet") ;;
    D1) MODELS=([planner]="google/gemini-2.5-flash" [mutator]="qwen/qwen3-coder-30b-a3b-instruct") ;;
    D2) MODELS=([planner]="claude-sonnet" [mutator]="qwen/qwen3-coder-30b-a3b-instruct") ;;
    *) echo "Unknown config: $CONFIG"; exit 1 ;;
esac

ARCH_CONTENT=$(cat "$ARCH_FILE")

# Helper: call a model and log results
call_and_log() {
    local layer="$1" sub_task="$2" model="$3" prompt_file="$4" workdir="$5" retry="${6:-0}"
    local start_ts=$(date +%s)

    local metrics
    metrics=$(python3 "$SPIKE_DIR/call-model.py" --model "$model" --prompt-file "$prompt_file" --workdir "$workdir" 2>/dev/null || echo '{"error":"call failed"}')

    local elapsed=$(( $(date +%s) - start_ts ))
    local cost=$(echo "$metrics" | python3 -c "import sys,json; print(json.load(sys.stdin).get('cost_usd',0))" 2>/dev/null || echo 0)
    local tokens_in=$(echo "$metrics" | python3 -c "import sys,json; print(json.load(sys.stdin).get('tokens_in',0))" 2>/dev/null || echo 0)
    local tokens_out=$(echo "$metrics" | python3 -c "import sys,json; print(json.load(sys.stdin).get('tokens_out',0))" 2>/dev/null || echo 0)
    local files=$(echo "$metrics" | python3 -c "import sys,json; print(json.load(sys.stdin).get('files_created',0))" 2>/dev/null || echo 0)
    local error=$(echo "$metrics" | python3 -c "import sys,json; print(json.load(sys.stdin).get('error',''))" 2>/dev/null || echo "")

    # Run gate if workdir has files
    local gate_pass="n/a"
    local review_verdict="n/a"
    local quality="-"

    echo -e "${APPROACH}\t${CONFIG}\t${ATTEMPT}\t${layer}\t${sub_task}\t${model}\t${cost}\t${elapsed}\t${tokens_in}\t${tokens_out}\t${gate_pass}\t${review_verdict}\t${quality}\t${retry}\t${files}\t${error}\t$(date -u '+%Y-%m-%dT%H:%M:%SZ')" >> "$RESULTS_FILE"

    log "  $layer/$sub_task: model=$model cost=\$$cost time=${elapsed}s files=$files ${error:+ERROR: $error}"
}

# Helper: run a gate command
run_gate() {
    local gate_cmd="$1" workdir="$2"
    local result
    result=$(bash "$SPIKE_DIR/validate-gate.sh" "$gate_cmd" "$workdir" 2>/dev/null || echo '{"pass":false}')
    echo "$result" | python3 -c "import sys,json; print(json.load(sys.stdin).get('pass', False))" 2>/dev/null
}

# ============================================================================
# APPROACH A: 4-Layer Quality-First
# ============================================================================
run_approach_a() {
    log "--- Approach A: 4-Layer Quality-First ---"
    local planner="${MODELS[planner]}"
    local executor="${MODELS[executor]}"
    local reviewer="${MODELS[reviewer]}"

    # Layer 2: Plan
    log "Layer 2: Planning ($planner)"
    local plan_dir="$ATTEMPT_DIR/plan"
    mkdir -p "$plan_dir"
    local plan_prompt="$plan_dir/prompt.txt"
    python3 -c "
t = open('$SPIKE_DIR/prompts/planner.md').read()
a = open('$ARCH_FILE').read()
open('$plan_prompt', 'w').write(t.replace('{{ARCHITECTURE}}', a))
"
    call_and_log "planner" "plan" "$planner" "$plan_prompt" "$plan_dir"

    # Parse plan JSON
    local plan_json="$plan_dir/plan.json"
    # Try to find JSON in the response
    if [ -f "$plan_dir/response.json" ]; then
        python3 -c "
import json
data = json.load(open('$plan_dir/response.json'))
content = data.get('choices',[{}])[0].get('message',{}).get('content','')
# Find JSON array in content
import re
match = re.search(r'\[.*\]', content, re.DOTALL)
if match:
    tasks = json.loads(match.group())
    json.dump(tasks, open('$plan_json', 'w'), indent=2)
    print(f'Parsed {len(tasks)} sub-tasks')
else:
    print('ERROR: No JSON array found in planner output')
" 2>/dev/null || echo "Plan parse failed"
    fi

    if [ ! -f "$plan_json" ]; then
        log "  PLAN FAILED — no valid sub-tasks JSON"
        return 1
    fi

    local task_count
    task_count=$(python3 -c "import json; print(len(json.load(open('$plan_json'))))" 2>/dev/null || echo 0)
    log "  Plan: $task_count sub-tasks"

    # Assembled output directory
    local assembled="$ATTEMPT_DIR/assembled/dep-doctor"
    mkdir -p "$assembled/lib" "$assembled/test" "$assembled/fixtures/valid" "$assembled/fixtures/malformed" "$assembled/fixtures/empty"

    # Layer 3+4: Execute + Review each sub-task
    python3 "$SPIKE_DIR/run-subtasks.py" \
        "$plan_json" "$ATTEMPT_DIR" "$assembled" "$SPIKE_DIR" "$ARCH_FILE" "$executor" "$reviewer" 2>&1

    log "--- Approach A complete ---"
}

# ============================================================================
# APPROACH B: Generate-and-Filter
# ============================================================================
run_approach_b() {
    log "--- Approach B: Generate-and-Filter ---"
    local planner="${MODELS[planner]:-${MODELS[plan_gen]:-google/gemini-2.5-flash}}"
    local executor="${MODELS[executor]:-qwen/qwen3-coder-30b-a3b-instruct}"
    local num_candidates=5

    # Step 1: Plan (reuse Approach A planner)
    log "Layer 2: Planning ($planner)"
    local plan_dir="$ATTEMPT_DIR/plan"
    mkdir -p "$plan_dir"
    local plan_prompt="$plan_dir/prompt.txt"
    python3 -c "
t = open('$SPIKE_DIR/prompts/planner.md').read()
a = open('$ARCH_FILE').read()
open('$plan_prompt', 'w').write(t.replace('{{ARCHITECTURE}}', a))
"
    call_and_log "planner" "plan" "$planner" "$plan_prompt" "$plan_dir"

    local plan_json="$plan_dir/plan.json"
    python3 -c "
import json, re
data = json.load(open('$plan_dir/response.json'))
content = data.get('choices',[{}])[0].get('message',{}).get('content','')
match = re.search(r'\[.*\]', content, re.DOTALL)
if match:
    tasks = json.loads(match.group())
    json.dump(tasks, open('$plan_json', 'w'), indent=2)
    print(f'Parsed {len(tasks)} sub-tasks')
else:
    print('ERROR: No JSON array found')
" 2>/dev/null || echo "Plan parse failed"

    if [ ! -f "$plan_json" ]; then
        log "  PLAN FAILED"
        return 1
    fi

    local assembled="$ATTEMPT_DIR/assembled/dep-doctor"
    mkdir -p "$assembled/lib" "$assembled/test" "$assembled/fixtures/valid" "$assembled/fixtures/malformed" "$assembled/fixtures/empty"

    # Step 2: Generate-and-filter execution
    python3 "$SPIKE_DIR/run-approach-b.py" \
        "$plan_json" "$ATTEMPT_DIR" "$assembled" "$SPIKE_DIR" "$ARCH_FILE" "$executor" "$num_candidates" 2>&1

    log "--- Approach B complete ---"
}

# ============================================================================
# APPROACH C: LLM Council
# ============================================================================
run_approach_c() {
    log "--- Approach C: LLM Council ---"
    local planner="${MODELS[planner]:-google/gemini-2.5-flash}"
    local council_a="${MODELS[council_a]:-qwen/qwen3-coder}"
    local council_b="${MODELS[council_b]:-minimax/minimax-m2.7}"
    local council_c="${MODELS[council_c]:-deepseek/deepseek-v3.2}"
    local chairman="${MODELS[chairman]:-google/gemini-2.5-flash}"

    # Plan (same as Approach A)
    log "Layer 2: Planning ($planner)"
    local plan_dir="$ATTEMPT_DIR/plan"
    mkdir -p "$plan_dir"
    local plan_prompt="$plan_dir/prompt.txt"
    python3 -c "
t = open('$SPIKE_DIR/prompts/planner.md').read()
a = open('$ARCH_FILE').read()
open('$plan_prompt', 'w').write(t.replace('{{ARCHITECTURE}}', a))
"
    call_and_log "planner" "plan" "$planner" "$plan_prompt" "$plan_dir"

    local plan_json="$plan_dir/plan.json"
    python3 -c "
import json, re
data = json.load(open('$plan_dir/response.json'))
content = data.get('choices',[{}])[0].get('message',{}).get('content','')
match = re.search(r'\[.*\]', content, re.DOTALL)
if match:
    tasks = json.loads(match.group())
    json.dump(tasks, open('$plan_json', 'w'), indent=2)
    print(f'Parsed {len(tasks)} sub-tasks')
" 2>/dev/null

    if [ ! -f "$plan_json" ]; then
        log "  PLAN FAILED"
        return 1
    fi

    local assembled="$ATTEMPT_DIR/assembled/dep-doctor"
    mkdir -p "$assembled/lib" "$assembled/test" "$assembled/fixtures/valid" "$assembled/fixtures/malformed" "$assembled/fixtures/empty"

    # Council execution: 3 models generate, pick best per gate
    python3 "$SPIKE_DIR/run-approach-c.py" \
        "$plan_json" "$ATTEMPT_DIR" "$assembled" "$SPIKE_DIR" "$ARCH_FILE" \
        "$council_a" "$council_b" "$council_c" "$chairman" 2>&1

    log "--- Approach C complete ---"
}

# ============================================================================
# APPROACH D: Evolutionary
# ============================================================================
run_approach_d() {
    log "--- Approach D: Evolutionary ---"
    local planner="${MODELS[planner]:-google/gemini-2.5-flash}"
    local mutator="${MODELS[mutator]:-qwen/qwen3-coder-30b-a3b-instruct}"

    # Plan
    log "Layer 2: Planning ($planner)"
    local plan_dir="$ATTEMPT_DIR/plan"
    mkdir -p "$plan_dir"
    local plan_prompt="$plan_dir/prompt.txt"
    python3 -c "
t = open('$SPIKE_DIR/prompts/planner.md').read()
a = open('$ARCH_FILE').read()
open('$plan_prompt', 'w').write(t.replace('{{ARCHITECTURE}}', a))
"
    call_and_log "planner" "plan" "$planner" "$plan_prompt" "$plan_dir"

    local plan_json="$plan_dir/plan.json"
    python3 -c "
import json, re
data = json.load(open('$plan_dir/response.json'))
content = data.get('choices',[{}])[0].get('message',{}).get('content','')
match = re.search(r'\[.*\]', content, re.DOTALL)
if match:
    tasks = json.loads(match.group())
    json.dump(tasks, open('$plan_json', 'w'), indent=2)
    print(f'Parsed {len(tasks)} sub-tasks')
" 2>/dev/null

    if [ ! -f "$plan_json" ]; then
        log "  PLAN FAILED"
        return 1
    fi

    local assembled="$ATTEMPT_DIR/assembled/dep-doctor"
    mkdir -p "$assembled/lib" "$assembled/test" "$assembled/fixtures/valid" "$assembled/fixtures/malformed" "$assembled/fixtures/empty"

    python3 "$SPIKE_DIR/run-approach-d.py" \
        "$plan_json" "$ATTEMPT_DIR" "$assembled" "$SPIKE_DIR" "$ARCH_FILE" \
        "$mutator" 3 3 2>&1

    log "--- Approach D complete ---"
}

# Run the selected approach
case "$APPROACH" in
    A) run_approach_a ;;
    B) run_approach_b ;;
    C) run_approach_c ;;
    D) run_approach_d ;;
    *) echo "Unknown approach: $APPROACH"; exit 1 ;;
esac

log "=== Experiment complete: $CONFIG-$ATTEMPT ==="
log "Results: $RESULTS_FILE"
