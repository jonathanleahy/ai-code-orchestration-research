#!/usr/bin/env python3
"""Approach C: LLM Council

For each sub-task:
1. 3 models generate independently
2. Run gate on each
3. Pick the one that passes (or best score if multiple pass)
If none pass, use chairman to synthesize from all 3.

Usage: python3 run-approach-c.py <plan.json> <attempt_dir> <assembled_dir> <spike_dir> <arch_file> <model_a> <model_b> <model_c> <chairman>
"""

import json, os, shutil, subprocess, sys

plan_json = sys.argv[1]
attempt_dir = sys.argv[2]
assembled = sys.argv[3]
spike_dir = sys.argv[4]
arch_file = sys.argv[5]
model_a = sys.argv[6]
model_b = sys.argv[7]
model_c = sys.argv[8]
chairman = sys.argv[9]

plan = json.load(open(plan_json))
arch_content = open(arch_file).read()
council = [("A", model_a), ("B", model_b), ("C", model_c)]

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


exec_template = open(os.path.join(spike_dir, "prompts", "executor.md")).read()

for i, task in enumerate(plan):
    task_id = task.get("id", f"ST-{i+1:02d}")
    title = task.get("title", "untitled")
    gate_cmd = task.get("gate_command", "")
    task_dir = os.path.join(attempt_dir, task_id)

    print(f"  Sub-task {task_id}: {title}", flush=True)

    arch_excerpt = arch_content[:3000]
    for marker in task.get("files_to_create", []):
        idx = arch_content.find(f"### {marker}")
        if idx > 0:
            arch_excerpt = arch_content[max(0, idx-200):idx+1500]
            break

    context = get_context()
    prompt = exec_template.replace("{{ARCHITECTURE}}", arch_excerpt)
    prompt = prompt.replace("{{SUB_TASK}}", json.dumps(task, indent=2))
    prompt = prompt.replace("{{CONTEXT_FILES}}", context)

    # Generate with all 3 council members
    results = []
    for label, model in council:
        cand_dir = os.path.join(task_dir, f"council-{label}")
        os.makedirs(cand_dir, exist_ok=True)

        pf = os.path.join(cand_dir, "prompt.txt")
        with open(pf, "w") as f:
            f.write(prompt)

        metrics = call_model(model, pf, cand_dir)
        cost = metrics.get("cost_usd", 0)
        stats["total_cost"] += cost if isinstance(cost, (int, float)) else 0
        stats["total_calls"] += 1

        # Test gate
        copy_to_assembled(cand_dir)
        if gate_cmd:
            gate_pass, gate_data = run_gate(gate_cmd, assembled)
        else:
            gate_pass = metrics.get("files_created", 0) > 0

        results.append({
            "label": label, "model": model, "dir": cand_dir,
            "gate_pass": gate_pass, "files": metrics.get("files_created", 0),
            "cost": cost
        })

        print(f"    Council {label} ({model.split('/')[-1]}): gate={'PASS' if gate_pass else 'FAIL'}, files={metrics.get('files_created', 0)}", flush=True)

    # Pick winner: first that passes gate, or the one with most files
    winners = [r for r in results if r["gate_pass"]]
    if winners:
        winner = winners[0]
        copy_to_assembled(winner["dir"])
        print(f"    Winner: {winner['label']} ({winner['model'].split('/')[-1]})", flush=True)
        stats["passed"] += 1
    else:
        # No one passed — pick best and mark failed
        results.sort(key=lambda r: -r["files"])
        if results[0]["files"] > 0:
            copy_to_assembled(results[0]["dir"])
        print(f"    NO WINNER — best effort from {results[0]['label']}", flush=True)
        stats["failed"] += 1

# Final validation
print(f"\n  === Final Validation ===", flush=True)
print(f"  Sub-tasks: {stats['passed']}/{stats['total']}", flush=True)
print(f"  Total calls: {stats['total_calls']}", flush=True)
print(f"  Total cost: ${stats['total_cost']:.4f}", flush=True)

test_file = os.path.join(assembled, "test", "dep-doctor.test.cjs")
if os.path.exists(test_file):
    try:
        r = subprocess.run(["node", test_file], capture_output=True, text=True, timeout=30, cwd=assembled)
        print(r.stdout, flush=True)
        print(f"  FINAL: {r.stdout.count('PASS')} PASS, {r.stdout.count('FAIL')} FAIL", flush=True)
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

summary = {**stats, "council": [m for _, m in council], "chairman": chairman}
with open(os.path.join(attempt_dir, "summary.json"), "w") as f:
    json.dump(summary, f, indent=2)
print(f"\n  Summary: {json.dumps(summary)}", flush=True)
