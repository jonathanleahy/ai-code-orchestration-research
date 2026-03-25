#!/usr/bin/env python3
"""Hardened file block parser (v2) — handles all known model output formats.

Extracts files from model output text and writes them to a workdir.
Returns list of created file paths.

Formats handled (in priority order):
1. --- FILE: path --- ... --- END FILE --- (canonical)
2. ```lang\n// FILE: path\n... ``` (code fence with filename comment)
3. ```lang filename\n... ``` (code fence with filename on opening line)
4. FILE: path hints near code blocks (loose matching)
5. Raw Go code — split on package declarations
6. Single code block fallback

Usage as module:
  from parse_blocks_v2 import extract_files

Usage as CLI:
  echo "model output" | python3 parse-blocks-v2.py --workdir ./out
"""

import json
import os
import re
import sys


def extract_files(content, workdir, expected_files=None):
    """Extract files from model output. Returns list of (path, content) tuples.

    Args:
        content: Raw model output text
        workdir: Directory to write files to
        expected_files: Optional list of expected file paths (helps disambiguation)

    Returns:
        List of created file paths (relative to workdir)
    """
    created = []

    # Strip any leading/trailing chat wrapper
    content = content.strip()

    # Pre-process: strip outer code fences that wrap FILE blocks
    # Models sometimes output ```go\n--- FILE: ... ---\n...\n--- END FILE ---\n```
    if re.match(r'^```\w*\n', content) and '--- FILE:' in content:
        content = re.sub(r'^```\w*\n', '', content)
        content = re.sub(r'\n```\s*$', '', content)

    # Pre-process: fix truncated END FILE markers (model ran out of tokens)
    # e.g., "--- END FILE --`" or "--- END FILE -" or just missing entirely
    content = re.sub(r'---\s*END\s*FILE\s*--[`\-]?\s*$', '--- END FILE ---', content)
    content = re.sub(r'---\s*END\s*FILE\s*-?\s*$', '--- END FILE ---', content)

    # === FORMAT 1: --- FILE: path --- blocks (canonical) ===
    # Very lenient: 2-4 dashes, optional trailing dashes, optional whitespace
    # Also handles: == FILE: path ==, ** FILE: path **
    f1_pattern = r'[-=*]{2,4}\s*FILE:\s*([\w./\-]+(?:\.\w+)?)\s*[-=*]{0,4}\n(.*?)(?:\n[-=*]{2,4}\s*END\s*FILE\s*[-=*]{0,4}|\n[-=*]{2,4}\s*FILE:|\Z)'
    f1_matches = list(re.finditer(f1_pattern, content, re.DOTALL))

    if f1_matches:
        for m in f1_matches:
            filepath = m.group(1).strip()
            filecontent = m.group(2)
            # Clean: remove leading/trailing blank lines
            filecontent = filecontent.strip('\n')
            if filepath and filecontent and len(filecontent) > 5:
                _write_file(workdir, filepath, filecontent, created)
        if created:
            return created

    # === FORMAT 1b: Splitting approach for truncated END markers ===
    file_markers = list(re.finditer(r'[-=*]{2,4}\s*FILE:\s*([\w./\-]+(?:\.\w+)?)\s*[-=*]{0,4}\n', content))
    if len(file_markers) > len(created):
        for i, m in enumerate(file_markers):
            filepath = m.group(1).strip()
            start = m.end()
            # Content goes until next FILE marker or end
            if i + 1 < len(file_markers):
                end = file_markers[i+1].start()
            else:
                end = len(content)
            body = content[start:end]
            # Strip END FILE marker if present
            body = re.sub(r'\n[-=*]{2,4}\s*END\s*FILE\s*[-=*]{0,4}\s*$', '', body)
            body = body.strip('\n')
            if filepath and body and len(body) > 5 and filepath not in created:
                _write_file(workdir, filepath, body, created)
        if created:
            return created

    # === FORMAT 2: Code fences with filename in comment or heading ===
    # Matches: ```go\n// FILE: model/task.go\n...```
    # Also: ```go\n// model/task.go\n...```
    # Also: ### model/task.go\n```go\n...```
    f2_blocks = list(re.finditer(
        r'(?:^#+\s*([\w./\-]+(?:\.\w+))\s*\n)?'  # optional ### heading with filename
        r'```(\w+)?\n'                              # opening fence
        r'(?://\s*(?:FILE:\s*)?([\w./\-]+(?:\.\w+)?)\s*\n)?'  # optional // FILE: comment
        r'(.*?)'                                    # content
        r'```',                                     # closing fence
        content, re.DOTALL | re.MULTILINE
    ))

    if f2_blocks:
        for m in f2_blocks:
            heading_path = m.group(1)
            lang = m.group(2) or ''
            comment_path = m.group(3)
            code = m.group(4)

            filepath = comment_path or heading_path
            if not filepath and code and len(code.strip()) > 20:
                # Try to infer from content
                filepath = _infer_filepath(code.strip(), lang, expected_files)

            if filepath and code and len(code.strip()) > 5:
                _write_file(workdir, filepath, code.strip(), created)

        if created:
            return created

    # === FORMAT 3: FILE: hints near code blocks (loose) ===
    # Extended to match .go, .graphql, .mod, .ts, .html, .svelte etc.
    file_ext = r'\.(?:go|graphql|mod|cjs|js|ts|json|sh|py|html|svelte|css|sql|yaml|yml|toml|md)'
    hint_pattern = rf'(?:FILE|file|File|Filename|filename):\s*([\w./\-]+{file_ext})\s*'
    block_pattern = r'```(?:\w+)?\n(.*?)```'

    hints = re.findall(hint_pattern, content)
    blocks = re.findall(block_pattern, content, re.DOTALL)

    if hints and blocks:
        for filepath, code in zip(hints, blocks):
            code = code.strip()
            if filepath and code and len(code) > 5:
                _write_file(workdir, filepath, code, created)
        if created:
            return created

    # === FORMAT 4: Raw Go code — split on package declarations ===
    if not created and ('package ' in content or '```go' in content):
        # First try ```go blocks
        go_blocks = re.findall(r'```go\n(.*?)```', content, re.DOTALL)

        if not go_blocks:
            # Split on package declarations
            parts = re.split(r'(?=^package \w+)', content, flags=re.MULTILINE)
            go_blocks = [p.strip() for p in parts if p.strip().startswith('package ')]

        for block in go_blocks:
            block = block.strip()
            if not block or len(block) < 20:
                continue
            filepath = _infer_go_filepath(block, expected_files)
            if filepath not in created:
                _write_file(workdir, filepath, block, created)

        if created:
            return created

    # === FORMAT 5: Single code block fallback ===
    if blocks and not created:
        code = blocks[0].strip()
        if len(code) > 20:
            filepath = _infer_filepath(code, '', expected_files)
            if filepath:
                _write_file(workdir, filepath, code, created)

    return created


