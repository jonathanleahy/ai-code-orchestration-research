#!/usr/bin/env python3
"""Unified model caller — OpenRouter API or claude -p.

Usage:
  python3 call-model.py --model qwen/qwen3-coder --prompt-file prompt.txt --workdir ./out
  python3 call-model.py --model claude-sonnet --prompt-file prompt.txt --workdir ./out

Env: OPENROUTER_API_KEY (for non-Claude models)

Outputs:
  - workdir/response.json (raw API response or claude log)
  - Extracted files written to workdir (from --- FILE: --- blocks)
  - JSON metrics to stdout: {cost, tokens_in, tokens_out, files_created, time_s}
"""

import argparse
import json
import os
import re
import subprocess
import sys
import time
import urllib.request
import urllib.error


def call_openrouter(model, prompt, max_tokens=8192):
    """Call OpenRouter API, return (content, usage_dict, raw_response)."""
    api_key = os.environ.get("OPENROUTER_API_KEY", "")
    if not api_key:
        return "", {"error": "OPENROUTER_API_KEY not set"}, {}

    payload = json.dumps({
        "model": model,
        "messages": [{"role": "user", "content": prompt}],
        "max_tokens": max_tokens,
        "temperature": 0.2
    })

    req = urllib.request.Request(
        "https://openrouter.ai/api/v1/chat/completions",
        data=payload.encode(),
        headers={
            "Authorization": f"Bearer {api_key}",
            "Content-Type": "application/json",
            "HTTP-Referer": "https://gyrum.ai"
        }
    )

    try:
        with urllib.request.urlopen(req, timeout=300) as resp:
            data = json.loads(resp.read())
    except urllib.error.HTTPError as e:
        body = e.read().decode()[:300] if e.fp else ""
        return "", {"error": f"HTTP {e.code}: {body}"}, {}
    except Exception as e:
        return "", {"error": str(e)}, {}

    content = data.get("choices", [{}])[0].get("message", {}).get("content", "")
    usage = data.get("usage", {})
    return content, usage, data


def call_claude_p(prompt, workdir, model="sonnet", budget="0.50"):
    """Call claude -p with imperative prompt, return (content, usage_dict, raw_log)."""
    cmd = [
        "claude", "-p",
        "--dangerously-skip-permissions",
        "--model", model,
        "--max-budget-usd", budget,
        "--output-format", "stream-json",
        "--verbose",
        prompt
    ]

    try:
        result = subprocess.run(
            cmd, capture_output=True, text=True, timeout=300, cwd=workdir
        )
        log = result.stdout + result.stderr
    except subprocess.TimeoutExpired:
        return "", {"error": "timeout"}, ""
    except Exception as e:
        return "", {"error": str(e)}, ""

    # Extract cost from stream-json
    cost = 0.0
    for line in log.split("\n"):
        if '"total_cost_usd"' in line:
            try:
                match = re.search(r'"total_cost_usd":([\d.]+)', line)
                if match:
                    cost = float(match.group(1))
            except:
                pass

    # Claude -p writes files via tools, not file blocks.
    # Content is empty (files are on disk), but we return the log.
    usage = {"total_cost_usd": cost, "prompt_tokens": 0, "completion_tokens": 0}
    return "", usage, log


