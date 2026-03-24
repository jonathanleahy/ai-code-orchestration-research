#!/usr/bin/env python3
"""Execute sub-tasks from a plan, with structural gates and assembly.

Usage: python3 run-subtasks.py <plan.json> <attempt_dir> <assembled_dir> <spike_dir> <arch_file> <executor_model> <reviewer_model>
"""

import json
import os
import shutil
import subprocess
import sys
import time

plan_json = sys.argv[1]
attempt_dir = sys.argv[2]
assembled = sys.argv[3]
spike_dir = sys.argv[4]
arch_file = sys.argv[5]
executor_model = sys.argv[6]
reviewer_model = sys.argv[7]

plan = json.load(open(plan_json))
arch_content = open(arch_file).read()

# Stats
stats = {"total": len(plan), "passed": 0, "failed": 0, "total_cost": 0.0}


def get_assembled_context(max_chars=4000):
    """Read all assembled files as context for the next sub-task."""
    context = ""
    for root, dirs, files in os.walk(assembled):
        for f in files:
            if f.endswith(('.cjs', '.json')) and not f.startswith('.'):
                fp = os.path.join(root, f)
                rel = os.path.relpath(fp, assembled)
                try:
                    content = open(fp).read()
                    entry = f"\n--- {rel} ---\n{content}\n"
                    if len(context) + len(entry) < max_chars:
                        context += entry
                except:
                    pass
    return context or "None yet — this is the first sub-task."


def call_model(model, prompt_file, workdir):
    """Call a model and return metrics dict."""
    try:
        result = subprocess.run(
            ["python3", os.path.join(spike_dir, "call-model.py"),
             "--model", model,
             "--prompt-file", prompt_file,
             "--workdir", workdir],
            capture_output=True, text=True, timeout=300
        )
        if result.stdout.strip():
            return json.loads(result.stdout.strip())
    except subprocess.TimeoutExpired:
        return {"error": "timeout", "cost_usd": 0, "files_created": 0}
    except:
        pass
    return {"error": "call_failed", "cost_usd": 0, "files_created": 0}


def run_gate(gate_cmd, workdir):
    """Run a gate command in workdir. Return (pass, details)."""
    try:
        result = subprocess.run(
            ["bash", os.path.join(spike_dir, "validate-gate.sh"), gate_cmd, workdir],
            capture_output=True, text=True, timeout=30
        )
        data = json.loads(result.stdout.strip()) if result.stdout.strip() else {"pass": False}
        return data.get("pass", False), data
    except:
        return False, {"pass": False, "error": "gate_exception"}


def copy_to_assembled(task_dir):
    """Copy generated files from task_dir to assembled directory."""
    copied = 0
    for root, dirs, files in os.walk(task_dir):
        for f in files:
            if f.endswith(('.cjs', '.json', '.txt', '.go', '.graphql', '.ts', '.svelte', '.css', '.html', '.mod')) and f not in ('response.json', 'prompt.txt', 'gate-result.json'):
                src = os.path.join(root, f)
                rel = os.path.relpath(src, task_dir)
                dst = os.path.join(assembled, rel)
                os.makedirs(os.path.dirname(dst), exist_ok=True)
                shutil.copy2(src, dst)
                copied += 1
    return copied


# Pre-seed assembled dir with go.mod if this is a Go project
go_mod_path = os.path.join(assembled, "go.mod")
if not os.path.exists(go_mod_path):
    # Check if architecture mentions Go
    if "go.mod" in arch_content or "package model" in arch_content:
        os.makedirs(assembled, exist_ok=True)
        with open(go_mod_path, "w") as f:
            f.write("module task-board\n\ngo 1.22\n")
        os.makedirs(os.path.join(assembled, "model"), exist_ok=True)
        print("  Pre-seeded assembled/ with go.mod", flush=True)

