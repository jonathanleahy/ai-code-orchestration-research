#!/usr/bin/env python3
"""Approach D: Evolutionary / Genetic

For each sub-task:
1. Gen 0: Generate pop_size candidates
2. Score each by running gate + counting passing checks
3. Keep top 2 survivors
4. Gen 1-N: Mutate survivors ("improve this code, fix these issues")
5. Best survivor after all generations wins

Usage: python3 run-approach-d.py <plan.json> <attempt_dir> <assembled_dir> <spike_dir> <arch_file> <mutator_model> <generations> <pop_size>
"""

import json, os, shutil, subprocess, sys, time

plan_json = sys.argv[1]
attempt_dir = sys.argv[2]
assembled = sys.argv[3]
spike_dir = sys.argv[4]
arch_file = sys.argv[5]
mutator_model = sys.argv[6]
generations = int(sys.argv[7]) if len(sys.argv) > 7 else 3
pop_size = int(sys.argv[8]) if len(sys.argv) > 8 else 3

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


def score_candidate(cand_dir, gate_cmd):
    """Score: 0 (broken) to 10 (passes gate + has content)."""
    # Check files exist
    files = [f for f in os.listdir(cand_dir) if f.endswith(('.cjs', '.json')) and f not in ('response.json', 'prompt.txt')]
    if not files:
        return 0

    score = 2  # Has files

    # Check file sizes (non-empty)
    for f in files:
        fp = os.path.join(cand_dir, f)
        if os.path.getsize(fp) > 50:
            score += 1

    # Copy to assembled temporarily and run gate
    copy_to_assembled(cand_dir)
    if gate_cmd:
        gate_pass, gate_data = run_gate(gate_cmd, assembled)
        if gate_pass:
            score += 5  # Big bonus for passing gate
    else:
        score += 3

    return min(score, 10)


exec_template = open(os.path.join(spike_dir, "prompts", "executor.md")).read()
mutator_template = open(os.path.join(spike_dir, "prompts", "mutator.md")).read()

for i, task in enumerate(plan):
    task_id = task.get("id", f"ST-{i+1:02d}")
    title = task.get("title", "untitled")
    gate_cmd = task.get("gate_command", "")
    task_dir = os.path.join(attempt_dir, task_id)
    os.makedirs(task_dir, exist_ok=True)

    print(f"  Sub-task {task_id}: {title}", flush=True)

    arch_excerpt = arch_content[:3000]
    for marker in task.get("files_to_create", []):
        idx = arch_content.find(f"### {marker}")
        if idx > 0:
            arch_excerpt = arch_content[max(0, idx-200):idx+1500]
            break

    context = get_context()
    best_score = 0
    best_dir = None

    for gen in range(generations):
        gen_dir = os.path.join(task_dir, f"gen-{gen:02d}")
        os.makedirs(gen_dir, exist_ok=True)

        candidates = []

        if gen == 0:
            # Gen 0: fresh generation
            for c in range(pop_size):
                cand_dir = os.path.join(gen_dir, f"cand-{c:02d}")
                os.makedirs(cand_dir, exist_ok=True)

                prompt = exec_template.replace("{{ARCHITECTURE}}", arch_excerpt)
                prompt = prompt.replace("{{SUB_TASK}}", json.dumps(task, indent=2))
                prompt = prompt.replace("{{CONTEXT_FILES}}", context)

                pf = os.path.join(cand_dir, "prompt.txt")
                with open(pf, "w") as f:
                    f.write(prompt)

                metrics = call_model(mutator_model, pf, cand_dir)
                cost = metrics.get("cost_usd", 0)
                stats["total_cost"] += cost if isinstance(cost, (int, float)) else 0
                stats["total_calls"] += 1

                score = score_candidate(cand_dir, gate_cmd)
                candidates.append((score, cand_dir))
        else:
            # Gen 1+: mutate top survivors
            for c, (_, parent_dir) in enumerate(survivors):
                cand_dir = os.path.join(gen_dir, f"cand-{c:02d}")
                os.makedirs(cand_dir, exist_ok=True)

                # Read parent code
                parent_code = ""
                for root, dirs, files in os.walk(parent_dir):
                    for f in files:
                        if f.endswith(('.cjs', '.json')) and f not in ('response.json', 'prompt.txt'):
                            fp = os.path.join(root, f)
                            rel = os.path.relpath(fp, parent_dir)
                            parent_code += f"\n--- {rel} ---\n{open(fp).read()}\n"

                # Get gate error if available
                test_results = ""
                if gate_cmd:
                    _, gate_data = run_gate(gate_cmd, parent_dir)
                    test_results = gate_data.get("stderr", "")[:500]

                prompt = mutator_template.replace("{{ARCHITECTURE_EXCERPT}}", arch_excerpt)
                prompt = prompt.replace("{{CURRENT_CODE}}", parent_code[:4000])
                prompt = prompt.replace("{{TEST_RESULTS}}", test_results or "Gate passed but may have quality issues.")
                # Replace {{FILE_PATH}} with the primary file
                files_to_create = task.get("files_to_create", ["unknown.cjs"])
                prompt = prompt.replace("{{FILE_PATH}}", files_to_create[0] if files_to_create else "output.cjs")

                pf = os.path.join(cand_dir, "prompt.txt")
                with open(pf, "w") as f:
                    f.write(prompt)

                metrics = call_model(mutator_model, pf, cand_dir)
                cost = metrics.get("cost_usd", 0)
                stats["total_cost"] += cost if isinstance(cost, (int, float)) else 0
                stats["total_calls"] += 1

                score = score_candidate(cand_dir, gate_cmd)
                candidates.append((score, cand_dir))

        # Select top 2 survivors
        candidates.sort(key=lambda x: -x[0])
        survivors = candidates[:2]
        gen_best = candidates[0][0] if candidates else 0

        print(f"    Gen {gen}: best={gen_best}, scores={[s for s,_ in candidates]}", flush=True)

        if gen_best > best_score:
            best_score = gen_best
            best_dir = candidates[0][1]

        # If gate passes, we're done
        if gen_best >= 7:
            break

    if best_dir and best_score >= 5:
        copy_to_assembled(best_dir)
        print(f"    Winner: score={best_score}", flush=True)
        stats["passed"] += 1
    else:
        print(f"    FAILED: best score={best_score}", flush=True)
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

summary = {**stats, "mutator_model": mutator_model, "generations": generations, "pop_size": pop_size}
with open(os.path.join(attempt_dir, "summary.json"), "w") as f:
    json.dump(summary, f, indent=2)
print(f"\n  Summary: {json.dumps(summary)}", flush=True)
