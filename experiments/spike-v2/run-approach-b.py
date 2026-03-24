#!/usr/bin/env python3
"""Approach B: Generate-and-Filter (AlphaCode pattern)

For each sub-task:
1. Generate N candidates with the cheapest model
2. Run gate on each — first to pass wins
3. If none pass, generate N more (up to 3 rounds)

Usage: python3 run-approach-b.py <plan.json> <attempt_dir> <assembled_dir> <spike_dir> <arch_file> <executor_model> <num_candidates>
"""

import json, os, shutil, subprocess, sys, time

plan_json = sys.argv[1]
attempt_dir = sys.argv[2]
assembled = sys.argv[3]
spike_dir = sys.argv[4]
arch_file = sys.argv[5]
executor_model = sys.argv[6]
num_candidates = int(sys.argv[7]) if len(sys.argv) > 7 else 5

plan = json.load(open(plan_json))
arch_content = open(arch_file).read()

stats = {"total": len(plan), "passed": 0, "failed": 0, "total_cost": 0.0, "total_calls": 0}


def get_context(max_chars=4000):
    ctx = ""
    for root, dirs, files in os.walk(assembled):
        for f in files:
            if f.endswith(('.cjs', '.json')) and not f.startswith('.'):
                fp = os.path.join(root, f)
                rel = os.path.relpath(fp, assembled)
                try:
                    entry = f"\n--- {rel} ---\n{open(fp).read()}\n"
                    if len(ctx) + len(entry) < max_chars:
                        ctx += entry
                except: pass
    return ctx or "None yet."


def call_model(model, prompt_file, workdir):
    try:
        r = subprocess.run(
            ["python3", os.path.join(spike_dir, "call-model.py"),
             "--model", model, "--prompt-file", prompt_file, "--workdir", workdir],
            capture_output=True, text=True, timeout=300
        )
        return json.loads(r.stdout.strip()) if r.stdout.strip() else {"error": "no output"}
    except subprocess.TimeoutExpired:
        return {"error": "timeout", "cost_usd": 0}
    except:
        return {"error": "call_failed", "cost_usd": 0}


def run_gate(gate_cmd, workdir):
    try:
        r = subprocess.run(
            ["bash", os.path.join(spike_dir, "validate-gate.sh"), gate_cmd, workdir],
            capture_output=True, text=True, timeout=30
        )
        data = json.loads(r.stdout.strip()) if r.stdout.strip() else {"pass": False}
        return data.get("pass", False), data
    except:
        return False, {"pass": False}


def copy_to_assembled(src_dir):
    copied = 0
    for root, dirs, files in os.walk(src_dir):
        for f in files:
            if f.endswith(('.cjs', '.json', '.txt')) and f not in ('response.json', 'prompt.txt'):
                src = os.path.join(root, f)
                rel = os.path.relpath(src, src_dir)
                dst = os.path.join(assembled, rel)
                os.makedirs(os.path.dirname(dst), exist_ok=True)
                shutil.copy2(src, dst)
                copied += 1
    return copied


# Build prompt template
prompt_template = open(os.path.join(spike_dir, "prompts", "executor.md")).read()

for i, task in enumerate(plan):
    task_id = task.get("id", f"ST-{i+1:02d}")
    title = task.get("title", "untitled")
    gate_cmd = task.get("gate_command", "")

    print(f"  Sub-task {task_id}: {title}", flush=True)

    # Build base prompt
    arch_excerpt = arch_content[:3000]
    for marker in task.get("files_to_create", []):
        idx = arch_content.find(f"### {marker}")
        if idx > 0:
            arch_excerpt = arch_content[max(0, idx-200):idx+1500]
            break

    base_prompt = prompt_template.replace("{{ARCHITECTURE}}", arch_excerpt)
    base_prompt = base_prompt.replace("{{SUB_TASK}}", json.dumps(task, indent=2))
    base_prompt = base_prompt.replace("{{CONTEXT_FILES}}", get_context())

    found_winner = False

    for round_num in range(3):  # Up to 3 rounds
        candidates_dir = os.path.join(attempt_dir, task_id, f"round-{round_num}")
        os.makedirs(candidates_dir, exist_ok=True)

        print(f"    Round {round_num}: generating {num_candidates} candidates...", flush=True)

        for c in range(num_candidates):
            cand_dir = os.path.join(candidates_dir, f"candidate-{c:02d}")
            os.makedirs(cand_dir, exist_ok=True)

            prompt_file = os.path.join(cand_dir, "prompt.txt")
            with open(prompt_file, "w") as f:
                f.write(base_prompt)

            metrics = call_model(executor_model, prompt_file, cand_dir)
            cost = metrics.get("cost_usd", 0)
            stats["total_cost"] += cost if isinstance(cost, (int, float)) else 0
            stats["total_calls"] += 1

            if metrics.get("files_created", 0) == 0:
                continue

            # Copy to assembled and test gate
            copy_to_assembled(cand_dir)

            if gate_cmd:
                gate_pass, _ = run_gate(gate_cmd, assembled)
            else:
                gate_pass = True

            if gate_pass:
                print(f"    Winner: candidate {c} (round {round_num}, ${cost:.4f})", flush=True)
                stats["passed"] += 1
                found_winner = True
                break

        if found_winner:
            break

    if not found_winner:
        print(f"    FAILED: no candidate passed after {3 * num_candidates} attempts", flush=True)
        stats["failed"] += 1

# Final validation
print(f"\n  === Final Validation ===", flush=True)
print(f"  Sub-tasks: {stats['passed']}/{stats['total']} passed", flush=True)
print(f"  Total calls: {stats['total_calls']}", flush=True)
print(f"  Total cost: ${stats['total_cost']:.4f}", flush=True)

test_file = os.path.join(assembled, "test", "dep-doctor.test.cjs")
if os.path.exists(test_file):
    try:
        r = subprocess.run(["node", test_file], capture_output=True, text=True, timeout=30, cwd=assembled)
        print(r.stdout, flush=True)
        pass_count = r.stdout.count("PASS")
        fail_count = r.stdout.count("FAIL")
        print(f"  FINAL: {pass_count} PASS, {fail_count} FAIL", flush=True)
    except:
        print("  FINAL: Tests failed to run", flush=True)
else:
    golden = os.path.join(spike_dir, "golden-master", "dep-doctor", "test", "dep-doctor.test.cjs")
    if os.path.exists(golden):
        try:
            r = subprocess.run(["node", golden], capture_output=True, text=True, timeout=30, cwd=assembled)
            print(r.stdout, flush=True)
            print(f"  FINAL (golden): {r.stdout.count('PASS')} PASS, {r.stdout.count('FAIL')} FAIL", flush=True)
        except:
            print("  FINAL: Golden tests failed", flush=True)

summary = {**stats, "executor_model": executor_model, "num_candidates": num_candidates}
with open(os.path.join(attempt_dir, "summary.json"), "w") as f:
    json.dump(summary, f, indent=2)
print(f"\n  Summary: {json.dumps(summary)}", flush=True)