def _write_file(workdir, filepath, content, created_list):
    """Write a file to workdir, create dirs as needed."""
    # Sanitize path — no ../ or absolute paths
    filepath = filepath.lstrip('/')
    if '..' in filepath:
        return

    full_path = os.path.join(workdir, filepath)
    os.makedirs(os.path.dirname(full_path) or workdir, exist_ok=True)
    with open(full_path, "w") as f:
        f.write(content)
    if filepath.endswith(('.sh', '.py')):
        os.chmod(full_path, 0o755)
    created_list.append(filepath)


def _infer_go_filepath(code, expected_files=None):
    """Infer Go file path from code content."""
    if expected_files:
        for ef in expected_files:
            if ef.endswith('_test.go') and 'func Test' in code:
                return ef
            if ef == 'main.go' and 'func main()' in code:
                return ef
            if ef.endswith('.go') and not ef.endswith('_test.go'):
                # Check if the package matches
                pkg_match = re.search(r'^package (\w+)', code)
                if pkg_match:
                    pkg = pkg_match.group(1)
                    if pkg in ef or (pkg == 'model' and 'model/' in ef):
                        return ef

    # Heuristic inference
    if 'func Test' in code and '"testing"' in code:
        if 'package model' in code:
            return "model/task_test.go"
        return "main_test.go"
    if 'type Task struct' in code or 'type Store struct' in code:
        return "model/task.go"
    if 'func main()' in code:
        return "main.go"
    if 'package model' in code:
        return "model/task.go"

    return "output.go"


def _infer_filepath(code, lang, expected_files=None):
    """Infer file path from code content and language hint."""
    if expected_files and len(expected_files) == 1:
        return expected_files[0]

    if lang == 'go' or code.startswith('package '):
        return _infer_go_filepath(code, expected_files)

    if lang == 'graphql' or 'type Task {' in code:
        return 'schema.graphql'

    if lang in ('javascript', 'js') or "'use strict'" in code or 'module.exports' in code:
        return 'output.cjs'

    if lang in ('typescript', 'ts') or 'interface ' in code or 'export ' in code:
        return 'output.ts'

    if lang == 'html' or '<!DOCTYPE' in code.upper() or '<html' in code.lower():
        return 'index.html'

    if lang == 'python' or code.startswith('#!/usr/bin/env python') or 'def ' in code:
        return 'output.py'

    if lang == 'bash' or lang == 'sh' or code.startswith('#!/'):
        return 'output.sh'

    return None


# === TEST SUITE ===