# Execute each sub-task
for i, task in enumerate(plan):
    task_id = task.get("id", f"ST-{i+1:02d}")
    title = task.get("title", "untitled")
    task_dir = os.path.join(attempt_dir, task_id)
    os.makedirs(task_dir, exist_ok=True)

    print(f"  Sub-task {task_id}: {title}", flush=True)

    # Build prompt with architecture excerpt and context
    # Select prompt template based on project type and model
    go_prompt = os.path.join(spike_dir, "prompts", "executor-go.md")
    v4_path = os.path.join(spike_dir, "prompts", "executor-v4.md")
    default_path = os.path.join(spike_dir, "prompts", "executor.md")

    primary_file = task.get("files_to_create", ["output.go"])[0] if task.get("files_to_create") else "output.go"

    if primary_file.endswith(".go") and os.path.exists(go_prompt):
        prompt_template = open(go_prompt).read()
    elif "30b" in executor_model or "7b" in executor_model:
        prompt_template = open(v4_path).read() if os.path.exists(v4_path) else open(default_path).read()
    else:
        prompt_template = open(default_path).read()

    # Use shorter architecture excerpt (just the relevant module spec)
    arch_excerpt = arch_content[:3000]
    # Try to find the specific module section
    for marker in task.get("files_to_create", []):
        basename = os.path.basename(marker).replace('.cjs', '')
        idx = arch_content.find(f"### {marker}")
        if idx == -1:
            idx = arch_content.find(f"### lib/{basename}")
        if idx == -1:
            idx = arch_content.find(basename)
        if idx > 0:
            arch_excerpt = arch_content[max(0, idx-200):idx+1500]
            break

    context = get_assembled_context()

    prompt = prompt_template.replace("{{ARCHITECTURE}}", arch_excerpt)
    prompt = prompt.replace("{{SUB_TASK}}", json.dumps(task, indent=2))
    prompt = prompt.replace("{{CONTEXT_FILES}}", context)
    # V4 template needs primary file path
    primary_file = task.get("files_to_create", ["output.cjs"])[0] if task.get("files_to_create") else "output.cjs"
    prompt = prompt.replace("{{PRIMARY_FILE}}", primary_file)

    prompt_file = os.path.join(task_dir, "prompt.txt")
    with open(prompt_file, "w") as f:
        f.write(prompt)

    gate_cmd = task.get("gate_command", "")
    passed = False

    for retry in range(3):
        # Clear previous files (except prompt)
        for root, dirs, files in os.walk(task_dir):
            for f in files:
                if f not in ("prompt.txt",) and not f.startswith("prompt"):
                    os.remove(os.path.join(root, f))

        # Call executor
        metrics = call_model(executor_model, prompt_file, task_dir)
        cost = metrics.get("cost_usd", 0)
        stats["total_cost"] += cost if isinstance(cost, (int, float)) else 0
        files_created = metrics.get("files_created", 0)
        error = metrics.get("error", "")

        if error:
            print(f"    Executor error: {error} (retry {retry})", flush=True)
            continue

        if files_created == 0:
            print(f"    No files created (retry {retry})", flush=True)
            # Append failure context for retry
            with open(prompt_file, "a") as f:
                f.write("\n\n## PREVIOUS ATTEMPT PRODUCED NO FILES\nYou MUST output file blocks using --- FILE: path --- format.\n")
            continue

        # Copy to assembled first, then run gate there
        copied = copy_to_assembled(task_dir)

        # Auto-fix: run Go tooling to fix trivial errors before gate
        if os.path.exists(os.path.join(assembled, "go.mod")):
            # goimports fixes unused/missing imports
            subprocess.run(["goimports", "-w", "."], cwd=assembled,
                         capture_output=True, timeout=10)
            # gofmt fixes formatting
            subprocess.run(["gofmt", "-w", "."], cwd=assembled,
                         capture_output=True, timeout=10)

        # Compile check: run go build as a pre-gate for Go projects
        compile_pass = True
        compile_error = ""
        if os.path.exists(os.path.join(assembled, "go.mod")):
            build = subprocess.run(["go", "build", "./..."], cwd=assembled,
                                  capture_output=True, text=True, timeout=30)
            if build.returncode != 0:
                compile_pass = False
                compile_error = build.stderr[:300]
                print(f"    Compile: FAIL — {compile_error[:80]}", flush=True)

        # Run gate against assembled directory (where all deps are)
        if gate_cmd and compile_pass:
            gate_pass, gate_data = run_gate(gate_cmd, assembled)
            # Save gate result
            with open(os.path.join(task_dir, "gate-result.json"), "w") as f:
                json.dump(gate_data, f, indent=2)
        elif not compile_pass:
            gate_pass = False
            gate_data = {"pass": False, "stderr": compile_error}
            with open(os.path.join(task_dir, "gate-result.json"), "w") as f:
                json.dump(gate_data, f, indent=2)
        else:
            gate_pass = files_created > 0

        if gate_pass:
            print(f"    Gate: PASS (retry {retry}, {copied} files copied, ${cost})", flush=True)
            stats["passed"] += 1
            passed = True
            break
        else:
            print(f"    Gate: FAIL (retry {retry}, compile_error={bool(compile_error)})", flush=True)
            stderr = gate_data.get("stderr", "")[:300] if isinstance(gate_data, dict) else ""
            # Build specific retry instructions
            retry_msg = f"\n\n## PREVIOUS ATTEMPT FAILED — FIX THESE ERRORS\n"
            if compile_error:
                retry_msg += f"Compile error:\n{compile_error}\n"
                # Add Go-specific hints for common errors
                if "unknown escape" in compile_error or "invalid character" in compile_error:
                    retry_msg += "\nHINT: Use Go raw string literals (backticks `) for HTML/JS content, not double quotes.\n"
                    retry_msg += "Example: const html = `<div onclick=\"alert('hi')\">` — backtick strings don't interpret escapes.\n"
                if "undefined" in compile_error:
                    retry_msg += "\nHINT: Make sure all imported packages and referenced types are defined or imported.\n"
                if "unused" in compile_error:
                    retry_msg += "\nHINT: Remove unused variables and imports.\n"
            if stderr and stderr != compile_error:
                retry_msg += f"Gate error:\n{stderr}\n"
            retry_msg += "\nFix ONLY the errors above. Output the COMPLETE corrected file."
            with open(prompt_file, "a") as f:
                f.write(retry_msg)

    if not passed:
        print(f"    FAILED after 3 retries", flush=True)
        stats["failed"] += 1

