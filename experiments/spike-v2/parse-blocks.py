#!/usr/bin/env python3
"""Extract --- FILE: path --- blocks from stdin or a file and write to a directory.

Usage:
  echo "response" | python3 parse-blocks.py --workdir ./out
  python3 parse-blocks.py --file response.txt --workdir ./out

Outputs JSON to stdout: {"files": ["path1", "path2"], "count": 2}
"""

import argparse
import json
import os
import re
import sys


def extract_blocks(text, workdir):
    pattern = r"---\s*FILE:\s*([\w./\-]+)\s*---\n(.*?)\n---\s*END FILE\s*---"
    matches = re.findall(pattern, text, re.DOTALL)

    created = []
    for filepath, content in matches:
        full_path = os.path.join(workdir, filepath)
        os.makedirs(os.path.dirname(full_path), exist_ok=True)
        with open(full_path, "w") as f:
            f.write(content)
        os.chmod(full_path, 0o755)
        created.append(filepath)

    return created


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