def run_tests():
    """Run parser tests against known model output patterns."""
    import tempfile

    passed = 0
    failed = 0

    tests = [
        # Test 1: Canonical format
        {
            "name": "canonical --- FILE: --- blocks",
            "input": """--- FILE: model/task.go ---
package model

type Task struct {
    ID string
}
--- END FILE ---

--- FILE: model/task_test.go ---
package model

import "testing"

func TestCreate(t *testing.T) {
    // test
}
--- END FILE ---""",
            "expected": ["model/task.go", "model/task_test.go"]
        },
        # Test 2: Truncated END marker
        {
            "name": "truncated END marker (last block)",
            "input": """--- FILE: model/task.go ---
package model

type Task struct { ID string }
--- END FILE ---

--- FILE: model/task_test.go ---
package model

import "testing"

func TestCreate(t *testing.T) {}""",
            "expected": ["model/task.go", "model/task_test.go"]
        },
        # Test 3: Code fences with ```go
        {
            "name": "```go code fences (no FILE markers)",
            "input": """Here is the model layer:

```go
package model

type Task struct {
    ID string
}

type Store struct {
    tasks map[string]*Task
}
```

And here is the main server:

```go
package main

import "net/http"

func main() {
    http.ListenAndServe(":8890", nil)
}
```""",
            "expected_contains": ["model/task.go", "main.go"]
        },
        # Test 4: FILE: hints with code blocks
        {
            "name": "FILE: hints near code blocks (Go files)",
            "input": """FILE: model/task.go

```go
package model

type Task struct { ID string }
```

FILE: schema.graphql

```graphql
type Task {
    id: ID!
    title: String!
}
```""",
            "expected": ["model/task.go", "schema.graphql"]
        },
        # Test 5: Two dashes instead of three
        {
            "name": "-- FILE: -- (two dashes)",
            "input": """-- FILE: model/task.go --
package model

type Task struct { ID string }
-- END FILE --

-- FILE: main.go --
package main

func main() {}
-- END FILE --""",
            "expected": ["model/task.go", "main.go"]
        },
        # Test 6: Code fence with filename comment
        {
            "name": "```go with // FILE: comment",
            "input": """```go
// FILE: model/task.go
package model

type Task struct {
    ID string
}
```

```go
// FILE: model/task_test.go
package model

import "testing"

func TestCreate(t *testing.T) {}
```""",
            "expected": ["model/task.go", "model/task_test.go"]
        },
        # Test 7: Heading + code block
        {
            "name": "### heading + code block",
            "input": """### schema.graphql
```graphql
type Task {
    id: ID!
    title: String!
}
```

### model/task.go
```go
package model

type Task struct { ID string }
```""",
            "expected": ["schema.graphql", "model/task.go"]
        },
        # Test 8: Mixed formats
        {
            "name": "model output with explanation text mixed in",
            "input": """I'll create both files for you.

First, the model layer:

--- FILE: model/task.go ---
package model

import "sync"

type Task struct {
    ID    string
    Title string
}

type Store struct {
    mu    sync.RWMutex
    tasks map[string]*Task
}
--- END FILE ---

Now the tests:

--- FILE: model/task_test.go ---
package model

import "testing"

func TestCreate(t *testing.T) {
    s := &Store{tasks: make(map[string]*Task)}
    _ = s
}
--- END FILE ---

That should work!""",
            "expected": ["model/task.go", "model/task_test.go"]
        },
    ]

    for test in tests:
        with tempfile.TemporaryDirectory() as tmpdir:
            result = extract_files(test["input"], tmpdir)

            if "expected" in test:
                expected = set(test["expected"])
                actual = set(result)
                if actual == expected:
                    passed += 1
                    print(f"  PASS: {test['name']}")
                else:
                    failed += 1
                    print(f"  FAIL: {test['name']}")
                    print(f"    expected: {sorted(expected)}")
                    print(f"    actual:   {sorted(actual)}")
            elif "expected_contains" in test:
                missing = [f for f in test["expected_contains"] if f not in result]
                if not missing:
                    passed += 1
                    print(f"  PASS: {test['name']} (got {result})")
                else:
                    failed += 1
                    print(f"  FAIL: {test['name']}")
                    print(f"    missing: {missing}")
                    print(f"    actual:  {result}")

    print(f"\n  {passed}/{passed+failed} tests passed")
    return failed == 0


if __name__ == "__main__":
    if "--test" in sys.argv:
        print("Running parser tests...")
        success = run_tests()
        sys.exit(0 if success else 1)
    elif "--workdir" in sys.argv:
        idx = sys.argv.index("--workdir")
        workdir = sys.argv[idx + 1] if idx + 1 < len(sys.argv) else "."
        content = sys.stdin.read()
        files = extract_files(content, workdir)
        print(json.dumps({"files_created": len(files), "file_paths": files}))
    else:
        print("Usage:")
        print("  python3 parse-blocks-v2.py --test")
        print("  echo 'content' | python3 parse-blocks-v2.py --workdir ./out")
