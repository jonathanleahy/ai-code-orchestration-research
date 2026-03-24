#!/usr/bin/env python3
"""Extract file blocks from model output. Handles multiple formats:

1. --- FILE: path --- ... --- END FILE ---  (preferred)
2. ```language ... ``` with FILE: path hints nearby
3. Single code block with filename guessing

Usage:
  echo "response" | python3 parse-blocks.py --workdir ./out
  python3 parse-blocks.py --file response.txt --workdir ./out

Outputs JSON: {"files": ["path1", "path2"], "count": 2}
"""

import argparse
import json
import os
import re
import sys


def extract_blocks(text, workdir):
    created = []

    # Format 1: --- FILE: path --- blocks (preferred, most reliable)
    pattern = r"---\s*FILE:\s*([\w./\-]+)\s*---\n(.*?)\n---\s*END FILE\s*---"
    matches = re.findall(pattern, text, re.DOTALL)
    if matches:
        for filepath, content in matches:
            _write_file(workdir, filepath, content)
            created.append(filepath)
        return created

    # Format 2: Look for FILE: path hints near code blocks
    # e.g., "FILE: lib/parser.cjs\n```javascript\n...\n```"
    file_hint = r"(?:FILE|file|File|filename|Filename):\s*([\w./\-]+\.(?:cjs|js|json|sh|py|ts))"
    code_block = r"```(?:\w+)?\n(.*?)```"

    hints = re.findall(file_hint, text)
    blocks = re.findall(code_block, text, re.DOTALL)

    if hints and blocks:
        for filepath, code in zip(hints, blocks):
            _write_file(workdir, filepath, code.strip())
            created.append(filepath)
        return created

    # Format 3: Multiple code blocks with filenames in preceding lines
    # Look for patterns like "### lib/validator.cjs\n```javascript"
    sections = re.split(r'(?=```)', text)
    for section in sections:
        # Check if there's a filename before the code block
        fname_match = re.search(r'([\w./\-]+\.(?:cjs|js|json|sh|py))\s*\n```', section)
        code_match = re.search(r'```(?:\w+)?\n(.*?)```', section, re.DOTALL)
        if fname_match and code_match:
            filepath = fname_match.group(1)
            code = code_match.group(1).strip()
            if code and len(code) > 10:
                _write_file(workdir, filepath, code)
                created.append(filepath)

    if created:
        return created

    # Format 4: Single code block — write as output.cjs
    if blocks:
        code = blocks[0].strip()
        if len(code) > 20 and ('module.exports' in code or 'require(' in code or "'use strict'" in code):
            _write_file(workdir, "output.cjs", code)
            created.append("output.cjs")

    return created


def _write_file(workdir, filepath, content):
    full_path = os.path.join(workdir, filepath)
    os.makedirs(os.path.dirname(full_path), exist_ok=True)
    with open(full_path, "w") as f:
        f.write(content)
    os.chmod(full_path, 0o755)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--workdir", required=True)
    parser.add_argument("--file", help="Input file (default: stdin)")
    args = parser.parse_args()

    if args.file:
        with open(args.file) as f:
            text = f.read()
    else:
        text = sys.stdin.read()

    files = extract_blocks(text, args.workdir)
    print(json.dumps({"files": files, "count": len(files)}))


if __name__ == "__main__":
    main()
