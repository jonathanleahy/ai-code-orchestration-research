#!/usr/bin/env python3
"""Call OpenRouter API and parse response into file blocks.

Usage: python3 call-openrouter.py <model> <workdir> <queue-add-path>

Env: OPENROUTER_API_KEY must be set.
Writes: workdir/tests/queue-add-model.test.sh, workdir/scripts/queue-add.sh, workdir/response.json
Outputs: JSON with metrics to stdout.
"""

import json, os, re, sys, urllib.request, urllib.error

model = sys.argv[1]
workdir = sys.argv[2]
queue_add_path = sys.argv[3]
api_key = os.environ["OPENROUTER_API_KEY"]

# Read current code
with open(queue_add_path) as f:
    current_code = f.read()
current_short = "\n".join(current_code.split("\n")[:100])

prompt = f"""You are modifying a bash script. Your task: add a --model flag to queue-add.sh.

## Current Code (first 100 lines of scripts/queue-add.sh)
```bash
{current_short}
```

## What to Add
1. Add VALID_MODELS=(default sonnet opus minimax-m2.7) near the VALID_SIZES line
2. Add --model) flag parsing in the existing case statement (same pattern as --size)
3. Validate model against VALID_MODELS (same pattern as size validation)
4. Include the model variable in the JSON body sent to the API

## Test Scenarios
- --model sonnet → accepted, model=sonnet
- --model gpt-4 → error, invalid model
- No --model → defaults to "default"
- --model minimax-m2.7 → accepted (dot in name)

## Output Format
Output EXACTLY two file blocks. Keep ALL existing code in queue-add.sh — only ADD the model support.

--- FILE: tests/queue-add-model.test.sh ---
[complete bash test script]
--- END FILE ---

--- FILE: scripts/queue-add.sh ---
[COMPLETE modified file — preserve ALL existing code, only ADD model support]
--- END FILE ---

Output ONLY the two file blocks. No explanation."""

payload = json.dumps({
    "model": model,
    "messages": [{"role": "user", "content": prompt}],
    "max_tokens": 8192,
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
    body = e.read().decode()[:200] if e.fp else ""
    data = {"error": {"message": f"HTTP {e.code}: {body}"}}
except Exception as e:
    data = {"error": {"message": str(e)}}

# Save raw response
with open(f"{workdir}/response.json", "w") as f:
    json.dump(data, f, indent=2)

error = data.get("error", {}).get("message", "")
if error:
    print(json.dumps({"error": error, "tokens": 0, "tests_created": 0, "files_changed": 0, "has_valid_models": False, "preserved": False}))
    sys.exit(0)

content = data.get("choices", [{}])[0].get("message", {}).get("content", "")
usage = data.get("usage", {})
tokens = usage.get("total_tokens", 0)

# Parse file blocks
def extract_block(text, filename):
    pattern = rf"---\s*FILE:\s*{re.escape(filename)}\s*---\n(.*?)\n---\s*END FILE\s*---"
    match = re.search(pattern, text, re.DOTALL)
    return match.group(1) if match else None

test_content = extract_block(content, "tests/queue-add-model.test.sh")
impl_content = extract_block(content, "scripts/queue-add.sh")

tests_created = 0
if test_content:
    os.makedirs(f"{workdir}/tests", exist_ok=True)
    with open(f"{workdir}/tests/queue-add-model.test.sh", "w") as f:
        f.write(test_content)
    os.chmod(f"{workdir}/tests/queue-add-model.test.sh", 0o755)
    tests_created = 1

files_changed = 0
if impl_content:
    with open(f"{workdir}/scripts/queue-add.sh", "w") as f:
        f.write(impl_content)
    os.chmod(f"{workdir}/scripts/queue-add.sh", 0o755)
    files_changed = 1

# Quality checks
try:
    impl_text = open(f"{workdir}/scripts/queue-add.sh").read()
except:
    impl_text = ""

has_valid_models = "VALID_MODELS" in impl_text
has_model_flag = "--model" in impl_text
impl_lines = len(impl_text.split("\n"))
orig_lines = len(current_code.split("\n"))
preserved = impl_lines >= (orig_lines * 0.5)

print(json.dumps({
    "tests_created": tests_created,
    "files_changed": files_changed,
    "has_valid_models": has_valid_models,
    "has_model_flag": has_model_flag,
    "impl_lines": impl_lines,
    "orig_lines": orig_lines,
    "preserved": preserved,
    "tokens": tokens
}))