# Final validation
print(f"\n  === Final Validation ===", flush=True)
print(f"  Sub-tasks: {stats['passed']}/{stats['total']} passed, {stats['failed']} failed", flush=True)
print(f"  Total cost: ${stats['total_cost']:.4f}", flush=True)

# List assembled files
print(f"  Assembled files:", flush=True)
for root, dirs, files in os.walk(assembled):
    for f in sorted(files):
        rel = os.path.relpath(os.path.join(root, f), assembled)
        print(f"    {rel}", flush=True)

# Final auto-fix pass on complete assembly
if os.path.exists(os.path.join(assembled, "go.mod")):
    subprocess.run(["goimports", "-w", "."], cwd=assembled, capture_output=True, timeout=10)
    subprocess.run(["gofmt", "-w", "."], cwd=assembled, capture_output=True, timeout=10)

# Run tests — detect project type
is_go = os.path.exists(os.path.join(assembled, "go.mod"))
is_node = os.path.exists(os.path.join(assembled, "test", "dep-doctor.test.cjs"))

if is_go:
    print("  Running Go tests on assembled output:", flush=True)
    try:
        # First try go build
        build = subprocess.run(["go", "build", "./..."], cwd=assembled,
                             capture_output=True, text=True, timeout=30)
        if build.returncode != 0:
            print(f"  Build failed:\n{build.stderr[:300]}", flush=True)

        # Run tests
        result = subprocess.run(["go", "test", "./...", "-v", "-count=1"], cwd=assembled,
                              capture_output=True, text=True, timeout=30)
        # Print just the test results (PASS/FAIL lines)
        for line in result.stdout.split("\n"):
            if "PASS" in line or "FAIL" in line or "ok " in line:
                print(f"  {line}", flush=True)
        if result.stderr:
            print(result.stderr[:300], flush=True)

        pass_count = result.stdout.count("--- PASS")
        fail_count = result.stdout.count("--- FAIL")
        print(f"  FINAL: {pass_count} PASS, {fail_count} FAIL", flush=True)
    except subprocess.TimeoutExpired:
        print("  FINAL: Tests timed out", flush=True)
    except Exception as e:
        print(f"  FINAL: Error running tests: {e}", flush=True)

elif is_node:
    test_file = os.path.join(assembled, "test", "dep-doctor.test.cjs")
    try:
        result = subprocess.run(["node", test_file], capture_output=True, text=True, timeout=30, cwd=assembled)
        print(result.stdout, flush=True)
        pass_count = result.stdout.count("PASS")
        fail_count = result.stdout.count("FAIL")
        print(f"  FINAL: {pass_count} PASS, {fail_count} FAIL", flush=True)
    except:
        print("  FINAL: Node tests failed to run", flush=True)

else:
    # Try golden master tests
    for gm_path in [
        os.path.join(spike_dir, "golden-master", "task-board"),
        os.path.join(spike_dir, "golden-master", "dep-doctor")
    ]:
        if os.path.exists(gm_path):
            # Copy golden master tests to assembled
            import glob
            for tf in glob.glob(os.path.join(gm_path, "**/*_test.go"), recursive=True):
                rel = os.path.relpath(tf, gm_path)
                dst = os.path.join(assembled, rel)
                os.makedirs(os.path.dirname(dst), exist_ok=True)
                shutil.copy2(tf, dst)
            print(f"  Copied golden master tests from {gm_path}", flush=True)
            try:
                result = subprocess.run(["go", "test", "./...", "-v", "-count=1"], cwd=assembled,
                                      capture_output=True, text=True, timeout=30)
                for line in result.stdout.split("\n"):
                    if "PASS" in line or "FAIL" in line or "ok " in line:
                        print(f"  {line}", flush=True)
                print(f"  FINAL: {result.stdout.count('--- PASS')} PASS, {result.stdout.count('--- FAIL')} FAIL", flush=True)
            except:
                print("  FINAL: Golden master tests failed", flush=True)
            break
    else:
        print("  FINAL: No test file found", flush=True)

# Write summary
summary = {
    "sub_tasks_total": stats["total"],
    "sub_tasks_passed": stats["passed"],
    "sub_tasks_failed": stats["failed"],
    "total_cost_usd": round(stats["total_cost"], 4),
    "executor_model": executor_model,
    "reviewer_model": reviewer_model
}
with open(os.path.join(attempt_dir, "summary.json"), "w") as f:
    json.dump(summary, f, indent=2)
print(f"\n  Summary: {json.dumps(summary)}", flush=True)