def extract_file_blocks(content, workdir):
    """Extract file blocks from model output. Tries multiple formats:
    1. --- FILE: path --- ... --- END FILE ---
    2. ```language ... ``` with FILE: path hint
    3. Single code block when only one file expected
    """
    created = []

    # Format 1: --- FILE: path --- blocks (preferred)
    pattern = r"---\s*FILE:\s*([\w./\-]+)\s*---\n(.*?)\n---\s*END FILE\s*---"
    matches = re.findall(pattern, content, re.DOTALL)

    if matches:
        for filepath, filecontent in matches:
            full_path = os.path.join(workdir, filepath)
            os.makedirs(os.path.dirname(full_path), exist_ok=True)
            with open(full_path, "w") as f:
                f.write(filecontent)
            os.chmod(full_path, 0o755)
            created.append(filepath)
        return created

    # Format 2: Look for FILE: path hints near code blocks
    file_hint_pattern = r"(?:FILE|file|File):\s*([\w./\-]+\.(?:cjs|js|json|sh|py))\s*(?:\n|---)"
    code_block_pattern = r"```(?:\w+)?\n(.*?)```"

    hints = re.findall(file_hint_pattern, content)
    blocks = re.findall(code_block_pattern, content, re.DOTALL)

    if hints and blocks:
        for filepath, code in zip(hints, blocks):
            full_path = os.path.join(workdir, filepath)
            os.makedirs(os.path.dirname(full_path), exist_ok=True)
            with open(full_path, "w") as f:
                f.write(code.strip())
            os.chmod(full_path, 0o755)
            created.append(filepath)
        return created

    # Format 3: Single code block — guess filename from content
    if blocks and not created:
        code = blocks[0].strip()
        if code.startswith("'use strict'") or code.startswith('"use strict"') or 'module.exports' in code:
            # Try to guess the filename from the prompt context
            # Default to output.cjs
            filepath = "output.cjs"
            full_path = os.path.join(workdir, filepath)
            with open(full_path, "w") as f:
                f.write(code)
            os.chmod(full_path, 0o755)
            created.append(filepath)

    return created


def main():
    parser = argparse.ArgumentParser(description="Unified model caller")
    parser.add_argument("--model", required=True, help="Model ID (e.g., qwen/qwen3-coder or claude-sonnet)")
    parser.add_argument("--prompt-file", required=True, help="Path to prompt file")
    parser.add_argument("--workdir", required=True, help="Output directory")
    parser.add_argument("--max-tokens", type=int, default=8192, help="Max output tokens")
    parser.add_argument("--budget", default="0.50", help="Budget for claude -p (ignored for OpenRouter)")
    args = parser.parse_args()

    os.makedirs(args.workdir, exist_ok=True)

    with open(args.prompt_file) as f:
        prompt = f.read()

    start = time.time()

    if args.model.startswith("claude"):
        # Use claude -p (subscription, free)
        claude_model = args.model.replace("claude-", "")  # claude-sonnet → sonnet
        content, usage, raw = call_claude_p(prompt, args.workdir, claude_model, args.budget)
        # Claude writes files via tools, so check workdir for files
        files_created = []
        for root, dirs, files in os.walk(args.workdir):
            for f in files:
                if f.endswith(('.js', '.cjs', '.json', '.md', '.sh')) and f != 'response.json':
                    files_created.append(os.path.relpath(os.path.join(root, f), args.workdir))
        # Save log as response
        with open(os.path.join(args.workdir, "response.json"), "w") as f:
            json.dump({"log": raw[:5000] if isinstance(raw, str) else "", "usage": usage}, f)
    else:
        # Use OpenRouter API
        content, usage, raw = call_openrouter(args.model, prompt, args.max_tokens)
        # Save raw response
        with open(os.path.join(args.workdir, "response.json"), "w") as f:
            json.dump(raw, f, indent=2)
        # Extract file blocks
        files_created = extract_file_blocks(content, args.workdir)

    elapsed = time.time() - start

    # Output metrics
    error = usage.get("error", "")
    metrics = {
        "model": args.model,
        "cost_usd": usage.get("total_cost_usd", 0) or (usage.get("total_tokens", 0) * 0.000003),
        "tokens_in": usage.get("prompt_tokens", 0),
        "tokens_out": usage.get("completion_tokens", 0),
        "files_created": len(files_created),
        "file_paths": files_created,
        "time_s": round(elapsed, 1),
        "error": error
    }
    print(json.dumps(metrics))


if __name__ == "__main__":
    main()
