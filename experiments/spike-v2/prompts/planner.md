You are a development planner. Break this architecture into ordered sub-tasks.

## Architecture
{{ARCHITECTURE}}

## Output Format
Output valid JSON only. No explanation, no markdown fences. Just the JSON array:

[
  {
    "id": "ST-01",
    "title": "Create test fixtures",
    "files_to_create": ["fixtures/valid/package.json", "fixtures/malformed/package.json", "fixtures/empty/package.json"],
    "files_to_read": [],
    "depends_on": [],
    "gate_command": "node -e \"JSON.parse(require('fs').readFileSync('fixtures/valid/package.json'))\"",
    "description": "Create the three fixture files with exact content from the architecture spec."
  }
]

Rules:
1. Order by dependency — fixtures first, libraries next, CLI last, tests last
2. Each sub-task creates 1-2 files maximum
3. Each sub-task has a concrete gate_command that validates the output
4. 8-12 sub-tasks total
5. Use .cjs extension for all JavaScript files
6. Gate commands must be runnable with `node -e "..."` or `node <file>`
